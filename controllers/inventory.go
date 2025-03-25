package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohamedhabas11/golang-api/database"
	"github.com/mohamedhabas11/golang-api/models"
	"gorm.io/gorm"
)

// CreateInventory creates a new inventory for a shop.
func CreateInventory(c *fiber.Ctx) error {
	var inventory models.Inventory

	// Parse the incoming request body.
	if err := c.BodyParser(&inventory); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Ensure the associated shop exists.
	var shop models.Shop
	if err := database.DB.First(&shop, inventory.ShopID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Shop not found")
	}

	// Create the inventory.
	if err := database.DB.Create(&inventory).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(inventory)
}

// GetInventories retrieves all inventories with their items.
func GetInventories(c *fiber.Ctx) error {
	var inventories []models.Inventory

	// Preload Items so that each inventory returns its items.
	if err := database.DB.Preload("Items").Find(&inventories).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(inventories)
}

// GetInventory retrieves a single inventory by ID, including its items.
func GetInventory(c *fiber.Ctx) error {
	id := c.Params("id")
	var inventory models.Inventory

	if err := database.DB.Preload("Items").First(&inventory, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Inventory not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(inventory)
}

// UpdateInventory updates an existing inventory record.
func UpdateInventory(c *fiber.Ctx) error {
	id := c.Params("id")
	var inventory models.Inventory

	// Find the inventory by ID.
	if err := database.DB.First(&inventory, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Inventory not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Parse update data.
	var updateData models.Inventory
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Update only allowed fields (for example, InventoryName).
	if updateData.InventoryName != "" {
		inventory.InventoryName = updateData.InventoryName
	}

	// Save updates.
	if err := database.DB.Save(&inventory).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(inventory)
}

// DeleteInventory deletes an inventory by its ID.
func DeleteInventory(c *fiber.Ctx) error {
	id := c.Params("id")
	var inventory models.Inventory

	// Ensure the inventory exists.
	if err := database.DB.First(&inventory, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Inventory not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Delete the inventory.
	if err := database.DB.Delete(&inventory).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).SendString("Inventory deleted successfully")
}
