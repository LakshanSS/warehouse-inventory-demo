#!/bin/bash

# Database Setup Script for Warehouse Inventory
# This script sets up the MSSQL database with schema and seed data

set -e

# Configuration
DB_SERVER="${DB_SERVER:-localhost}"
DB_PORT="${DB_PORT:-1433}"
DB_USER="${DB_USER:-sa}"
DB_PASSWORD="${DB_PASSWORD:-YourStrong@Passw0rd}"
DB_NAME="${DB_NAME:-InventoryDB}"

echo "=== Warehouse Inventory Database Setup ==="

# Wait for SQL Server to be ready
echo "Waiting for SQL Server to be ready..."
for i in {1..30}; do
    if /opt/mssql-tools18/bin/sqlcmd -S "$DB_SERVER,$DB_PORT" -U "$DB_USER" -P "$DB_PASSWORD" -C -Q "SELECT 1" &>/dev/null; then
        echo "SQL Server is ready!"
        break
    fi
    echo "Waiting... ($i/30)"
    sleep 2
done

# Create database
echo "Creating database $DB_NAME..."
/opt/mssql-tools18/bin/sqlcmd -S "$DB_SERVER,$DB_PORT" -U "$DB_USER" -P "$DB_PASSWORD" -C -Q "
IF NOT EXISTS (SELECT name FROM sys.databases WHERE name = N'$DB_NAME')
BEGIN
    CREATE DATABASE [$DB_NAME];
END
"

# Run schema script
echo "Creating schema..."
/opt/mssql-tools18/bin/sqlcmd -S "$DB_SERVER,$DB_PORT" -U "$DB_USER" -P "$DB_PASSWORD" -C -d "$DB_NAME" -i /scripts/schema.sql

# Run seed script
echo "Seeding data..."
/opt/mssql-tools18/bin/sqlcmd -S "$DB_SERVER,$DB_PORT" -U "$DB_USER" -P "$DB_PASSWORD" -C -d "$DB_NAME" -i /scripts/seed.sql

echo "=== Database setup complete! ==="
echo "Database: $DB_NAME"
echo "Server: $DB_SERVER:$DB_PORT"
