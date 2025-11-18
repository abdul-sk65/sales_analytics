package api

import (
	"time"

	"sales_analytics/config"
	"sales_analytics/pkg/repository"
	"sales_analytics/pkg/scheduler"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests
type Handler struct {
	repo      *repository.MongoRepository
	config    *config.Config
	scheduler *scheduler.Scheduler
}

// NewHandler creates a new handler
func NewHandler(repo *repository.MongoRepository, cfg *config.Config, sched *scheduler.Scheduler) *Handler {
	return &Handler{
		repo:      repo,
		config:    cfg,
		scheduler: sched,
	}
}

// HealthCheck returns the health status of the API
func (h *Handler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "healthy",
		"time":   time.Now(),
	})
}
