# Self-Service Infrastructure Workflows in OpenChoreo

This document explains how to create self-service infrastructure pipelines using OpenChoreo's Generic Workflows. Developers can trigger workflows to provision infrastructure (like databases) without needing direct access to cloud consoles or Terraform.

---

## Overview

OpenChoreo provides **Generic Workflows** - standalone automation pipelines not tied to application builds. These are ideal for:

- **Infrastructure Provisioning** - Create databases, storage, networking via Terraform/Pulumi
- **Data Pipelines** - ETL jobs, data migrations
- **Testing** - Integration tests, load tests
- **Package Publishing** - Publish to npm, PyPI, Maven
- **Operations** - Backups, cleanup jobs, certificate rotation

```
┌─────────────────────────────────────────────────────────────────┐
│                    Self-Service Infrastructure                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   Developer                                                      │
│      │                                                           │
│      │ kubectl apply -f provision-database-run.yaml              │
│      ▼                                                           │
│   ┌──────────────┐                                               │
│   │ WorkflowRun  │  (Trigger with parameters)                    │
│   └──────┬───────┘                                               │
│          │                                                       │
│          ▼                                                       │
│   ┌──────────────┐     ┌─────────────────────────┐               │
│   │   Workflow   │────▶│ ClusterWorkflowTemplate │               │
│   │   (Schema)   │     │   (Argo Workflows)      │               │
│   └──────────────┘     └───────────┬─────────────┘               │
│                                    │                             │
│                                    ▼                             │
│                        ┌─────────────────────────┐               │
│                        │   Terraform Container   │               │
│                        │   - terraform init      │               │
│                        │   - terraform plan      │               │
│                        │   - terraform apply     │               │
│                        └───────────┬─────────────┘               │
│                                    │                             │
│                                    ▼                             │
│                        ┌─────────────────────────┐               │
│                        │   Azure SQL Database    │               │
│                        │   (Provisioned)         │               │
│                        └─────────────────────────┘               │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Generic Workflows vs Component Workflows

| Aspect | Generic Workflows | Component Workflows |
|--------|-------------------|---------------------|
| **Purpose** | Any automation task | Application builds |
| **Tied to Component** | No | Yes |
| **Trigger** | Manual WorkflowRun CR | Auto or ComponentWorkflowRun |
| **Schema** | `parameters` only | `systemParameters` + `parameters` |
| **Use Cases** | Infrastructure, ETL, testing | Build from source, GitOps release |
| **CR Types** | Workflow, WorkflowRun | ComponentWorkflow, ComponentWorkflowRun |

---

## The 3-Resource Pattern

Generic workflows use three Kubernetes resources:

### 1. ClusterWorkflowTemplate (Argo Workflows)

Defines the **execution logic** - what containers to run, in what order, with what commands.

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ClusterWorkflowTemplate
metadata:
  name: terraform-azure-database
spec:
  entrypoint: provision-database
  templates:
    - name: provision-database
      steps:
        - - name: terraform-init
            template: init
        - - name: terraform-plan
            template: plan
        - - name: terraform-apply
            template: apply

    - name: init
      container:
        image: hashicorp/terraform:1.7
        command: [sh, -c]
        args:
          - |
            cd /workspace/terraform
            terraform init

    - name: plan
      container:
        image: hashicorp/terraform:1.7
        command: [sh, -c]
        args:
          - |
            terraform plan -var="db_name={{workflow.parameters.db-name}}"

    - name: apply
      container:
        image: hashicorp/terraform:1.7
        command: [sh, -c]
        args:
          - |
            terraform apply -auto-approve
```

### 2. Workflow (OpenChoreo CR)

Defines the **parameter schema** and maps parameters to the ClusterWorkflowTemplate.

