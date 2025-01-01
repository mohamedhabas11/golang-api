package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohamedhabas11/golang-api/controllers"
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

	// Group API routes under `/api` prefix
	api := app.Group("/api")

	// Customer endpoints
	api.Get("/customers", handlers.GetCustomers)
	api.Post("/customers", handlers.CreateCustomer)
	api.Get("/customers/items", handlers.GetCustomersItems)

	// User endpoints
	userGroup := api.Group("/users")               // Group user-related routes
	userGroup.Get("/", handlers.GetUsers)          // Fetch all users
	userGroup.Put("/:id", handlers.UpdateUser)     // Update a specific user by ID
	userGroup.Get("/:id", handlers.GetUser)        // Fetch a specific user by ID
	userGroup.Delete("/:id", handlers.DeleteUser)  // Delete a specific user by ID
	userGroup.Post("/signup", handlers.CreateUser) // Create a new user
	userGroup.Post("/login", controllers.Login)    // Login with existing user, return jwt token

	// Inventory endpoints
	api.Get("/inventories", handlers.GetInventories)
	api.Post("/inventories", handlers.CreateInventory)

	// Item endpoints
	api.Get("/items", handlers.GetItems)
	api.Post("/items", handlers.CreateItem)

	// Catch-all route for undefined endpoints
	app.Use(handlers.NotFoundRoute)
}
