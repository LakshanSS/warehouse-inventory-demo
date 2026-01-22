// Warehouse Inventory API
// A RESTful service for managing warehouse inventory with MSSQL backend

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/openchoreo/inventory-api/api/routes"
	"github.com/openchoreo/inventory-api/internal/config"
	"github.com/openchoreo/inventory-api/internal/repositories"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database connection
	repo, err := repositories.NewInventoryRepository(cfg.DatabaseURL)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}
	defer repo.Close()

	// Create Gin router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Initialize routes
	routes.Initialize(router, repo)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logrus.WithFields(logrus.Fields{"port": cfg.Port}).Info("Inventory API service is starting...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logrus.Info("SIGTERM received: gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("Server shutdown error: %v", err)
	}

	logrus.Info("Server exited")
}