```yaml
apiVersion: openchoreo.dev/v1alpha1
kind: Workflow
metadata:
  name: provision-azure-database
  namespace: default
spec:
  schema:
    parameters:
      database:
        name: string | description="Database name"
        sku: string | default="Basic" | description="Database SKU"
      environment:
        name: string | description="Target environment"

  runTemplate:
    apiVersion: argoproj.io/v1alpha1
    kind: Workflow
    spec:
      arguments:
        parameters:
          - name: db-name
            value: ${parameters.database.name}
          - name: db-sku
            value: ${parameters.database.sku}
      workflowTemplateRef:
        clusterScope: true
        name: terraform-azure-database
```

### 3. WorkflowRun (OpenChoreo CR)

**Triggers** the workflow with specific parameter values.

```yaml
apiVersion: openchoreo.dev/v1alpha1
kind: WorkflowRun
metadata:
  name: provision-inventory-db-001
spec:
  workflow:
    name: provision-azure-database
    parameters:
      database:
        name: "inventory-db"
        sku: "Standard_S0"
      environment:
        name: "development"
```

---

## Demo: Azure SQL Database Provisioning

For the warehouse inventory demo, we can create a self-service workflow to provision Azure SQL databases.

### Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     Database Provisioning Flow                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. Developer creates WorkflowRun                                │
│     ┌─────────────────────────────────────────────────────┐     │
│     │ parameters:                                          │     │
│     │   database:                                          │     │
│     │     name: "pge-demo-db"                              │     │
│     │     sku: "Standard_S0"                               │     │
│     │   server:                                            │     │
│     │     name: "pge-demo-db-server"                       │     │
│     │     region: "eastus"                                 │     │
│     └─────────────────────────────────────────────────────┘     │
│                          │                                       │
│                          ▼                                       │
│  2. Argo Workflow Executes                                       │
│     ┌─────────────────────────────────────────────────────┐     │
│     │ Step 1: Clone Terraform configs from Git             │     │
│     │ Step 2: terraform init (configure backend)           │     │
│     │ Step 3: terraform plan (show changes)                │     │
│     │ Step 4: terraform apply (create resources)           │     │
│     │ Step 5: Output connection string to Secret           │     │
│     └─────────────────────────────────────────────────────┘     │
│                          │                                       │
│                          ▼                                       │
│  3. Azure Resources Created                                      │
│     ┌─────────────────────────────────────────────────────┐     │
│     │ - Azure SQL Server: pge-demo-db-server               │     │
│     │ - Azure SQL Database: pge-demo-db                    │     │
│     │ - Firewall rules configured                          │     │
│     │ - Connection string saved to K8s Secret              │     │
│     └─────────────────────────────────────────────────────┘     │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Implementation Files

### File 1: Terraform Configuration

Create `demo/infrastructure/terraform/azure-sql/main.tf`:

```hcl
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
  }

  backend "azurerm" {
    resource_group_name  = "terraform-state-rg"
    storage_account_name = "tfstatestore"
    container_name       = "tfstate"
    key                  = "database.tfstate"
  }
}

provider "azurerm" {
  features {}
}

variable "resource_group_name" {
  description = "Resource group name"
  type        = string
}

variable "location" {
  description = "Azure region"
  type        = string
  default     = "eastus"
}

variable "server_name" {
  description = "SQL Server name"
  type        = string
}

variable "database_name" {
  description = "Database name"
  type        = string
}

variable "database_sku" {
  description = "Database SKU"
  type        = string
  default     = "Basic"
}

variable "admin_username" {
  description = "Admin username"
  type        = string
  sensitive   = true
}

variable "admin_password" {
  description = "Admin password"
  type        = string
  sensitive   = true
}

resource "azurerm_mssql_server" "main" {
  name                         = var.server_name
  resource_group_name          = var.resource_group_name
  location                     = var.location
  version                      = "12.0"
  administrator_login          = var.admin_username
  administrator_login_password = var.admin_password

  tags = {
    environment = "demo"
    managed_by  = "openchoreo"
  }
}

resource "azurerm_mssql_database" "main" {
  name         = var.database_name
  server_id    = azurerm_mssql_server.main.id
  sku_name     = var.database_sku
  collation    = "SQL_Latin1_General_CP1_CI_AS"
  max_size_gb  = 2

  tags = {
    environment = "demo"
    managed_by  = "openchoreo"
  }
}

resource "azurerm_mssql_firewall_rule" "allow_azure" {
  name             = "AllowAzureServices"
  server_id        = azurerm_mssql_server.main.id
  start_ip_address = "0.0.0.0"
  end_ip_address   = "0.0.0.0"
}

output "server_fqdn" {
  value = azurerm_mssql_server.main.fully_qualified_domain_name
}

output "database_name" {
  value = azurerm_mssql_database.main.name
}

output "connection_string" {
  value     = "Server=tcp:${azurerm_mssql_server.main.fully_qualified_domain_name},1433;Database=${azurerm_mssql_database.main.name};User ID=${var.admin_username};Password=${var.admin_password};Encrypt=true;TrustServerCertificate=false;"
  sensitive = true
}
```

