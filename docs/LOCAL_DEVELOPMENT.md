# Local Development Guide

This guide covers running the Warehouse Inventory Demo locally for development and testing.

## Prerequisites

- **Go 1.22+** - For the backend API
- **Node.js 18+** - For the frontend
- **Docker** - For local MSSQL database (optional)

## Quick Start

### 1. Start the Database (Optional)

The backend returns mock data when no database is connected, so this step is optional.

```bash
cd database
docker-compose up -d
```

Wait for MSSQL to be ready (~30 seconds), then create the database:

```bash
docker exec -it inventory-mssql /opt/mssql-tools18/bin/sqlcmd \
  -S localhost -U sa -P "YourStrong@Passw0rd" -C \
  -Q "CREATE DATABASE InventoryDB"
```

Run the schema and seed scripts using DBeaver, Azure Data Studio, or any SQL client. See [DATABASE_SETUP.md](DATABASE_SETUP.md) for details.

### 2. Start the Backend

```bash
cd inventory-api
go mod download
go run main.go
```

The API will be available at `http://localhost:9090`.

**Test it:**
```bash
curl http://localhost:9090/health
curl http://localhost:9090/inventory
```

### 3. Start the Frontend

```bash
cd inventory-dashboard
npm install
npm run dev
```

The dashboard will be available at `http://localhost:3000`.

## Environment Variables

### Backend (inventory-api)

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | API server port | `9090` |
| `ENVIRONMENT` | Runtime environment | `development` |
| `DB_SERVER` | MSSQL server address | `localhost` |
| `DB_PORT` | MSSQL server port | `1433` |
| `DB_USER` | Database username | `sa` |
| `DB_PASSWORD` | Database password | - |
| `DB_NAME` | Database name | `InventoryDB` |

**Note:** When `DB_SERVER` is not set or connection fails, the backend returns mock data (15 sample items).

### Frontend (inventory-dashboard)

The frontend automatically detects the API URL based on the environment:
- **Local development:** Uses `http://localhost:9090`
- **OpenChoreo:** Detects from hostname pattern

No environment variables needed for basic local development.

## Project Structure

```
inventory-api/
├── main.go              # Entry point
├── api/routes/          # Route registration
└── internal/
    ├── config/          # Configuration loader
    ├── controllers/     # HTTP handlers (includes mock data)
    ├── models/          # Data models
    └── repositories/    # Database layer

inventory-dashboard/
├── src/
│   ├── app/             # Next.js App Router pages
│   ├── components/      # React components
│   └── lib/             # Utilities and types
└── public/              # Static assets
```

## Common Tasks

### Running Tests

```bash
# Backend (if tests exist)
cd inventory-api
go test ./...

# Frontend
cd inventory-dashboard
npm test
```

### Building Docker Images

```bash
# Backend
cd inventory-api
docker build -t inventory-api:local .

# Frontend
cd inventory-dashboard
docker build -t inventory-dashboard:local .
```

### Connecting to Local Database

Use DBeaver or Azure Data Studio:
- **Server:** `localhost`
- **Port:** `1433`
- **Username:** `sa`
- **Password:** `YourStrong@Passw0rd`
- **Database:** `InventoryDB`

## Troubleshooting

### Backend won't start
- Check if port 9090 is already in use
- Verify Go version: `go version`

### Frontend won't start
- Check if port 3000 is already in use
- Clear node_modules: `rm -rf node_modules && npm install`

### Database connection fails
- Verify Docker container is running: `docker ps`
- Check MSSQL logs: `docker logs inventory-mssql`
- The backend will fall back to mock data automatically

### API returns 404
- Ensure you're using the correct path: `/inventory` not `/api/inventory`
- Check the backend is running: `curl http://localhost:9090/health`
