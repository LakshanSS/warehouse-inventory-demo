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
	repo *repositories.InventoryRepository
}

func NewInventoryController(repo *repositories.InventoryRepository) *InventoryController {
	return &InventoryController{repo: repo}
}

// Health check endpoint
func (c *InventoryController) Health(ctx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Second)
	defer cancel()

	dbStatus := "healthy"
	if err := c.repo.Ping(reqCtx); err != nil {
		dbStatus = "unhealthy"
	}

	ctx.JSON(http.StatusOK, models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Database:  dbStatus,
	})
}

// Get all inventory items
func (c *InventoryController) GetInventory(ctx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
	defer cancel()

	products, err := c.repo.GetAllProducts(reqCtx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve inventory",
		})
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
