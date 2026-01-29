# Secrets Management in OpenChoreo

This guide explains how to securely configure database credentials and other secrets in OpenChoreo.

## Overview

OpenChoreo uses the External Secrets Operator pattern:

```
ClusterSecretStore     →  SecretReference  →  Workload  →  Pod
(stores credentials)      (references key)    (uses ref)   (gets env var)
```

**Key principle:** Credentials are stored in ClusterSecretStore, never in Git.

## Quick Setup

### Step 1: Add Credentials to ClusterSecretStore

```bash
kubectl edit clustersecretstore default
```

Add under `spec.provider.fake.data`:

```yaml
spec:
  provider:
    fake:
      data:
        # ... existing entries ...
        - key: inventory-db-server
          value: "your-server.database.windows.net"
        - key: inventory-db-port
          value: "1433"
        - key: inventory-db-user
          value: "your-username"
        - key: inventory-db-password
          value: "your-password"
        - key: inventory-db-name
          value: "your-database"
```

### Step 2: Apply SecretReference Resources

```bash
kubectl apply -f .openchoreo/db-secrets.yaml
```

This creates SecretReference CRs that point to keys in ClusterSecretStore.

### Step 3: Update Workload to Use Secrets

Edit the inventory-api workload:

```bash
kubectl edit workload inventory-api-workload
```

Update the `spec.containers.main.env` section:

```yaml
spec:
  containers:
    main:
      env:
        - key: ENVIRONMENT
          value: production
        - key: DB_SERVER
          valueFrom:
            secretRef:
              name: inventory-db-server
              key: value
        - key: DB_PORT
          valueFrom:
            secretRef:
              name: inventory-db-port
              key: value
        - key: DB_USER
          valueFrom:
            secretRef:
              name: inventory-db-user
              key: value
        - key: DB_PASSWORD
          valueFrom:
            secretRef:
              name: inventory-db-password
              key: value
        - key: DB_NAME
          valueFrom:
            secretRef:
              name: inventory-db-name
              key: value
```

### Step 4: Verify

```bash
# Check SecretReferences exist
kubectl get secretreference

# Check ExternalSecrets are synced
kubectl get externalsecret -A

# Check secrets are created
kubectl get secret -A | grep inventory

# Verify pod environment (after pod restarts)
POD=$(kubectl get pods -A -l openchoreo.dev/component=inventory-api -o jsonpath='{.items[0].metadata.name}')
NS=$(kubectl get pods -A -l openchoreo.dev/component=inventory-api -o jsonpath='{.items[0].metadata.namespace}')
kubectl exec -n $NS $POD -- env | grep DB_
```

## How It Works

### 1. ClusterSecretStore

Central store for all secrets. In development, uses a "fake" provider with inline values. In production, connects to real secret stores (AWS Secrets Manager, Azure Key Vault, GCP Secret Manager, HashiCorp Vault).

### 2. SecretReference (OpenChoreo CR)

Points to a key in ClusterSecretStore. Example from `db-secrets.yaml`:

```yaml
apiVersion: openchoreo.dev/v1alpha1
kind: SecretReference
metadata:
  name: inventory-db-server
spec:
  refreshInterval: 1h
  data:
    - secretKey: value
      remoteRef:
        key: inventory-db-server
```

### 3. ExternalSecret (Auto-created)

OpenChoreo automatically creates ExternalSecret resources that sync to Kubernetes Secrets.

### 4. Workload

Uses `valueFrom.secretRef` to inject secrets as environment variables:

```yaml
env:
  - key: DB_PASSWORD
    valueFrom:
      secretRef:
        name: inventory-db-password
        key: value
```

## Files Reference

| File | Purpose |
|------|---------|
| `.openchoreo/db-secrets.yaml` | SecretReference definitions for DB credentials |

## Adding New Secrets

1. Add key-value to ClusterSecretStore
2. Create a SecretReference CR
3. Reference it in Workload using `valueFrom.secretRef`

## Security Notes

- Never commit credentials to Git
- Use `refreshInterval` for automatic rotation support
- The "fake" provider is for development only
- Production should use real secret stores with proper access controls

## Troubleshooting

### SecretReference not syncing

```bash
kubectl describe secretreference inventory-db-server
kubectl get externalsecret -A
```

### Pod not getting env vars

- Verify the secret exists: `kubectl get secret -A | grep inventory`
- Check pod events: `kubectl describe pod <pod-name>`
- Pod may need restart after secret changes

### Key not found

- Verify the key exists in ClusterSecretStore
- Check spelling matches exactly
