-- Warehouse Inventory Database Schema
-- MSSQL Server

-- Create database (run this separately if needed)
-- CREATE DATABASE InventoryDB;
-- GO

-- Use the database
-- USE InventoryDB;
-- GO

-- Create Products table
IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'Products')
BEGIN
    CREATE TABLE Products (
        ID INT IDENTITY(1,1) PRIMARY KEY,
        SKU NVARCHAR(50) NOT NULL UNIQUE,
        Name NVARCHAR(255) NOT NULL,
        Quantity INT NOT NULL DEFAULT 0,
        Category NVARCHAR(100) NULL,
        Location NVARCHAR(100) NULL,
        UpdatedAt DATETIME2 NOT NULL DEFAULT GETUTCDATE(),

        CONSTRAINT CK_Products_Quantity CHECK (Quantity >= 0)
    );

    -- Create indexes for common queries
    CREATE INDEX IX_Products_SKU ON Products(SKU);
    CREATE INDEX IX_Products_Category ON Products(Category);
    CREATE INDEX IX_Products_Quantity ON Products(Quantity);
END
GO

-- Create a trigger to update the UpdatedAt timestamp
IF EXISTS (SELECT * FROM sys.triggers WHERE name = 'TR_Products_UpdateTimestamp')
    DROP TRIGGER TR_Products_UpdateTimestamp;
GO

CREATE TRIGGER TR_Products_UpdateTimestamp
ON Products
AFTER UPDATE
AS
BEGIN
    SET NOCOUNT ON;
    UPDATE Products
    SET UpdatedAt = GETUTCDATE()
    FROM Products p
    INNER JOIN inserted i ON p.ID = i.ID;
END
GO
