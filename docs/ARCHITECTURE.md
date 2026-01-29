# Architecture Overview

Technical details of the Warehouse Inventory Demo architecture.

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        OpenChoreo Platform                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────────┐       ┌─────────────────────────────┐  │
│  │  inventory-dashboard│       │      inventory-api          │  │
│  │     (Next.js 14)    │──────▶│         (Go/Gin)            │  │
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

### Backend: inventory-api

| Aspect | Details |
|--------|---------|
| **Language** | Go 1.22+ |
| **Framework** | Gin Web Framework |
| **Database Driver** | go-mssqldb |
| **Port** | 9090 |

**API Endpoints:**

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check (returns DB status) |
| GET | `/inventory` | List all inventory items |
| GET | `/inventory/:id` | Get single item |
| POST | `/inventory` | Create new item |
| PUT | `/inventory/:id` | Update item quantity |
| DELETE | `/inventory/:id` | Delete item |

**Key Features:**
- Returns mock data (15 items) when database unavailable
- Graceful shutdown handling
- CORS enabled for frontend access
- Health/readiness probes for Kubernetes

**Code Structure:**
```
inventory-api/
├── main.go                 # Entry point, graceful shutdown
├── api/routes/routes.go    # Route registration
└── internal/
    ├── config/config.go    # Environment configuration
    ├── controllers/        # HTTP handlers + mock data
    ├── models/product.go   # Data model
    └── repositories/       # Database operations
```

### Frontend: inventory-dashboard

| Aspect | Details |
|--------|---------|
| **Framework** | Next.js 14 (App Router) |
| **UI Components** | TanStack Table, Tailwind CSS |
| **Icons** | Lucide React |
| **Port** | 3000 |

**Features:**
- Stock overview table with search/filter/sort
- Pagination (10 items per page)
- Low stock alerts (Yellow: ≤10, Red: ≤5)
- Stats cards (total products, total items, alerts)

**Dynamic API URL Detection:**

The frontend automatically detects the correct API URL:
```javascript
function getApiUrl() {
  if (typeof window === 'undefined') return 'http://localhost:9090';

  const hostname = window.location.hostname;
  // OpenChoreo pattern: component-env.openchoreoapis.localhost
  if (hostname.includes('openchoreoapis')) {
    const env = hostname.split('-').pop().split('.')[0];
    return `http://${env}.openchoreoapis.localhost:19080/inventory-api`;
  }
  return 'http://localhost:9090';
}
```

**Code Structure:**
```
inventory-dashboard/
├── src/
│   ├── app/
│   │   ├── page.tsx        # Main dashboard + API detection
│   │   ├── layout.tsx      # Root layout
│   │   └── globals.css     # Global styles
│   ├── components/
│   │   ├── InventoryTable.tsx  # TanStack Table
│   │   ├── StockBadge.tsx      # Quantity indicators
│   │   └── StockIndicator.tsx  # Alert banners
│   └── lib/
│       └── types.ts        # TypeScript interfaces
└── public/                 # Static assets
```

### Database

| Aspect | Details |
|--------|---------|
| **Type** | Microsoft SQL Server |
| **Table** | Products |
| **Indexes** | SKU, Category, Quantity |

See [DATABASE_SETUP.md](DATABASE_SETUP.md) for schema details.

## OpenChoreo Configuration

### Component Types

| Component | Type | Description |
|-----------|------|-------------|
| inventory-api | `deployment/service` | Backend service with internal routing |
| inventory-dashboard | `deployment/web-application` | Frontend with external access |

### Resource Allocation

Both components:
- **CPU:** 100m request, 500m limit
- **Memory:** 256Mi request, 512Mi limit

### Health Probes

| Component | Path | Port |
|-----------|------|------|
| inventory-api | `/health` | 9090 |
| inventory-dashboard | `/` | 3000 |

## URL Patterns (OpenChoreo)

| Type | Pattern | Example |
|------|---------|---------|
| Backend API | `http://{env}.openchoreoapis.localhost:{port}/{component}/...` | `http://development.openchoreoapis.localhost:19080/inventory-api/inventory` |
| Frontend | `http://{component}-{env}.openchoreoapis.localhost:{port}/` | `http://inventory-dashboard-development.openchoreoapis.localhost:19080/` |

## Data Flow

```
User Browser
    │
    ▼
inventory-dashboard (Next.js)
    │
    │ HTTP GET /inventory
    ▼
inventory-api (Go/Gin)
    │
    │ SQL Query
    ▼
MSSQL Database
    │
    │ Results
    ▼
JSON Response → Rendered Table
```

## Security

- Database credentials managed via External Secrets Operator
- Non-root container users (UID 10014)
- CORS configured for frontend origin only
- No credentials in Git repository

## Extending the Demo

### Adding New Endpoints

1. Add handler in `internal/controllers/inventory_controller.go`
2. Register route in `api/routes/routes.go`
3. Update OpenAPI spec in `docs/openapi.yaml`

### Adding Frontend Features

1. Create component in `src/components/`
2. Import and use in `src/app/page.tsx`
3. Add types in `src/lib/types.ts` if needed

### Changing Database

1. Update repository in `internal/repositories/`
2. Modify connection string in `internal/config/config.go`
3. Update schema in `database/schema.sql`
