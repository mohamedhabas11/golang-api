package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohamedhabas11/golang-api/database"
	"github.com/mohamedhabas11/golang-api/models"
)

// CreateCustomer handler for creating a new customer
func CreateCustomer(c *fiber.Ctx) error {
	var customer models.Customer
	// Parse the incoming request body to a customer model
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// check if a customer with the same email already exists
	var existingcustomer models.Customer
	if err := database.DB.Where("email = ?", customer.Email).First(&existingcustomer).Error; err == nil {
		return c.Status(fiber.StatusConflict).SendString("Customer with this email already exists")
	}

	// Insert the customer into the database
	if result := database.DB.Create(&customer); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(customer)
}

// GetCustomers handler to fetch all customers
func GetCustomers(c *fiber.Ctx) error {
	var customers []models.Customer
	// Query the database for all customers
	if result := database.DB.Find(&customers); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	// Return the list of customers
	return c.Status(fiber.StatusOK).JSON(customers)
}

// GetCustomersItems handler to fetch all items
func GetCustomersItems(c *fiber.Ctx) error {
	var customers_items []models.Customer
	// Eager load inventories and items
	if result := database.DB.Preload("Inventories.Items").Find(&customers_items); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.Status(fiber.StatusOK).JSON(customers_items)
}
