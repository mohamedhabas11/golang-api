package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohamedhabas11/golang-api/handlers"
)

// SetupRoutes defines all the routes for the application
func SetupRoutes(app *fiber.App) {
	// Define Default route
	app.Get("/", handlers.DefaultRoute)

	// Add a basic health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("OK")
	})

	// Define routes for customers
	api := app.Group("/api") // Grouping all routes under /api prefix
	api.Get("/customers", handlers.GetCustomers)
	api.Post("/customers", handlers.CreateCustomer)
	api.Get("/customers/items", handlers.GetCustomersItems)

	// Define routes for inventories
	api.Get("/inventories", handlers.GetInventories)
	api.Post("/inventories", handlers.CreateInventory)

	// Define routes for items
	api.Get("/items", handlers.GetItems)
	api.Post("/items", handlers.CreateItem)

	// Catch-all route for undefined endpoints
	app.Use(handlers.NotFoundRoute)
}