### File 2: ClusterWorkflowTemplate

Create `demo/infrastructure/workflows/cluster-workflow-template-terraform-azure-sql.yaml`:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ClusterWorkflowTemplate
metadata:
  name: terraform-azure-sql
spec:
  entrypoint: provision-database

  arguments:
    parameters:
      - name: git-repo
        value: ""
      - name: git-branch
        value: "main"
      - name: terraform-path
        value: "/infrastructure/terraform/azure-sql"
      - name: resource-group
        value: ""
      - name: location
        value: "eastus"
      - name: server-name
        value: ""
      - name: database-name
        value: ""
      - name: database-sku
        value: "Basic"

  templates:
    # Main orchestration template
    - name: provision-database
      steps:
        - - name: clone-repo
            template: clone
        - - name: terraform-init
            template: init
        - - name: terraform-plan
            template: plan
        - - name: terraform-apply
            template: apply
            arguments:
              parameters:
                - name: plan-output
                  value: "{{steps.terraform-plan.outputs.parameters.plan-file}}"
        - - name: save-outputs
            template: save-to-secret

    # Step 1: Clone repository with Terraform configs
    - name: clone
      container:
        image: alpine/git:2.43.0
        command:
          - sh
          - -c
          - |
            set -e
            echo "Cloning {{workflow.parameters.git-repo}}..."
            git clone --depth 1 --branch {{workflow.parameters.git-branch}} \
              {{workflow.parameters.git-repo}} /mnt/vol/source
            echo "Clone complete"
            ls -la /mnt/vol/source{{workflow.parameters.terraform-path}}
        volumeMounts:
          - name: workspace
            mountPath: /mnt/vol

    # Step 2: Terraform init
    - name: init
      container:
        image: hashicorp/terraform:1.7
        workingDir: /mnt/vol/source{{workflow.parameters.terraform-path}}
        command:
          - sh
          - -c
          - |
            set -e
            echo "Initializing Terraform..."
            terraform init \
              -backend-config="resource_group_name={{workflow.parameters.resource-group}}" \
              -backend-config="key={{workflow.parameters.database-name}}.tfstate"
            echo "Init complete"
        env:
          - name: ARM_SUBSCRIPTION_ID
            valueFrom:
              secretKeyRef:
                name: azure-credentials
                key: subscription-id
          - name: ARM_TENANT_ID
            valueFrom:
              secretKeyRef:
                name: azure-credentials
                key: tenant-id
          - name: ARM_CLIENT_ID
            valueFrom:
              secretKeyRef:
                name: azure-credentials
                key: client-id
          - name: ARM_CLIENT_SECRET
            valueFrom:
              secretKeyRef:
                name: azure-credentials
                key: client-secret
        volumeMounts:
          - name: workspace
            mountPath: /mnt/vol

    # Step 3: Terraform plan
    - name: plan
      container:
        image: hashicorp/terraform:1.7
        workingDir: /mnt/vol/source{{workflow.parameters.terraform-path}}
        command:
          - sh
          - -c
          - |
            set -e
            echo "Planning Terraform changes..."
            terraform plan \
              -var="resource_group_name={{workflow.parameters.resource-group}}" \
              -var="location={{workflow.parameters.location}}" \
              -var="server_name={{workflow.parameters.server-name}}" \
              -var="database_name={{workflow.parameters.database-name}}" \
              -var="database_sku={{workflow.parameters.database-sku}}" \
              -var="admin_username=$DB_ADMIN_USER" \
              -var="admin_password=$DB_ADMIN_PASS" \
              -out=/mnt/vol/tfplan
            echo "Plan complete"
            echo "/mnt/vol/tfplan" > /mnt/vol/plan-path.txt
        env:
          - name: ARM_SUBSCRIPTION_ID
            valueFrom:
              secretKeyRef:
                name: azure-credentials
                key: subscription-id
          - name: ARM_TENANT_ID
            valueFrom:
              secretKeyRef:
                name: azure-credentials
                key: tenant-id
          - name: ARM_CLIENT_ID
            valueFrom:
              secretKeyRef:
                name: azure-credentials
                key: client-id
          - name: ARM_CLIENT_SECRET
            valueFrom:
              secretKeyRef:
                name: azure-credentials
                key: client-secret
          - name: DB_ADMIN_USER
            valueFrom:
              secretKeyRef:
                name: azure-credentials
                key: db-admin-username
          - name: DB_ADMIN_PASS
            valueFrom:
              secretKeyRef:
                name: azure-credentials
                key: db-admin-password
        volumeMounts:
          - name: workspace
            mountPath: /mnt/vol
      outputs:
        parameters:
          - name: plan-file
            valueFrom:
              path: /mnt/vol/plan-path.txt

    # Step 4: Terraform apply
    - name: apply
      inputs:
        parameters:
          - name: plan-output
      container:
        image: hashicorp/terraform:1.7
        workingDir: /mnt/vol/source{{workflow.parameters.terraform-path}}
        command:
          - sh
          - -c
          - |
            set -e
            echo "Applying Terraform plan..."
            terraform apply -auto-approve {{inputs.parameters.plan-output}}

            echo "Extracting outputs..."
            terraform output -json > /mnt/vol/tf-outputs.json
            terraform output -raw server_fqdn > /mnt/vol/server-fqdn.txt
            terraform output -raw database_name > /mnt/vol/database-name.txt

            echo "Apply complete!"
            echo "Server: $(cat /mnt/vol/server-fqdn.txt)"
            echo "Database: $(cat /mnt/vol/database-name.txt)"
        env:
          - name: ARM_SUBSCRIPTION_ID
            valueFrom:
              secretKeyRef:
                name: azure-credentials
                key: subscription-id
          - name: ARM_TENANT_ID
            valueFrom:
              secretKeyRef:
                name: azure-credentials
                key: tenant-id
          - name: ARM_CLIENT_ID
            valueFrom:
              secretKeyRef:
                name: azure-credentials
                key: client-id
          - name: ARM_CLIENT_SECRET
            valueFrom:
              secretKeyRef:
                name: azure-credentials
                key: client-secret
        volumeMounts:
          - name: workspace
            mountPath: /mnt/vol

    # Step 5: Save outputs to Kubernetes Secret
    - name: save-to-secret
      container:
        image: bitnami/kubectl:1.28
        command:
          - sh
          - -c
          - |
            set -e
            echo "Saving database connection info to Secret..."

            SERVER_FQDN=$(cat /mnt/vol/server-fqdn.txt)
            DB_NAME=$(cat /mnt/vol/database-name.txt)

            kubectl create secret generic {{workflow.parameters.database-name}}-connection \
              --from-literal=server=$SERVER_FQDN \
              --from-literal=database=$DB_NAME \
              --from-literal=port=1433 \
              --dry-run=client -o yaml | kubectl apply -f -

            echo "Secret created: {{workflow.parameters.database-name}}-connection"
        volumeMounts:
          - name: workspace
            mountPath: /mnt/vol

  # Persistent volume for workflow
  volumeClaimTemplates:
    - metadata:
        name: workspace
      spec:
        accessModes:
          - ReadWriteOnce
        resources:
          requests:
            storage: 1Gi

  # Auto-cleanup after 2 hours
  ttlStrategy:
    secondsAfterCompletion: 7200
    secondsAfterSuccess: 7200
    secondsAfterFailure: 86400  # Keep failed runs for 24h for debugging
