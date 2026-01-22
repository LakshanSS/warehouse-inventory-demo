# Warehouse Inventory Demo for OpenChoreo

A demo application showcasing OpenChoreo deployment with a Go backend API and Next.js frontend.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        OpenChoreo Platform                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────────┐       ┌─────────────────────────────┐  │
│  │  inventory-dashboard│       │      inventory-api          │  │
│  │     (Next.js)       │──────▶│         (Go/Gin)            │  │
│  │      Port: 3000     │       │       Port: 9090            │  │
│  └─────────────────────┘       └──────────────┬──────────────┘  │
│                                               │                  │
│                                               ▼                  │
│                                ┌─────────────────────────────┐  │
│                                │         MSSQL Database      │  │
│                                │      (External/Managed)     │  │
│                                └─────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

## Components

### Backend: inventory-api (Go)
- **Tech Stack**: Go 1.22+, Gin Web Framework, go-mssqldb driver
- **Port**: 9090
- **Endpoints**:
  - `GET /health` - Health check for K8s probes
  - `GET /inventory` - Fetch all stock items
  - `GET /inventory/:id` - Fetch single item
  - `POST /inventory` - Create new item
  - `PUT /inventory/:id` - Update stock quantity
  - `DELETE /inventory/:id` - Remove item

### Frontend: inventory-dashboard (Next.js)
- **Tech Stack**: Next.js 14 (App Router), TanStack Table, Tailwind CSS, Lucide Icons
- **Port**: 3000
- **Features**:
  - Stock overview with search/filter
  - Real-time low-stock indicators (Red/Yellow)
  - Sortable columns
  - Pagination

### Database: MSSQL
- **Schema**: Products table with ID, SKU, Name, Quantity, Category, Location
- See `database/` folder for setup scripts

## Local Development

### 1. Start Database
```bash
cd database
docker-compose up -d

# Wait for MSSQL to be ready, then run setup
docker exec -it inventory-mssql /opt/mssql-tools18/bin/sqlcmd \
  -S localhost -U sa -P "YourStrong@Passw0rd" -C \
  -Q "CREATE DATABASE InventoryDB"

# Run schema and seed scripts using your preferred SQL client
```

### 2. Start Backend
```bash
cd inventory-api
go mod download
go run main.go
```

### 3. Start Frontend
```bash
cd inventory-dashboard
npm install
npm run dev
```

## OpenChoreo Deployment

### Prerequisites
- OpenChoreo cluster running
- GitHub repository with the demo code
- MSSQL database accessible from the cluster

### Configuration Files
- `.openchoreo/inventory-api.yaml` - Backend component config
- `.openchoreo/inventory-dashboard.yaml` - Frontend component config

### Deploy Steps

1. **Update repository URLs** in `.openchoreo/*.yaml` files to point to your GitHub repo

2. **Create database secret**:
```bash
kubectl create secret generic inventory-db-credentials \
  --from-literal=server=your-mssql-server \
  --from-literal=port=1433 \
  --from-literal=username=sa \
  --from-literal=password=YourStrong@Passw0rd \
  --from-literal=database=InventoryDB
```

3. **Apply configurations**:
```bash
kubectl apply -f .openchoreo/inventory-api.yaml
kubectl apply -f .openchoreo/inventory-dashboard.yaml
```

## Environment Variables

### inventory-api
| Variable | Description | Default |
|----------|-------------|---------|
| PORT | API server port | 9090 |
| ENVIRONMENT | Runtime environment | development |
| DB_SERVER | MSSQL server address | localhost |
| DB_PORT | MSSQL server port | 1433 |
| DB_USER | Database username | sa |
| DB_PASSWORD | Database password | - |
| DB_NAME | Database name | InventoryDB |

### inventory-dashboard
| Variable | Description | Default |
|----------|-------------|---------|
| NEXT_PUBLIC_API_URL | Backend API URL | http://localhost:9090 |

## API Documentation

OpenAPI specification available at `inventory-api/docs/openapi.yaml`

## Project Structure

```
demo/
├── inventory-api/           # Go backend
│   ├── api/routes/          # Route handlers
│   ├── internal/
│   │   ├── config/          # Configuration
│   │   ├── controllers/     # Business logic
│   │   ├── models/          # Data models
│   │   └── repositories/    # Database layer
│   ├── docs/                # OpenAPI spec
│   ├── Dockerfile
│   ├── workload.yaml        # OpenChoreo workload descriptor
│   └── main.go
├── inventory-dashboard/     # Next.js frontend
│   ├── src/
│   │   ├── app/             # Next.js App Router
│   │   ├── components/      # React components
│   │   └── lib/             # Utilities & types
│   ├── Dockerfile
│   └── workload.yaml        # OpenChoreo workload descriptor
├── database/                # Database scripts
│   ├── schema.sql           # Table definitions
│   ├── seed.sql             # Sample data
│   └── docker-compose.yml   # Local MSSQL setup
└── .openchoreo/             # OpenChoreo configurations
    ├── inventory-api.yaml
    └── inventory-dashboard.yaml
```
