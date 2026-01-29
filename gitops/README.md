# GitOps CI/CD for Warehouse Inventory Demo

This directory contains ComponentWorkflowRun manifests to trigger GitOps-based builds and releases.

## Components

| Component | Source | Deployment Method |
|-----------|--------|-------------------|
| `inventory-api` | `/demo/inventory-api` | Direct (kubectl apply) |
| `inventory-dashboard` | `/demo/inventory-dashboard` | Direct (kubectl apply) |
| `inventory-api-gitops` | `/demo/inventory-api` | **GitOps workflow** |
| `inventory-dashboard-gitops` | `/demo/inventory-dashboard` | **GitOps workflow** |

The `-gitops` components use the same source code but are deployed via the GitOps workflow, allowing you to test GitOps without affecting the working demo.

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         GitOps Workflow Flow                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. kubectl apply -f inventory-api-gitops-run.yaml                          │
│           │                                                                  │
│           ▼                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    BUILD PHASE (Argo Workflow)                       │    │
│  │  clone-source → build-image → push-image → extract-descriptor        │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│           │                                                                  │
│           ▼                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                   RELEASE PHASE (Argo Workflow)                      │    │
│  │  clone-gitops → create-branch → generate-manifests → create-pr       │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│           │                                                                  │
│           ▼                                                                  │
│  2. Pull Request created in GitOps repository                               │
│           │                                                                  │
│           ▼                                                                  │
│  3. Review and merge PR                                                      │
│           │                                                                  │
│           ▼                                                                  │
│  4. OpenChoreo syncs and deploys                                            │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Quick Start

### Step 1: Install Workflow Templates

```bash
kubectl apply -f samples/gitops-workflows/component-workflows/build-and-release/docker/docker-gitops-release-template.yaml
kubectl apply -f samples/gitops-workflows/component-workflows/build-and-release/docker/docker-gitops-release.yaml

# Verify
kubectl get clusterworkflowtemplate docker-gitops-release
kubectl get componentworkflow docker-gitops-release
```

### Step 2: Configure GitHub Tokens

```bash
# Set your tokens
SOURCE_GIT_TOKEN="ghp_your_source_repo_token"      # For private source repos
GITOPS_GIT_TOKEN="ghp_your_gitops_repo_token"      # Required - needs 'repo' scope

# Add to ClusterSecretStore
kubectl patch clustersecretstore default --type='json' -p="[
  {\"op\": \"add\", \"path\": \"/spec/provider/fake/data/-\", \"value\": {\"key\": \"git-token\", \"value\": \"${SOURCE_GIT_TOKEN}\"}},
  {\"op\": \"add\", \"path\": \"/spec/provider/fake/data/-\", \"value\": {\"key\": \"gitops-token\", \"value\": \"${GITOPS_GIT_TOKEN}\"}}
]"

# Verify
kubectl get clustersecretstore default -o jsonpath='{.spec.provider.fake.data[*].key}' | tr ' ' '\n'
```

### Step 3: Create the GitOps Components

```bash
kubectl apply -f demo/.openchoreo/inventory-api-gitops.yaml
kubectl apply -f demo/.openchoreo/inventory-dashboard-gitops.yaml

# Verify
kubectl get components
```

### Step 4: Trigger GitOps Builds

```bash
# Build inventory-api-gitops
kubectl apply -f demo/gitops/inventory-api-gitops-run.yaml

# Build inventory-dashboard-gitops
kubectl apply -f demo/gitops/inventory-dashboard-gitops-run.yaml
```

### Step 5: Monitor Progress

```bash
# Watch ComponentWorkflowRun status
kubectl get componentworkflowrun -w

# View Argo Workflow
kubectl get workflow -n openchoreo-ci-default

# View logs
kubectl logs -n openchoreo-ci-default -l workflows.argoproj.io/workflow=inventory-api-gitops-release-001 --all-containers -f
```

### Step 6: Merge the PR

Once the workflow completes, a PR will be created in your GitOps repository. Merge it to trigger deployment.

## Re-triggering Builds

To re-run a workflow, increment the name:

```yaml
metadata:
  name: inventory-api-gitops-release-002  # Changed from 001
```

Then apply:
```bash
kubectl apply -f demo/gitops/inventory-api-gitops-run.yaml
```

## Files

```
demo/.openchoreo/
├── inventory-api-gitops.yaml           # Component definition
└── inventory-dashboard-gitops.yaml     # Component definition

demo/gitops/
├── README.md                           # This file
├── inventory-api-gitops-run.yaml       # Workflow trigger for API
└── inventory-dashboard-gitops-run.yaml # Workflow trigger for Dashboard
```

## GitOps Repository Setup

Your GitOps repository needs these resources for the workflow to succeed:

```
warehouse-inventory-gitops/
├── projects/
│   └── default/
│       └── project.yaml
├── components/
│   ├── inventory-api-gitops/
│   │   └── component.yaml
│   └── inventory-dashboard-gitops/
│       └── component.yaml
├── environments/
│   └── development/
│       └── environment.yaml
└── pipelines/
    └── simple-pipeline/
        └── pipeline.yaml
```

## Troubleshooting

| Issue | Check |
|-------|-------|
| Workflow not starting | `kubectl describe componentworkflowrun <name>` |
| Clone fails | Token permissions, repo URL |
| Build fails | Dockerfile path, build context |
| PR not created | `gitops-token` has `repo` scope |

```bash
# Check ExternalSecrets
kubectl get externalsecret -n openchoreo-ci-default

# Check workflow pods
kubectl get pods -n openchoreo-ci-default

# Get detailed logs
kubectl logs -n openchoreo-ci-default <pod-name> -c main
```
