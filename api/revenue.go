package api

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetTotalRevenue calculates total revenue for a date range
func (h *Handler) GetTotalRevenue(c *fiber.Ctx) error {
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	revenue, err := h.repo.CalculateTotalRevenue(ctx, startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate total revenue",
		})
	}

	return c.JSON(fiber.Map{
		"start_date":    startDate.Format("2006-01-02"),
		"end_date":      endDate.Format("2006-01-02"),
		"total_revenue": revenue,
	})
}

// GetRevenueByProduct calculates revenue grouped by product
func (h *Handler) GetRevenueByProduct(c *fiber.Ctx) error {
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	results, err := h.repo.CalculateRevenueByProduct(ctx, startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate revenue by product",
		})
	}

	return c.JSON(fiber.Map{
		"start_date":       startDate.Format("2006-01-02"),
		"end_date":         endDate.Format("2006-01-02"),
		"products_revenue": results,
	})
}

// GetRevenueByCategory calculates revenue grouped by category
func (h *Handler) GetRevenueByCategory(c *fiber.Ctx) error {
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	results, err := h.repo.CalculateRevenueByCategory(ctx, startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate revenue by category",
		})
	}

	return c.JSON(fiber.Map{
		"start_date":         startDate.Format("2006-01-02"),
		"end_date":           endDate.Format("2006-01-02"),
		"categories_revenue": results,
	})
}

// GetRevenueByRegion calculates revenue grouped by region
func (h *Handler) GetRevenueByRegion(c *fiber.Ctx) error {
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	results, err := h.repo.CalculateRevenueByRegion(ctx, startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate revenue by region",
		})
	}

	return c.JSON(fiber.Map{
		"start_date":      startDate.Format("2006-01-02"),
		"end_date":        endDate.Format("2006-01-02"),
		"regions_revenue": results,
	})
}

// parseDateRange extracts and validates start_date and end_date from query params
func (h *Handler) parseDateRange(c *fiber.Ctx) (time.Time, time.Time, error) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		return time.Time{}, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "start_date and end_date are required")
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "invalid start_date format, use YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "invalid end_date format, use YYYY-MM-DD")
	}

	if endDate.Before(startDate) {
		return time.Time{}, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "end_date must be after start_date")
	}

	// Add 23:59:59 to end date to include the entire day
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	return startDate, endDate, nil
}
