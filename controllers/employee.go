package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mohamedhabas11/golang-api/database"
	"github.com/mohamedhabas11/golang-api/models"
	"github.com/mohamedhabas11/golang-api/utils"
	"gorm.io/gorm"
)

// CreateEmployee registers a new shop employee.
func CreateEmployee(c *fiber.Ctx) error {
	var employee models.ShopEmployee
	if err := c.BodyParser(&employee); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Validate email format.
	if !utils.ValidateEmail(employee.Email) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid email format")
	}

	// Check if an employee with the same email already exists.
	var existing models.ShopEmployee
	if err := database.DB.Where("email = ?", employee.Email).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).SendString("Employee with this email already exists")
	} else if err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).SendString("Database error")
	}

	// Ensure the associated shop exists.
	var shop models.Shop
	if err := database.DB.First(&shop, employee.ShopID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Associated shop not found")
	}

	// Validate employee password (using minimum length of 8).
	if err := utils.ValidatePassword(employee.Password, utils.NewPasswordValidationConfig(8)); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid password: " + err.Error())
	}

	// Hash the password.
	hashedPassword, err := utils.HashPassword(employee.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error hashing password")
	}
	employee.Password = hashedPassword

	// Create the employee.
	if err := database.DB.Create(&employee).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(employee)
}

// GetEmployees retrieves all shop employees.
func GetEmployees(c *fiber.Ctx) error {
	var employees []models.ShopEmployee
	if err := database.DB.Find(&employees).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Prepare a response without sensitive info.
	type EmployeeResponse struct {
		ID        uint      `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		ShopID    uint      `json:"shop_id"`
	}

	var response []EmployeeResponse
	for _, emp := range employees {
		response = append(response, EmployeeResponse{
			ID:        emp.ID,
			CreatedAt: emp.CreatedAt,
			UpdatedAt: emp.UpdatedAt,
			Name:      emp.Name,
			Email:     emp.Email,
			ShopID:    emp.ShopID,
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetEmployee retrieves a single shop employee by ID.
func GetEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	var employee models.ShopEmployee
	if err := database.DB.First(&employee, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Employee not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Hide the password in the response.
	employee.Password = ""
	return c.Status(fiber.StatusOK).JSON(employee)
}

// UpdateEmployee updates an existing shop employee.
func UpdateEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	var employee models.ShopEmployee

	// Find the employee.
	if err := database.DB.First(&employee, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Employee not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	var updateData models.ShopEmployee
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Update allowed fields.
	if updateData.Name != "" {
		employee.Name = updateData.Name
	}
	if updateData.Email != "" && updateData.Email != employee.Email {
		if !utils.ValidateEmail(updateData.Email) {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid email format")
		}
		// Check if new email is already taken.
		var exists models.ShopEmployee
		if err := database.DB.Where("email = ?", updateData.Email).First(&exists).Error; err == nil {
			return c.Status(fiber.StatusConflict).SendString("Another employee with this email already exists")
		}
		employee.Email = updateData.Email
	}

	// Update password if provided.
	if updateData.Password != "" {
		if err := utils.ValidatePassword(updateData.Password, utils.NewPasswordValidationConfig(8)); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid password: " + err.Error())
		}
		hashedPassword, err := utils.HashPassword(updateData.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error hashing password")
		}
		employee.Password = hashedPassword
	}

	// Optionally update ShopID.
	if updateData.ShopID != 0 && updateData.ShopID != employee.ShopID {
		var shop models.Shop
		if err := database.DB.First(&shop, updateData.ShopID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("New associated shop not found")
		}
		employee.ShopID = updateData.ShopID
	}

	if err := database.DB.Save(&employee).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Hide the password in the response.
	employee.Password = ""
	return c.Status(fiber.StatusOK).JSON(employee)
}

// DeleteEmployee deletes a shop employee by ID.
func DeleteEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	var employee models.ShopEmployee

	if err := database.DB.First(&employee, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Employee not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if err := database.DB.Delete(&employee).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).SendString("Employee deleted successfully")
}