```

### File 3: Workflow (OpenChoreo CR)

Create `demo/infrastructure/workflows/workflow-terraform-azure-sql.yaml`:

```yaml
# OpenChoreo Workflow for Azure SQL Database Provisioning
# This defines the parameter schema and maps to the ClusterWorkflowTemplate

apiVersion: openchoreo.dev/v1alpha1
kind: Workflow
metadata:
  name: provision-azure-sql-database
  namespace: default
  labels:
    category: infrastructure
    provider: azure
    resource-type: database
spec:
  # Parameter schema - defines what users can configure
  schema:
    parameters:
      # Source repository containing Terraform configs
      repository:
        url:
          type: string
          description: "Git repository URL containing Terraform configurations"
          default: "https://github.com/LakshanSS/warehouse-inventory-demo"
        branch:
          type: string
          description: "Git branch"
          default: "main"
        terraformPath:
          type: string
          description: "Path to Terraform configs within repo"
          default: "/demo/infrastructure/terraform/azure-sql"

      # Azure resource configuration
      azure:
        resourceGroup:
          type: string
          description: "Azure Resource Group name"
        location:
          type: string
          description: "Azure region"
          default: "eastus"
          enum:
            - eastus
            - westus2
            - centralus
            - northeurope
            - westeurope

      # Database configuration
      database:
        serverName:
          type: string
          description: "SQL Server name (must be globally unique)"
        name:
          type: string
          description: "Database name"
        sku:
          type: string
          description: "Database SKU/tier"
          default: "Basic"
          enum:
            - Basic
            - Standard_S0
            - Standard_S1
            - Standard_S2
            - Premium_P1

  # Template that generates the Argo Workflow
  runTemplate:
    apiVersion: argoproj.io/v1alpha1
    kind: Workflow
    metadata:
      name: ${metadata.workflowRunName}
      namespace: openchoreo-ci-${metadata.orgName}
      labels:
        openchoreo.dev/workflow: provision-azure-sql-database
        openchoreo.dev/workflow-run: ${metadata.workflowRunName}
    spec:
      serviceAccountName: workflow-sa
      arguments:
        parameters:
          - name: git-repo
            value: ${parameters.repository.url}
          - name: git-branch
            value: ${parameters.repository.branch}
          - name: terraform-path
            value: ${parameters.repository.terraformPath}
          - name: resource-group
            value: ${parameters.azure.resourceGroup}
          - name: location
            value: ${parameters.azure.location}
          - name: server-name
            value: ${parameters.database.serverName}
          - name: database-name
            value: ${parameters.database.name}
          - name: database-sku
            value: ${parameters.database.sku}
      workflowTemplateRef:
        clusterScope: true
        name: terraform-azure-sql
