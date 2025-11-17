package api

import (
	"time"

	"sales_analytics/config"
	"sales_analytics/pkg/repository"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests
type Handler struct {
	repo   *repository.MongoRepository
	config *config.Config
}

// NewHandler creates a new handler
func NewHandler(repo *repository.MongoRepository, cfg *config.Config) *Handler {
	return &Handler{
		repo:   repo,
		config: cfg,
	}
}

// HealthCheck returns the health status of the API
func (h *Handler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "healthy",
		"time":   time.Now(),
	})
}
