package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests
type Handler struct{}

// NewHandler creates a new handler
func NewHandler() *Handler {
	return &Handler{}
}

// HealthCheck returns the health status of the API
func (h *Handler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "healthy",
		"time":   time.Now(),
	})
}