```

### File 4: WorkflowRun Example

Create `demo/infrastructure/workflows/workflow-run-provision-inventory-db.yaml`:

```yaml
# WorkflowRun to provision the inventory demo database
# Usage: kubectl apply -f workflow-run-provision-inventory-db.yaml

apiVersion: openchoreo.dev/v1alpha1
kind: WorkflowRun
metadata:
  name: provision-inventory-db-001
spec:
  workflow:
    name: provision-azure-sql-database
    parameters:
      repository:
        url: "https://github.com/LakshanSS/warehouse-inventory-demo"
        branch: "main"
        terraformPath: "/demo/infrastructure/terraform/azure-sql"
      azure:
        resourceGroup: "openchoreo-demo-rg"
        location: "eastus"
      database:
        serverName: "pge-demo-db-server"
        name: "pge-demo-db"
        sku: "Standard_S0"
```

---

## Setup Steps

### Step 1: Add Azure Credentials to ClusterSecretStore

```bash
kubectl patch clustersecretstore default --type='json' -p='[
  {
    "op": "add",
    "path": "/spec/provider/fake/data/-",
    "value": {
      "key": "azure-subscription-id",
      "value": "YOUR_SUBSCRIPTION_ID"
    }
  },
  {
    "op": "add",
    "path": "/spec/provider/fake/data/-",
    "value": {
      "key": "azure-tenant-id",
      "value": "YOUR_TENANT_ID"
    }
  },
  {
    "op": "add",
    "path": "/spec/provider/fake/data/-",
    "value": {
      "key": "azure-client-id",
      "value": "YOUR_CLIENT_ID"
    }
  },
  {
    "op": "add",
    "path": "/spec/provider/fake/data/-",
    "value": {
      "key": "azure-client-secret",
      "value": "YOUR_CLIENT_SECRET"
    }
  },
  {
    "op": "add",
    "path": "/spec/provider/fake/data/-",
    "value": {
      "key": "db-admin-username",
      "value": "sqladmin"
    }
  },
  {
    "op": "add",
    "path": "/spec/provider/fake/data/-",
    "value": {
      "key": "db-admin-password",
      "value": "YOUR_SECURE_PASSWORD"
    }
  }
]'
```

### Step 2: Create Azure Credentials Secret

```bash
kubectl create secret generic azure-credentials \
  --from-literal=subscription-id=YOUR_SUBSCRIPTION_ID \
  --from-literal=tenant-id=YOUR_TENANT_ID \
  --from-literal=client-id=YOUR_CLIENT_ID \
  --from-literal=client-secret=YOUR_CLIENT_SECRET \
  --from-literal=db-admin-username=sqladmin \
  --from-literal=db-admin-password=YOUR_SECURE_PASSWORD
