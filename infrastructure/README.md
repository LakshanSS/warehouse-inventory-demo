# Self-Service Infrastructure Provisioning

Provision Azure SQL databases using OpenChoreo Generic Workflows. Developers can create databases without direct cloud console access.

## Quick Start

### Prerequisites

- OpenChoreo cluster with Argo Workflows
- Azure credentials configured (see [Setup Azure Credentials](#setup-azure-credentials))

### Install Workflow

```bash
kubectl apply -f workflows/cluster-workflow-template-terraform-azure-sql.yaml
kubectl apply -f workflows/workflow-terraform-azure-sql.yaml
```

### Provision a Database

```bash
kubectl apply -f workflows/workflow-run-provision-inventory-db.yaml
```

### Monitor Progress

```bash
kubectl get workflowrun -w
```

## How It Works

OpenChoreo Generic Workflows use three resources:

```
WorkflowRun  →  Workflow  →  ClusterWorkflowTemplate
(trigger)       (schema)     (Argo execution steps)
```

1. **ClusterWorkflowTemplate** - Defines the Terraform execution steps (init, plan, apply)
2. **Workflow** - Defines the parameter schema (database name, SKU, region)
3. **WorkflowRun** - Triggers execution with specific values

## Files

| File | Purpose |
|------|---------|
| `workflows/cluster-workflow-template-terraform-azure-sql.yaml` | Argo steps for Terraform |
| `workflows/workflow-terraform-azure-sql.yaml` | Parameter schema |
| `workflows/workflow-run-provision-inventory-db.yaml` | Example trigger |
| `terraform/azure-sql/main.tf` | Terraform configuration |

## Setup Azure Credentials

Create a secret with Azure service principal credentials:

```bash
kubectl create secret generic azure-credentials \
  --from-literal=subscription-id=YOUR_SUBSCRIPTION_ID \
  --from-literal=tenant-id=YOUR_TENANT_ID \
  --from-literal=client-id=YOUR_CLIENT_ID \
  --from-literal=client-secret=YOUR_CLIENT_SECRET \
  --from-literal=db-admin-username=sqladmin \
  --from-literal=db-admin-password=YOUR_SECURE_PASSWORD
```

## Configuration Options

Edit `workflow-run-provision-inventory-db.yaml` to customize:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `database.serverName` | SQL Server name (globally unique) | - |
| `database.name` | Database name | - |
| `database.sku` | Pricing tier | `Basic` |
| `azure.resourceGroup` | Azure resource group | - |
| `azure.location` | Azure region | `eastus` |

Available SKUs: `Basic`, `Standard_S0`, `Standard_S1`, `Standard_S2`, `Premium_P1`

## Verify Provisioning

After workflow completes:

```bash
# Check the connection secret was created
kubectl get secret pge-demo-db-connection -o yaml

# Get server FQDN
kubectl get secret pge-demo-db-connection -o jsonpath='{.data.server}' | base64 -d
```

## Troubleshooting

### Workflow stuck in pending
```bash
kubectl get clusterworkflowtemplate terraform-azure-sql
kubectl describe workflow provision-azure-sql-database
```

### Terraform fails
```bash
# View logs
kubectl logs -n openchoreo-ci-default \
  -l workflows.argoproj.io/workflow=provision-inventory-db-001 \
  --all-containers=true
```

Common issues:
- Azure credentials missing or invalid
- Server name already taken (must be globally unique)
- Resource group doesn't exist
- Insufficient Azure permissions

### Cleanup

```bash
kubectl delete workflowrun provision-inventory-db-001
```

## Extending

The same pattern works for other infrastructure:

| Resource | Provider |
|----------|----------|
| Azure Storage | azurerm |
| Azure Redis | azurerm |
| AWS RDS | aws |
| AWS S3 | aws |
| GCP Cloud SQL | google |

Create new Terraform configs in `terraform/` and update the workflow templates.
