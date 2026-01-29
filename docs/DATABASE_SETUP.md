# Database Setup Guide

This guide covers setting up and connecting to the MSSQL database for the Warehouse Inventory Demo.

## Database Overview

- **Database Type:** Microsoft SQL Server
- **Schema:** Single `Products` table with inventory data
- **Sample Data:** 25 items across 5 categories

## Schema

The `Products` table structure:

| Column | Type | Description |
|--------|------|-------------|
| `ID` | INT | Primary key (auto-increment) |
| `SKU` | NVARCHAR(50) | Unique product identifier |
| `Name` | NVARCHAR(255) | Product name |
| `Quantity` | INT | Stock quantity (>= 0) |
| `Category` | NVARCHAR(100) | Product category |
| `Location` | NVARCHAR(100) | Warehouse location |
| `UpdatedAt` | DATETIME2 | Last update timestamp |

**Indexes:** SKU, Category, Quantity

## Option 1: Local Docker Database

### Start the Container

```bash
cd database
docker-compose up -d
```

### Create the Database

```bash
docker exec -it inventory-mssql /opt/mssql-tools18/bin/sqlcmd \
  -S localhost -U sa -P "YourStrong@Passw0rd" -C \
  -Q "CREATE DATABASE InventoryDB"
```

### Run Schema and Seed

Using DBeaver or Azure Data Studio:

1. Connect to `localhost:1433` with user `sa`
2. Open and run `database/schema.sql`
3. Open and run `database/seed.sql`

### Connection Details

| Setting | Value |
|---------|-------|
| Server | `localhost` |
| Port | `1433` |
| Username | `sa` |
| Password | `YourStrong@Passw0rd` |
| Database | `InventoryDB` |

## Option 2: Azure SQL Database

For production or shared environments, use Azure SQL Database.

### Provisioning

You can provision Azure SQL Database:
- **Manually** via Azure Portal
- **Self-service** using OpenChoreo Infrastructure Workflows (see [infrastructure/README.md](../infrastructure/README.md))

### Running Schema and Seed

1. Connect to your Azure SQL Database using DBeaver or Azure Data Studio
2. Run `database/schema.sql` to create the table
3. Run `database/seed.sql` to insert sample data

### Verify Data

```sql
SELECT * FROM Products;

-- Check stock levels by category
SELECT
    Category,
    COUNT(*) as ItemCount,
    SUM(Quantity) as TotalQuantity
FROM Products
GROUP BY Category;
```

## Sample Data

The seed script inserts 25 items across 5 categories:

| Category | Items | Examples |
|----------|-------|----------|
| Electronics | 5 | Sensors, Power Supplies, Displays |
| Hardware | 5 | Bolts, Brackets, Hinges |
| Safety Equipment | 5 | Goggles, Gloves, Hard Hats |
| Tools | 5 | Drills, Wrenches, Multimeters |
| Packaging | 5 | Boxes, Bubble Wrap, Tape |

Some items have low stock levels to demonstrate the dashboard's alert features.

## Connecting the Backend

Set these environment variables for the backend:

```bash
export DB_SERVER=your-server.database.windows.net
export DB_PORT=1433
export DB_USER=your-username
export DB_PASSWORD=your-password
export DB_NAME=InventoryDB
```

Or in OpenChoreo, use SecretReferences. See [SECRETS_MANAGEMENT.md](SECRETS_MANAGEMENT.md).

## Troubleshooting

### Connection Timeout
- Check firewall rules allow your IP
- For Azure SQL, add your client IP in the Azure Portal

### Login Failed
- Verify username and password
- Ensure the user has access to the database

### Table Not Found
- Run `schema.sql` first
- Verify you're connected to the correct database

### Backend Uses Mock Data
- Backend falls back to mock data when DB connection fails
- Check backend logs for connection errors
- Verify all DB_* environment variables are set
