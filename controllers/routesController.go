package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohamedhabas11/golang-api/middlewares"
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

// SetupRoutes defines all the routes for the application
func SetupRoutes(app *fiber.App) {
	// Default route
	app.Get("/", DefaultRoute)
	app.Get("/health", func(c *fiber.Ctx) error { return c.Status(fiber.StatusOK).SendString("OK") })

	// Group API routes
	api := app.Group("/api")

	// Public routes (No authentication required)
	api.Post("/user/signup", CreateUser)
	api.Post("/users/login", Login)

	// Protected routes (Login required)
	protected := api.Group("/")
	protected.Use(middlewares.RequireAuth)

	// Customer endpoints
	protected.Post("/customers", CreateCustomer)
	protected.Get("/customers", GetCustomers)
	protected.Get("/customers/items", GetCustomersItems)

	// Inventory endpoints
	protected.Post("/inventories", CreateInventory)
	protected.Get("/inventories", GetInventories)

	// Item endpoint
	protected.Post("/items", CreateItem)
	protected.Get("/items", GetItems)

	// User endpoints
	userGroup := protected.Group("/users") // Group user-related routes
	userGroup.Get("/", GetUsers)           // Fetch all users
	userGroup.Put("/:id", UpdateUser)      // Update a specific user by ID
	userGroup.Get("/:id", GetUser)         // Fetch a specific user by ID
	userGroup.Delete("/:id", DeleteUser)   // Delete a specific user by ID

	// Catch-all route for undefined endpoints
	app.Use(NotFoundRoute)
}
