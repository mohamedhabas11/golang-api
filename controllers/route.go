package controllers

import (
	"github.com/gofiber/fiber/v2"
)

// DefaultRoute handles the root endpoint.
func DefaultRoute(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Welcome to the Go API!",
	})
}

// NotFoundRoute handles undefined endpoints.
func NotFoundRoute(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error":   "Route not found",
		"message": "The requested resource does not exist.",
	})
}

// SetupRoutes defines all the API routes.
func SetupRoutes(app *fiber.App) {
	// Default routes.
	app.Get("/", DefaultRoute)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// API group.
	api := app.Group("/api")

	// Public routes.
	// Customer registration and login.
	api.Post("/customer/signup", CreateCustomer) // Create a Customer
	api.Post("/customer/login", LoginCustomer)   // Customer login

	// Shop registration and ShopOwner login.
	api.Post("/shop/signup", CreateShop)    // Create a Shop with its ShopOwner
	api.Post("/shop/login", LoginShopOwner) // ShopOwner login

	// (Optional) ShopEmployee signup can be added similarly.
	// api.Post("/employee/signup", CreateEmployee)

	// Protected routes (authentication required).
	protected := api.Group("/")
	// protected.Use(middlewares.RequireAuth) // Uncomment when auth middleware is configured.

	// Customer endpoints.
	protected.Get("/customers", GetCustomers)

	// Shop endpoints.
	protected.Get("/shops", GetShops)
	// (Additional shop update/delete endpoints can be added here)

	// Inventory endpoints.
	protected.Post("/inventories", CreateInventory)
	protected.Get("/inventories", GetInventories)
	protected.Get("/inventories/:id", GetInventory)
	protected.Put("/inventories/:id", UpdateInventory)
	protected.Delete("/inventories/:id", DeleteInventory)

	// Item endpoints.
	protected.Post("/items", CreateItem)
	protected.Get("/items", GetItems)
	protected.Get("/items/:id", GetItem)
	protected.Put("/items/:id", UpdateItem)
	protected.Delete("/items/:id", DeleteItem)

	// ShopEmployee endpoints.
	protected.Post("/employees", CreateEmployee)
	protected.Get("/employees", GetEmployees)
	protected.Get("/employees/:id", GetEmployee)
	protected.Put("/employees/:id", UpdateEmployee)
	protected.Delete("/employees/:id", DeleteEmployee)

	// Catch-all route.
	app.Use(NotFoundRoute)
}
