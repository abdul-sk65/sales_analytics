package api

import (
	"context"
	"log"
	"time"

	"sales_analytics/pkg/repository"

	"github.com/gofiber/fiber/v2"
)

// RefreshData triggers a data refresh from CSV
func (h *Handler) RefreshData(c *fiber.Ctx) error {
	log.Println("Data refresh triggered")

	// Create a background context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create loader and load data
	loader := repository.NewDataLoader(h.repo, h.config.WorkerPoolSize)

	// Run in goroutine for async processing
	go func() {
		if err := loader.LoadCSV(ctx, h.config.CSVFilePath); err != nil {
			log.Printf("Data refresh failed: %v", err)
		}
	}()

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Data refresh initiated",
		"status":  "processing",
	})
}
