package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohamedhabas11/golang-api/handlers"
)

// SetupRoutes defines all the routes for the application
func SetupRoutes(app *fiber.App) {
	// Define Default route
	app.Get("/", handlers.DefaultRoute)

	// Define routes for customers
	app.Get("/customers", handlers.GetCustomers)
	app.Post("/customers", handlers.CreateCustomer)

	// Define routes for inventories
	app.Get("/inventories", handlers.GetInventories)
	app.Post("/inventories", handlers.CreateInventory)

	// Define routes for items
	app.Get("/items", handlers.GetItems)
	app.Post("/items", handlers.CreateItem)
}