```

### Step 3: Install the Workflow Resources

```bash
# Install ClusterWorkflowTemplate
kubectl apply -f demo/infrastructure/workflows/cluster-workflow-template-terraform-azure-sql.yaml

# Install OpenChoreo Workflow
kubectl apply -f demo/infrastructure/workflows/workflow-terraform-azure-sql.yaml

# Verify
kubectl get clusterworkflowtemplate terraform-azure-sql
kubectl get workflow provision-azure-sql-database
```

### Step 4: Trigger Database Provisioning

```bash
# Apply the WorkflowRun
kubectl apply -f demo/infrastructure/workflows/workflow-run-provision-inventory-db.yaml

# Monitor progress
kubectl get workflowrun -w

# View Argo workflow
kubectl get workflow -n openchoreo-ci-default

# Check logs
kubectl logs -n openchoreo-ci-default \
  -l workflows.argoproj.io/workflow=provision-inventory-db-001 \
  --all-containers=true -f
```

### Step 5: Verify Database Created

```bash
# Check the connection secret was created
kubectl get secret pge-demo-db-connection -o yaml

# Use the connection info in your application
kubectl get secret pge-demo-db-connection -o jsonpath='{.data.server}' | base64 -d
```

---

## Workflow Monitoring

### Check WorkflowRun Status

```bash
kubectl get workflowrun
kubectl describe workflowrun provision-inventory-db-001
```

### View Argo Workflow Steps

```bash
# List workflows in build plane
kubectl get workflow -n openchoreo-ci-default

# Describe specific workflow
kubectl describe workflow provision-inventory-db-001 -n openchoreo-ci-default

# View step-by-step progress
kubectl get pods -n openchoreo-ci-default -l workflows.argoproj.io/workflow=provision-inventory-db-001
```

### Check Logs for Each Step

```bash
# All containers
kubectl logs -n openchoreo-ci-default \
  -l workflows.argoproj.io/workflow=provision-inventory-db-001 \
  --all-containers=true

