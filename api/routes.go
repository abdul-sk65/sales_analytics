package api

import (
	"sales_analytics/config"
	"sales_analytics/pkg/repository"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all API routes
func SetupRoutes(app *fiber.App, repo *repository.MongoRepository, cfg *config.Config) {
	handler := NewHandler(repo, cfg)

	// Health check
	app.Get("/health", handler.HealthCheck)

	// Data refresh endpoints
	api := app.Group("/api/v1")
	dataRefresh := api.Group("/data")
	dataRefresh.Post("/refresh", handler.RefreshData)
	dataRefresh.Get("/logs", handler.GetRefreshLogs)

	// Revenue analytics endpoints
	revenue := api.Group("/revenue")
	revenue.Get("/total", handler.GetTotalRevenue)
	revenue.Get("/by-product", handler.GetRevenueByProduct)
	revenue.Get("/by-category", handler.GetRevenueByCategory)
	revenue.Get("/by-region", handler.GetRevenueByRegion)
}
