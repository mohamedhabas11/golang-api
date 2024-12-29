package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohamedhabas11/golang-api/database"
	"github.com/mohamedhabas11/golang-api/models"
)

// CreateInventory handler for creating a new inventory
func CreateInventory(c *fiber.Ctx) error {
	var inventory models.Inventory
	if err := c.BodyParser(&inventory); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if result := database.DB.Db.Create(&inventory); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(inventory)
}

// GetInventories handler to fetch all inventories
func GetInventories(c *fiber.Ctx) error {
	var inventories []models.Inventory
	if result := database.DB.Db.Find(&inventories); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.Status(fiber.StatusOK).JSON(inventories)
}
