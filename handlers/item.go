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

	if result := database.DB.Db.Create(&item); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(item)
}

// GetItems handler to fetch all items
func GetItems(c *fiber.Ctx) error {
	var items []models.Item
	if result := database.DB.Db.Find(&items); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.Status(fiber.StatusOK).JSON(items)
}
