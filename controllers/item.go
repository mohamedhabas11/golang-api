package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohamedhabas11/golang-api/database"
	"github.com/mohamedhabas11/golang-api/models"
	"gorm.io/gorm"
)

// CreateItem creates a new item under a given inventory.
func CreateItem(c *fiber.Ctx) error {
	var item models.Item

	// Parse the incoming request body.
	if err := c.BodyParser(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Ensure the associated inventory exists.
	var inventory models.Inventory
	if err := database.DB.First(&inventory, item.InventoryID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Inventory not found")
	}

	// Create the item.
	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(item)
}

// GetItems retrieves all items.
func GetItems(c *fiber.Ctx) error {
	var items []models.Item
	if err := database.DB.Find(&items).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(items)
}

// GetItem retrieves a single item by its ID.
func GetItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.Item

	if err := database.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Item not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(item)
}

// UpdateItem updates an existing item.
func UpdateItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.Item

	// Find the item.
	if err := database.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Item not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Parse update data.
	var updateData models.Item
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Update allowed fields.
	if updateData.Name != "" {
		item.Name = updateData.Name
	}
	if updateData.Quantity != 0 {
		item.Quantity = updateData.Quantity
	}
	// Optionally allow changing the inventory, but check that the new inventory exists.
	if updateData.InventoryID != 0 && updateData.InventoryID != item.InventoryID {
		var newInventory models.Inventory
		if err := database.DB.First(&newInventory, updateData.InventoryID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("New inventory not found")
		}
		item.InventoryID = updateData.InventoryID
	}

	// Save updates.
	if err := database.DB.Save(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(item)
}

// DeleteItem deletes an item by its ID.
func DeleteItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.Item

	// Check if the item exists.
	if err := database.DB.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Item not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Delete the item.
	if err := database.DB.Delete(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).SendString("Item deleted successfully")
}
