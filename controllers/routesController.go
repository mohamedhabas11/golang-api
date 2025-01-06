package controllers

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

// SetupRoutes defines all the routes for the application
func SetupRoutes(app *fiber.App) {
	// Define Default route
	app.Get("/", DefaultRoute)

	// Add a basic health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("OK")
	})

	// Group API routes under `/api` prefix
	api := app.Group("/api")

	// Customer endpoints
	api.Get("/customers", GetCustomers)
	api.Post("/customers", CreateCustomer)
	api.Get("/customers/items", GetCustomersItems)

	// User endpoints
	userGroup := api.Group("/users")      // Group user-related routes
	userGroup.Get("/", GetUsers)          // Fetch all users
	userGroup.Put("/:id", UpdateUser)     // Update a specific user by ID
	userGroup.Get("/:id", GetUser)        // Fetch a specific user by ID
	userGroup.Delete("/:id", DeleteUser)  // Delete a specific user by ID
	userGroup.Post("/signup", CreateUser) // Create a new user
	userGroup.Post("/login", Login)       // Login with existing user, return jwt token

	// Inventory endpoints
	api.Get("/inventories", GetInventories)
	api.Post("/inventories", CreateInventory)

	// Item endpoints
	api.Get("/items", GetItems)
	api.Post("/items", CreateItem)

	// Catch-all route for undefined endpoints
	app.Use(NotFoundRoute)
}
