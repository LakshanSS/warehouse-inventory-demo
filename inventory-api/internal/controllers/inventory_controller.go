package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/openchoreo/inventory-api/internal/models"
	"github.com/openchoreo/inventory-api/internal/repositories"
)

type InventoryController struct {
	repo     *repositories.InventoryRepository
	mockData []models.Product
}

func NewInventoryController(repo *repositories.InventoryRepository) *InventoryController {
	return &InventoryController{
		repo:     repo,
		mockData: getMockData(),
	}
}

func getMockData() []models.Product {
	now := time.Now().UTC()
	return []models.Product{
		{ID: 1, SKU: "ELEC-001", Name: "Industrial Sensor Module", Quantity: 150, Category: "Electronics", Location: "Aisle A, Shelf 1", UpdatedAt: now},
		{ID: 2, SKU: "ELEC-002", Name: "Power Supply Unit 500W", Quantity: 75, Category: "Electronics", Location: "Aisle A, Shelf 2", UpdatedAt: now},
		{ID: 3, SKU: "ELEC-003", Name: "Control Panel Display", Quantity: 8, Category: "Electronics", Location: "Aisle A, Shelf 3", UpdatedAt: now},
		{ID: 4, SKU: "HW-001", Name: "Steel Bolts M10 (Box)", Quantity: 500, Category: "Hardware", Location: "Aisle B, Shelf 1", UpdatedAt: now},
		{ID: 5, SKU: "HW-002", Name: "Aluminum Brackets Set", Quantity: 120, Category: "Hardware", Location: "Aisle B, Shelf 2", UpdatedAt: now},
		{ID: 6, SKU: "HW-003", Name: "Rubber Gaskets Pack", Quantity: 3, Category: "Hardware", Location: "Aisle B, Shelf 4", UpdatedAt: now},
		{ID: 7, SKU: "SAFE-001", Name: "Safety Goggles", Quantity: 50, Category: "Safety Equipment", Location: "Aisle C, Shelf 1", UpdatedAt: now},
		{ID: 8, SKU: "SAFE-002", Name: "Hard Hat", Quantity: 25, Category: "Safety Equipment", Location: "Aisle C, Shelf 3", UpdatedAt: now},
		{ID: 9, SKU: "SAFE-003", Name: "First Aid Kit", Quantity: 5, Category: "Safety Equipment", Location: "Aisle C, Shelf 4", UpdatedAt: now},
		{ID: 10, SKU: "TOOL-001", Name: "Power Drill Set", Quantity: 20, Category: "Tools", Location: "Aisle D, Shelf 1", UpdatedAt: now},
		{ID: 11, SKU: "TOOL-002", Name: "Digital Multimeter", Quantity: 15, Category: "Tools", Location: "Aisle D, Shelf 3", UpdatedAt: now},
		{ID: 12, SKU: "TOOL-003", Name: "Precision Screwdriver Set", Quantity: 2, Category: "Tools", Location: "Aisle D, Shelf 5", UpdatedAt: now},
		{ID: 13, SKU: "PKG-001", Name: "Cardboard Boxes Large", Quantity: 250, Category: "Packaging", Location: "Aisle E, Shelf 1", UpdatedAt: now},
		{ID: 14, SKU: "PKG-002", Name: "Bubble Wrap Roll", Quantity: 80, Category: "Packaging", Location: "Aisle E, Shelf 2", UpdatedAt: now},
		{ID: 15, SKU: "PKG-003", Name: "Foam Padding Sheets", Quantity: 4, Category: "Packaging", Location: "Aisle E, Shelf 4", UpdatedAt: now},
	}
}

// Health check endpoint
func (c *InventoryController) Health(ctx *gin.Context) {
	dbStatus := "not_configured"

	if c.repo != nil {
		reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Second)
		defer cancel()

		dbStatus = "healthy"
		if err := c.repo.Ping(reqCtx); err != nil {
			dbStatus = "unhealthy"
		}
	}

	ctx.JSON(http.StatusOK, models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Database:  dbStatus,
	})
}

// Get all inventory items
func (c *InventoryController) GetInventory(ctx *gin.Context) {
	if c.repo == nil {
		// Return mock data when database is not available
		ctx.JSON(http.StatusOK, c.mockData)
		return
	}

	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
	defer cancel()

	products, err := c.repo.GetAllProducts(reqCtx)
	if err != nil {
		// Fallback to mock data on database error
		ctx.JSON(http.StatusOK, c.mockData)
		return
	}

	if products == nil {
		products = []models.Product{}
	}

	ctx.JSON(http.StatusOK, products)
}

// Get single inventory item by ID
func (c *InventoryController) GetInventoryItem(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_id",
			Message: "ID must be a valid integer",
		})
		return
	}

	if c.repo == nil {
		// Return mock data when database is not available
		for _, p := range c.mockData {
			if p.ID == id {
				ctx.JSON(http.StatusOK, p)
				return
			}
		}
		ctx.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Product not found",
		})
		return
	}

	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	product, err := c.repo.GetProductByID(reqCtx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve product",
		})
		return
	}

	if product == nil {
		ctx.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Product not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// Create new inventory item
func (c *InventoryController) CreateInventoryItem(ctx *gin.Context) {
	var req models.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	if c.repo == nil {
		// Return mock created product in demo mode
		product := models.Product{
			ID:        100,
			SKU:       req.SKU,
			Name:      req.Name,
			Quantity:  req.Quantity,
			Category:  req.Category,
			Location:  req.Location,
			UpdatedAt: time.Now().UTC(),
		}
		ctx.JSON(http.StatusCreated, product)
		return
	}

	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	product, err := c.repo.CreateProduct(reqCtx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to create product",
		})
		return
	}

	ctx.JSON(http.StatusCreated, product)
}

// Update inventory quantity
func (c *InventoryController) UpdateInventoryQuantity(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_id",
			Message: "ID must be a valid integer",
		})
		return
	}

	var req models.UpdateQuantityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	if c.repo == nil {
		// Return mock updated product in demo mode
		for _, p := range c.mockData {
			if p.ID == id {
				p.Quantity = req.Quantity
				p.UpdatedAt = time.Now().UTC()
				ctx.JSON(http.StatusOK, p)
				return
			}
		}
		ctx.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Product not found",
		})
		return
	}

	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	product, err := c.repo.UpdateProductQuantity(reqCtx, id, req.Quantity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to update product",
		})
		return
	}

	if product == nil {
		ctx.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Product not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// Delete inventory item
func (c *InventoryController) DeleteInventoryItem(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_id",
			Message: "ID must be a valid integer",
		})
		return
	}

	if c.repo == nil {
		// Return success in demo mode
		for _, p := range c.mockData {
			if p.ID == id {
				ctx.Status(http.StatusNoContent)
				return
			}
		}
		ctx.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Product not found",
		})
		return
	}

	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	if err := c.repo.DeleteProduct(reqCtx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to delete product",
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
