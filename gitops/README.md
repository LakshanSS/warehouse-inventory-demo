# GitOps CI/CD for Warehouse Inventory Demo

This directory contains ComponentWorkflowRun manifests to trigger GitOps-based builds and releases.

## Components

| Component | Source | Deployment Method | Endpoint |
|-----------|--------|-------------------|----------|
| `inventory-api` | `/inventory-api` | Direct (kubectl apply) | `/inventory-api/` |
| `inventory-dashboard` | `/inventory-dashboard` | Direct (kubectl apply) | `inventory-dashboard-development.*` |
| `inventory-api-gitops` | `/inventory-api` | **GitOps workflow** | `/<release-name>/` |
| `inventory-dashboard-gitops` | `/inventory-dashboard` | **GitOps workflow** | `<release-name>-development.*` |

The `-gitops` components use the same source code but are deployed via the GitOps workflow.

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
│  4. Flux auto-syncs (or manual kubectl apply)                               │
│           │                                                                  │
│           ▼                                                                  │
│  5. OpenChoreo deploys to data plane                                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Prerequisites

- OpenChoreo installed with build plane
- GitHub Personal Access Token with `repo` scope
- GitOps repository: https://github.com/LakshanSS/warehouse-inventory-gitops

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
SOURCE_GIT_TOKEN="ghp_your_source_repo_token"      # For private source repos (optional for public)
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
kubectl get components | grep gitops
```

### Step 4: Trigger GitOps Builds

```bash
kubectl apply -f demo/gitops/inventory-api-gitops-run.yaml
kubectl apply -f demo/gitops/inventory-dashboard-gitops-run.yaml
```

### Step 5: Monitor Build Progress

```bash
# Watch ComponentWorkflowRun status
kubectl get componentworkflowrun -w

# View Argo Workflow
kubectl get workflow -n openchoreo-ci-default

# View logs
kubectl logs -n openchoreo-ci-default -l workflows.argoproj.io/workflow=inventory-api-gitops-release-001 --all-containers -f
```

### Step 6: Merge PRs and Deploy

Once workflows complete, PRs are created in the GitOps repo. Merge them, then either:

**Option A: Automatic sync with Flux** (recommended)
```bash
# Flux auto-syncs every ~1 minute after PR merge
flux get kustomizations
```

**Option B: Manual apply**
```bash
cd /path/to/warehouse-inventory-gitops
git pull
kubectl apply -f projects/default/components/inventory-api-gitops/releases/
kubectl apply -f projects/default/components/inventory-api-gitops/bindings/
```

### Step 7: Verify Deployment

```bash
# Check pods
kubectl get pods -A | grep gitops

# Check services
kubectl get svc -A | grep gitops

# Test endpoints (release name varies)
curl http://development.openchoreoapis.localhost:19080/inventory-api-gitops-development-<hash>/health
```

## Setting Up Flux (Auto-Sync)

Install Flux to automatically sync the GitOps repo to the cluster:

```bash
# Install Flux
brew install fluxcd/tap/flux
flux install

# Create GitRepository source
flux create source git warehouse-inventory-gitops \
  --url=https://github.com/LakshanSS/warehouse-inventory-gitops \
  --branch=main \
  --interval=1m

# Create Kustomization to sync projects folder
flux create kustomization projects \
  --source=GitRepository/warehouse-inventory-gitops \
  --path="./projects" \
  --prune=true \
  --interval=1m

# Verify
flux get kustomizations
flux get sources git
```

## GitOps Repository Structure

The GitOps repository requires this structure:

```
warehouse-inventory-gitops/
├── platform/                              # Platform resources (optional - may already exist in cluster)
│   ├── component-types/
│   │   ├── service.yaml
│   │   └── web-application.yaml
│   └── infrastructure/
│       ├── environments/
│       │   └── development.yaml
│       └── deployment-pipelines/
│           └── simple-pipeline.yaml
└── projects/
    └── default/
        ├── project.yaml
        ├── kustomization.yaml             # For Flux
        └── components/
            ├── inventory-api-gitops/
            │   ├── component.yaml
            │   ├── workload.yaml          # Generated by workflow
            │   ├── releases/
            │   │   ├── kustomization.yaml
            │   │   └── *.yaml             # Generated by workflow
            │   └── bindings/
            │       └── *.yaml             # Generated by workflow
            └── inventory-dashboard-gitops/
                └── (same structure)
```

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

## Known Issues & Fixes

### 1. Container Name Mismatch

**Issue**: ComponentRelease template uses `workload.containers["app"]` but Workload defines container as `main`.

**Error**: `CEL evaluation error: no such key: app`

**Fix**: Edit the ComponentRelease files in GitOps repo to use `main`:
```yaml
# Change this:
image: ${workload.containers["app"].image}
name: app

# To this:
image: ${workload.containers["main"].image}
name: main
```

### 2. Port Configuration

**Issue**: Service targetPort defaults to 80 but API runs on 9090.

**Fix**: Set port in ComponentRelease `componentProfile.parameters`:
```yaml
spec:
  componentProfile:
    parameters:
      port: 9090  # For API (3000 for Dashboard)
```

### 3. Flux Reconciliation Conflicts

**Issue**: `spec.componentProfile is immutable` error when Flux tries to update existing resources.

**Fix**: Add force annotation in kustomization.yaml:
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - inventory-api-gitops-9dec2481.yaml
commonAnnotations:
  kustomize.toolkit.fluxcd.io/force: "enabled"
```

### 4. OpenAPI Schema File

**Issue**: Workflow fails with `failed to read schema file`.

**Fix**: Remove or comment out `schemaFile` in workload.yaml if the OpenAPI file isn't needed.

## Troubleshooting

| Issue | Check |
|-------|-------|
| Workflow not starting | `kubectl describe componentworkflowrun <name>` |
| Clone fails | Token permissions, repo URL |
| Build fails | Dockerfile path, build context |
| PR not created | `gitops-token` has `repo` scope |
| No deployment after merge | Check ReleaseBinding: `kubectl describe releasebinding <name>` |
| Service port wrong | Check ComponentRelease `componentProfile.parameters.port` |
| Flux not syncing | `flux get kustomizations` and `flux logs` |

```bash
# Check ExternalSecrets
kubectl get externalsecret -n openchoreo-ci-default

# Check workflow pods
kubectl get pods -n openchoreo-ci-default

# Get detailed logs
kubectl logs -n openchoreo-ci-default <pod-name> -c main

# Check ReleaseBinding status
kubectl describe releasebinding inventory-api-gitops-development

# Check Flux status
flux get kustomizations
flux logs --level=error
```

## Related Documentation

- [OpenChoreo GitOps Workflows](../../samples/gitops-workflows/component-workflows/build-and-release/)
- [Docker Workflow README](../../samples/gitops-workflows/component-workflows/build-and-release/docker/README.md)
- [GitOps Repository](https://github.com/LakshanSS/warehouse-inventory-gitops)
