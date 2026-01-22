package repositories

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/microsoft/go-mssqldb"

	"github.com/openchoreo/inventory-api/internal/models"
)

type InventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(connectionString string) (*InventoryRepository, error) {
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return nil, err
	}

	// Configure connection pool for Azure SQL
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute) // Close idle connections before Azure drops them

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return &InventoryRepository{db: db}, nil
}

func (r *InventoryRepository) Close() error {
	return r.db.Close()
}

func (r *InventoryRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *InventoryRepository) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	query := `
		SELECT ID, SKU, Name, Quantity, ISNULL(Category, ''), ISNULL(Location, ''), UpdatedAt
		FROM Products
		ORDER BY Name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Quantity, &p.Category, &p.Location, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *InventoryRepository) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	query := `
		SELECT ID, SKU, Name, Quantity, ISNULL(Category, ''), ISNULL(Location, ''), UpdatedAt
		FROM Products
		WHERE ID = @p1
	`

	var p models.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.SKU, &p.Name, &p.Quantity, &p.Category, &p.Location, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *InventoryRepository) CreateProduct(ctx context.Context, req models.CreateProductRequest) (*models.Product, error) {
	query := `
		INSERT INTO Products (SKU, Name, Quantity, Category, Location, UpdatedAt)
		OUTPUT INSERTED.ID, INSERTED.SKU, INSERTED.Name, INSERTED.Quantity,
		       ISNULL(INSERTED.Category, ''), ISNULL(INSERTED.Location, ''), INSERTED.UpdatedAt
		VALUES (@p1, @p2, @p3, @p4, @p5, GETUTCDATE())
	`

	var p models.Product
	err := r.db.QueryRowContext(ctx, query, req.SKU, req.Name, req.Quantity, req.Category, req.Location).Scan(
		&p.ID, &p.SKU, &p.Name, &p.Quantity, &p.Category, &p.Location, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *InventoryRepository) UpdateProductQuantity(ctx context.Context, id int, quantity int) (*models.Product, error) {
	query := `
		UPDATE Products
		SET Quantity = @p1, UpdatedAt = GETUTCDATE()
		OUTPUT INSERTED.ID, INSERTED.SKU, INSERTED.Name, INSERTED.Quantity,
		       ISNULL(INSERTED.Category, ''), ISNULL(INSERTED.Location, ''), INSERTED.UpdatedAt
		WHERE ID = @p2
	`

	var p models.Product
	err := r.db.QueryRowContext(ctx, query, quantity, id).Scan(
		&p.ID, &p.SKU, &p.Name, &p.Quantity, &p.Category, &p.Location, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *InventoryRepository) DeleteProduct(ctx context.Context, id int) error {
	query := `DELETE FROM Products WHERE ID = @p1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
