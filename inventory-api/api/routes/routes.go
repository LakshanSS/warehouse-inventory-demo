package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/openchoreo/inventory-api/internal/controllers"
	"github.com/openchoreo/inventory-api/internal/repositories"
)

func Initialize(router *gin.Engine, repo *repositories.InventoryRepository) {
	controller := controllers.NewInventoryController(repo)

	// Health check endpoint
	router.GET("/health", controller.Health)

	// Inventory API endpoints
	api := router.Group("/inventory")
	{
		api.GET("", controller.GetInventory)
		api.GET("/:id", controller.GetInventoryItem)
		api.POST("", controller.CreateInventoryItem)
		api.PUT("/:id", controller.UpdateInventoryQuantity)
		api.DELETE("/:id", controller.DeleteInventoryItem)
	}
}
