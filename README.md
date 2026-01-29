# Warehouse Inventory Demo

A sample application for demonstrating OpenChoreo deployments. Features a Go backend API, Next.js frontend dashboard, and MSSQL database.

## Quick Start (OpenChoreo)

### Prerequisites

- OpenChoreo cluster running
- `kubectl` configured to access the cluster

### Deploy

```bash
# Deploy backend and frontend
kubectl apply -f .openchoreo/inventory-api.yaml
kubectl apply -f .openchoreo/inventory-dashboard.yaml

# Monitor deployment
kubectl get components -w
```

### Access the Demo

| Component | URL |
|-----------|-----|
| **Dashboard** | http://inventory-dashboard-development.openchoreoapis.localhost:19080/ |
| **API** | http://development.openchoreoapis.localhost:19080/inventory-api/inventory |

The backend returns mock data by default, so you can see the demo immediately without database setup.

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

To connect to a real MSSQL database instead of mock data:

1. **Add credentials to ClusterSecretStore** - See [docs/SECRETS_MANAGEMENT.md](docs/SECRETS_MANAGEMENT.md)
2. **Apply SecretReferences**:
   ```bash
   kubectl apply -f .openchoreo/db-secrets.yaml
   ```
3. **Update the workload** to use secret references
4. **Set up the database** - See [docs/DATABASE_SETUP.md](docs/DATABASE_SETUP.md)

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
