package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohamedhabas11/golang-api/database"
	"github.com/mohamedhabas11/golang-api/models"
	"github.com/mohamedhabas11/golang-api/utils"
)

// CreateCustomer handler for creating a new customer
func CreateCustomer(c *fiber.Ctx) error {
	var customer models.Customer
	// Parse the incoming request body to a customer model
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// validate email format
	if !utils.ValidateEmail(customer.Email) {
		return c.Status(fiber.StatusBadRequest).SendString("invalid email format")
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

func GetCustomers(c *fiber.Ctx) error {
	var customers []models.Customer

	// Exclude relationships to avoid loading unnecessary data
	if result := database.DB.Model(&models.Customer{}).Select("id", "created_at", "updated_at", "name", "email").Find(&customers); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.Status(fiber.StatusOK).JSON(customers)
}

// GetCustomersItems handler to fetch all inventories and their items for customers
func GetCustomersItems(c *fiber.Ctx) error {
	var customers []models.Customer

	// Eager load inventories and items while excluding unrelated data
	if err := database.DB.
		Preload("Inventories.Items").
		Find(&customers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Transform the data to return only inventories and items
	var response []struct {
		CustomerName string `json:"customer_name"`
		Inventories  []struct {
			InventoryName string `json:"inventory_name"`
			Items         []struct {
				Name     string `json:"name"`
				Quantity int    `json:"quantity"`
			} `json:"items"`
		} `json:"inventories"`
	}

	for _, customer := range customers {
		customerData := struct {
			CustomerName string `json:"customer_name"`
			Inventories  []struct {
				InventoryName string `json:"inventory_name"`
				Items         []struct {
					Name     string `json:"name"`
					Quantity int    `json:"quantity"`
				} `json:"items"`
			} `json:"inventories"`
		}{
			CustomerName: customer.Name,
		}

		for _, inventory := range customer.Inventories {
			inventoryData := struct {
				InventoryName string `json:"inventory_name"`
				Items         []struct {
					Name     string `json:"name"`
					Quantity int    `json:"quantity"`
				} `json:"items"`
			}{
				InventoryName: inventory.InventoryName,
			}

			for _, item := range inventory.Items {
				inventoryData.Items = append(inventoryData.Items, struct {
					Name     string `json:"name"`
					Quantity int    `json:"quantity"`
				}{
					Name:     item.Name,
					Quantity: item.Quantity,
				})
			}

			customerData.Inventories = append(customerData.Inventories, inventoryData)
		}

		response = append(response, customerData)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
