package models

import "time"

type Product struct {
	ID        int       `json:"id"`
	SKU       string    `json:"sku"`
	Name      string    `json:"name"`
	Quantity  int       `json:"quantity"`
	Category  string    `json:"category,omitempty"`
	Location  string    `json:"location,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateQuantityRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}

type CreateProductRequest struct {
	SKU      string `json:"sku" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Quantity int    `json:"quantity" binding:"min=0"`
	Category string `json:"category"`
	Location string `json:"location"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Database  string `json:"database"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
