# Warehouse Inventory Demo

A sample application for demonstrating OpenChoreo deployments. Features a Go backend API, Next.js frontend dashboard, and MSSQL database.

## Quick Start (OpenChoreo)

### Prerequisites

- OpenChoreo cluster running
- `kubectl` configured to access the cluster

### Step 1: Set Up Database Secrets

The backend workload requires database secret references to be configured before deployment. The API returns mock data when database connection fails, so you can use placeholder values if you don't have a real database.

**Add credentials to ClusterSecretStore:**

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
        value: placeholder    # Or your real DB server
      - key: inventory-db-port
        value: "1433"
      - key: inventory-db-user
        value: placeholder    # Or your real DB username
      - key: inventory-db-password
        value: placeholder    # Or your real DB password
      - key: inventory-db-name
        value: placeholder    # Or your real DB name
```

**Apply SecretReferences:**

```bash
kubectl apply -f .openchoreo/db-secrets.yaml
```

**Verify secrets are synced:**

```bash
kubectl get secretreference
kubectl get externalsecret -A
```

### Step 2: Deploy Components

```bash
# Deploy backend and frontend
kubectl apply -f .openchoreo/inventory-api.yaml
kubectl apply -f .openchoreo/inventory-dashboard.yaml

# Monitor deployment
kubectl get components -w
```

### Step 3: Access the Demo

| Component | URL |
|-----------|-----|
| **Dashboard** | http://inventory-dashboard-development.openchoreoapis.localhost:19080/ |
| **API** | http://development.openchoreoapis.localhost:19080/inventory-api/inventory |

The backend returns mock data when database connection fails, so you can see the demo working even with placeholder secrets.

## What You'll See

The dashboard displays a warehouse inventory management interface:

- **Inventory Table** - Sortable, filterable list of products
- **Stock Alerts** - Visual indicators for low stock (yellow ≤10, red ≤5)
- **Stats Cards** - Total products, total items, alert counts
- **Search** - Filter products by name, SKU, or category

## Components

| Component | Tech Stack | Port | Description |
|-----------|------------|------|-------------|
| inventory-api | Go, Gin | 9090 | REST API for inventory operations |
| inventory-dashboard | Next.js, Tailwind | 3000 | Web dashboard UI |

## Connect a Real Database (Optional)

To switch from mock data to a real MSSQL database:

1. **Update ClusterSecretStore** with real credentials:
   ```bash
   kubectl edit clustersecretstore default
   ```
   Replace placeholder values with your actual database credentials.

2. **Set up the database schema** - See [docs/DATABASE_SETUP.md](docs/DATABASE_SETUP.md)

3. **Restart the API pod** to pick up the new credentials:
   ```bash
   kubectl delete pod -l openchoreo.dev/component=inventory-api -A
   ```

For more details on secrets management, see [docs/SECRETS_MANAGEMENT.md](docs/SECRETS_MANAGEMENT.md).

## Project Structure

```
├── .openchoreo/              # OpenChoreo component configs
│   ├── inventory-api.yaml
│   ├── inventory-dashboard.yaml
│   └── db-secrets.yaml
├── inventory-api/            # Go backend
├── inventory-dashboard/      # Next.js frontend
├── database/                 # SQL schema and seed scripts
├── infrastructure/           # Self-service infra provisioning
└── docs/                     # Additional documentation
```

## Documentation

| Guide | Description |
|-------|-------------|
| [SECRETS_MANAGEMENT.md](docs/SECRETS_MANAGEMENT.md) | Configure database credentials |
| [DATABASE_SETUP.md](docs/DATABASE_SETUP.md) | Set up and seed the database |
| [LOCAL_DEVELOPMENT.md](docs/LOCAL_DEVELOPMENT.md) | Run locally without OpenChoreo |
| [ARCHITECTURE.md](docs/ARCHITECTURE.md) | Technical architecture details |
| [infrastructure/README.md](infrastructure/README.md) | Self-service database provisioning |

## Useful Commands

```bash
# Check deployment status
kubectl get components
kubectl get workloads
kubectl get pods -A | grep inventory

# View logs
kubectl logs -n <namespace> <pod-name>

# Check secrets
kubectl get secretreference
kubectl get externalsecret -A
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/inventory` | List all items |
| GET | `/inventory/:id` | Get single item |
| POST | `/inventory` | Create item |
| PUT | `/inventory/:id` | Update item |
| DELETE | `/inventory/:id` | Delete item |

Full OpenAPI spec: [inventory-api/docs/openapi.yaml](inventory-api/docs/openapi.yaml)
