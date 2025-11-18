package api

import (
	"sales_analytics/config"
	"sales_analytics/pkg/repository"
	"sales_analytics/pkg/scheduler"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all API routes
func SetupRoutes(app *fiber.App, repo *repository.MongoRepository, cfg *config.Config, sched *scheduler.Scheduler) {
	handler := NewHandler(repo, cfg, sched)

	// Health check
	app.Get("/health", handler.HealthCheck)

	api := app.Group("/api/v1")

	// Data refresh endpoints
	dataRefresh := api.Group("/data")
	dataRefresh.Post("/refresh", handler.RefreshData)
	dataRefresh.Get("/logs", handler.GetRefreshLogs)

	// Cron job management endpoints
	cron := api.Group("/cron")
	cron.Post("/create", handler.CreateCronJob)
	cron.Delete("/delete", handler.DeleteCronJob)
	cron.Get("/status", handler.GetCronStatus)

	// Revenue analytics endpoints
	revenue := api.Group("/revenue")
	revenue.Get("/total", handler.GetTotalRevenue)
	revenue.Get("/product", handler.GetRevenueByProduct)
	revenue.Get("/category", handler.GetRevenueByCategory)
	revenue.Get("/region", handler.GetRevenueByRegion)
}
