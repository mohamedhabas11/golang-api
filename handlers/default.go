package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// DefaultRoute handler for the root endpoint
func DefaultRoute(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Welcome to the Go API!",
	})
}

// NotFoundRoute handler for undefined routes
func NotFoundRoute(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error":   "Route not found",
		"message": "The requested resource does not exist.",
	})
}
