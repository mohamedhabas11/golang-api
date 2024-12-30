package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohamedhabas11/golang-api/database"
	"github.com/mohamedhabas11/golang-api/models"
)

// CreateItem handler for creating a new item
func CreateItem(c *fiber.Ctx) error {
	var item models.Item
	if err := c.BodyParser(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Check if the associated inventory exists
	var inventory models.Inventory
	if err := database.DB.First(&inventory, item.InventoryID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Inventory not found")
	}

	if result := database.DB.Create(&item); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(item)
}

// GetItems handler to fetch all items
func GetItems(c *fiber.Ctx) error {
	var items []models.Item
	if result := database.DB.Find(&items); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.Status(fiber.StatusOK).JSON(items)
}
