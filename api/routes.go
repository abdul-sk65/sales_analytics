package api

import (
	"sales_analytics/config"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all API routes
func SetupRoutes(app *fiber.App, cfg *config.Config) {
	handler := NewHandler()

	// Health check
	app.Get("/health", handler.HealthCheck)

}
