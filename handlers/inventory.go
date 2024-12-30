package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohamedhabas11/golang-api/database"
	"github.com/mohamedhabas11/golang-api/models"
)

// CreateInventory handler for creating a new inventory
func CreateInventory(c *fiber.Ctx) error {
	var inventory models.Inventory
	// Parse the incoming request body to an inventory model
	if err := c.BodyParser(&inventory); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Check if the associated customer exists
	var customer models.Customer
	if err := database.DB.First(&customer, inventory.CustomerID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Customer not found")
	}

	if result := database.DB.Create(&inventory); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(inventory)
}

// GetInventories handler to fetch all inventories
func GetInventories(c *fiber.Ctx) error {
	var inventories []models.Inventory
	if result := database.DB.Find(&inventories); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.Status(fiber.StatusOK).JSON(inventories)
}