# Specific step (e.g., terraform-apply)
kubectl logs -n openchoreo-ci-default \
  provision-inventory-db-001-terraform-apply-XXXX -c main
```

---

## Extending the Pattern

### Other Infrastructure Resources

The same pattern can provision:

| Resource | Terraform Provider | Use Case |
|----------|-------------------|----------|
| Azure Storage | azurerm | Blob storage, file shares |
| Azure Redis | azurerm | Caching layer |
| Azure Service Bus | azurerm | Message queues |
| AWS RDS | aws | PostgreSQL, MySQL |
| AWS S3 | aws | Object storage |
| GCP Cloud SQL | google | Managed databases |
| GCP GCS | google | Object storage |

### Multi-Environment Support

Create environment-specific WorkflowRuns:

```yaml
# Development
apiVersion: openchoreo.dev/v1alpha1
kind: WorkflowRun
metadata:
  name: provision-db-dev-001
spec:
  workflow:
    name: provision-azure-sql-database
    parameters:
      azure:
        resourceGroup: "dev-rg"
      database:
        name: "inventory-db-dev"
        sku: "Basic"

---
# Production
apiVersion: openchoreo.dev/v1alpha1
kind: WorkflowRun
metadata:
  name: provision-db-prod-001
spec:
  workflow:
    name: provision-azure-sql-database
    parameters:
      azure:
        resourceGroup: "prod-rg"
      database:
        name: "inventory-db-prod"
        sku: "Premium_P1"
```

---

## Troubleshooting

### Workflow Stuck in Pending

```bash
# Check if ClusterWorkflowTemplate exists
kubectl get clusterworkflowtemplate terraform-azure-sql

# Check Workflow CR
kubectl describe workflow provision-azure-sql-database

# Check events
kubectl get events --sort-by='.lastTimestamp'
```

### Terraform Init Fails

- Verify Azure credentials secret exists
- Check backend storage account is accessible
- Ensure service principal has required permissions

### Terraform Apply Fails

```bash
# Check terraform plan output
kubectl logs -n openchoreo-ci-default <plan-pod-name> -c main

# Common issues:
# - Resource group doesn't exist
# - Server name already taken (must be globally unique)
# - Insufficient permissions
```

### Secret Not Created

- Verify the workflow-sa service account has permission to create secrets
- Check the save-to-secret step logs

---

## Security Considerations

1. **Credential Rotation**: Rotate Azure service principal credentials regularly
2. **Least Privilege**: Grant only required permissions to the service principal
3. **State Encryption**: Enable encryption for Terraform state in Azure Storage
4. **Network Security**: Configure firewall rules appropriately
5. **Audit Logging**: Enable Azure Activity Log for compliance

---

## Directory Structure

After implementation:

```
demo/
├── infrastructure/
│   ├── terraform/
│   │   └── azure-sql/
│   │       ├── main.tf
│   │       ├── variables.tf
│   │       └── outputs.tf
│   └── workflows/
│       ├── cluster-workflow-template-terraform-azure-sql.yaml
│       ├── workflow-terraform-azure-sql.yaml
│       └── workflow-run-provision-inventory-db.yaml
├── database/           # Existing SQL scripts
├── inventory-api/      # Go backend
├── inventory-dashboard/ # Next.js frontend
└── .openchoreo/        # Component configs
```

---

## Quick Reference

```bash
# Install workflow
kubectl apply -f demo/infrastructure/workflows/cluster-workflow-template-terraform-azure-sql.yaml
kubectl apply -f demo/infrastructure/workflows/workflow-terraform-azure-sql.yaml

# Provision database
kubectl apply -f demo/infrastructure/workflows/workflow-run-provision-inventory-db.yaml

# Monitor
kubectl get workflowrun -w
kubectl logs -n openchoreo-ci-default -l workflows.argoproj.io/workflow=<name> -f

# Cleanup (re-run with different name)
kubectl delete workflowrun provision-inventory-db-001
```
