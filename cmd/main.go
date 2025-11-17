package main

import (
	"log"

	"sales_analytics/api"
	"sales_analytics/config"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	// Setup routes
	api.SetupRoutes(app, cfg)

	// Start server
	port := cfg.Port
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
