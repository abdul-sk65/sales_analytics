package api

import (
	"github.com/gofiber/fiber/v2"
)

// CreateCronJobRequest represents the request body for creating a cron job
type CreateCronJobRequest struct {
	Interval string `json:"interval"`
}

// CreateCronJob creates or replaces a cron job for automated data refresh
func (h *Handler) CreateCronJob(c *fiber.Ctx) error {
	req := new(CreateCronJobRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if req.Interval == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "interval is required (e.g., '1h', '30m', '24h')",
		})
	}

	if err := h.scheduler.CreateJob(req.Interval); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "Cron job created successfully (replaced existing if any)",
		"interval": req.Interval,
		"status":   h.scheduler.GetJobStatus(),
	})
}

// DeleteCronJob deletes the active cron job
func (h *Handler) DeleteCronJob(c *fiber.Ctx) error {
	if err := h.scheduler.DeleteJob(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Cron job deleted successfully",
	})
}

// GetCronStatus returns the status of the cron job
func (h *Handler) GetCronStatus(c *fiber.Ctx) error {
	status := h.scheduler.GetJobStatus()

	return c.JSON(fiber.Map{
		"status": status,
	})
}
