package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohamedhabas11/golang-api/database"
	"github.com/mohamedhabas11/golang-api/models"
	"github.com/mohamedhabas11/golang-api/utils"
)

// CreateUser creates a new user
func CreateUser(c *fiber.Ctx) error {
	var user models.User

	// Parse request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Validate email and password
	if !utils.ValidateEmail(user.Email) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid email format")
	}
	if !utils.ValidatePassword(user.Password) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid password")
	}

	// Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).SendString("User with this email already exists")
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error hashing password")
	}
	user.Password = hashedPassword

	// Save user to the database
	if result := database.DB.Create(&user); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetUsers retrieves all users
func GetUsers(c *fiber.Ctx) error {
	var users []struct {
		ID        uint   `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Name      string `json:"name"`
		Email     string `json:"email"`
	}

	// Query only selected fields, excluding sensitive or nested data
	if result := database.DB.Model(&models.User{}).Select("id, created_at, updated_at, name, email").Scan(&users); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user struct {
		ID        uint   `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Name      string `json:"name"`
		Email     string `json:"email"`
	}

	// Query only selected fields by ID
	if err := database.DB.Model(&models.User{}).Select("id, created_at, updated_at, name, email").Where("id = ?", id).Scan(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("User not found")
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// UpdateUser updates a user by ID
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	// Find existing user
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("User not found")
	}

	// Parse request body
	var updates models.User
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Update fields conditionally
	if updates.Email != "" && updates.Email != user.Email {
		if !utils.ValidateEmail(updates.Email) {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid email format")
		}
		user.Email = updates.Email
	}

	if updates.Password != "" {
		if !utils.ValidatePassword(updates.Password) {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid password")
		}
		hashedPassword, err := utils.HashPassword(updates.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error hashing password")
		}
		user.Password = hashedPassword
	}

	if updates.Name != "" {
		user.Name = updates.Name
	}

	// Save updates
	if result := database.DB.Save(&user); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(result.Error.Error())
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// DeleteUser deletes a user by ID
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	// Check if the user exists
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("User not found")
	}

	// Delete the user
	if result := database.DB.Delete(&user); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error deleting user")
	}

	return c.Status(fiber.StatusOK).SendString("User deleted successfully")
}
