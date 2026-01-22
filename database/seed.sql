-- Seed data for Warehouse Inventory
-- USE InventoryDB;
-- GO

-- Clear existing data (optional - remove in production)
-- DELETE FROM Products;

-- Insert sample inventory items
INSERT INTO Products (SKU, Name, Quantity, Category, Location)
VALUES
    -- Electronics
    ('ELEC-001', 'Industrial Sensor Module', 150, 'Electronics', 'Aisle A, Shelf 1'),
    ('ELEC-002', 'Power Supply Unit 500W', 75, 'Electronics', 'Aisle A, Shelf 2'),
    ('ELEC-003', 'Control Panel Display', 30, 'Electronics', 'Aisle A, Shelf 3'),
    ('ELEC-004', 'Circuit Breaker 30A', 200, 'Electronics', 'Aisle A, Shelf 4'),
    ('ELEC-005', 'LED Light Panel', 8, 'Electronics', 'Aisle A, Shelf 5'),

    -- Hardware
    ('HW-001', 'Steel Bolts M10 (Box)', 500, 'Hardware', 'Aisle B, Shelf 1'),
    ('HW-002', 'Aluminum Brackets Set', 120, 'Hardware', 'Aisle B, Shelf 2'),
    ('HW-003', 'Industrial Hinges (Pair)', 45, 'Hardware', 'Aisle B, Shelf 3'),
    ('HW-004', 'Rubber Gaskets Pack', 3, 'Hardware', 'Aisle B, Shelf 4'),
    ('HW-005', 'Stainless Steel Screws', 1000, 'Hardware', 'Aisle B, Shelf 5'),

    -- Safety Equipment
    ('SAFE-001', 'Safety Goggles', 50, 'Safety Equipment', 'Aisle C, Shelf 1'),
    ('SAFE-002', 'Work Gloves (Pair)', 100, 'Safety Equipment', 'Aisle C, Shelf 2'),
    ('SAFE-003', 'Hard Hat', 25, 'Safety Equipment', 'Aisle C, Shelf 3'),
    ('SAFE-004', 'First Aid Kit', 5, 'Safety Equipment', 'Aisle C, Shelf 4'),
    ('SAFE-005', 'Fire Extinguisher', 12, 'Safety Equipment', 'Aisle C, Shelf 5'),

    -- Tools
    ('TOOL-001', 'Power Drill Set', 20, 'Tools', 'Aisle D, Shelf 1'),
    ('TOOL-002', 'Wrench Set (Metric)', 35, 'Tools', 'Aisle D, Shelf 2'),
    ('TOOL-003', 'Digital Multimeter', 15, 'Tools', 'Aisle D, Shelf 3'),
    ('TOOL-004', 'Soldering Station', 10, 'Tools', 'Aisle D, Shelf 4'),
    ('TOOL-005', 'Precision Screwdriver Set', 2, 'Tools', 'Aisle D, Shelf 5'),

    -- Packaging
    ('PKG-001', 'Cardboard Boxes Large', 250, 'Packaging', 'Aisle E, Shelf 1'),
    ('PKG-002', 'Bubble Wrap Roll', 80, 'Packaging', 'Aisle E, Shelf 2'),
    ('PKG-003', 'Packing Tape', 150, 'Packaging', 'Aisle E, Shelf 3'),
    ('PKG-004', 'Foam Padding Sheets', 4, 'Packaging', 'Aisle E, Shelf 4'),
    ('PKG-005', 'Shipping Labels (Roll)', 500, 'Packaging', 'Aisle E, Shelf 5');

GO

-- Verify the data
SELECT
    Category,
    COUNT(*) as ItemCount,
    SUM(Quantity) as TotalQuantity,
    SUM(CASE WHEN Quantity <= 5 THEN 1 ELSE 0 END) as CriticalStock,
    SUM(CASE WHEN Quantity > 5 AND Quantity <= 10 THEN 1 ELSE 0 END) as LowStock
FROM Products
GROUP BY Category
ORDER BY Category;
GO
