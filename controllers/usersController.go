package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mohamedhabas11/golang-api/database"
	"github.com/mohamedhabas11/golang-api/middlewares"
	"github.com/mohamedhabas11/golang-api/models"
	"github.com/mohamedhabas11/golang-api/utils"
	"gorm.io/gorm"
)

// User password config
var userPasswordConfig = utils.NewPasswordValidationConfig(8) // Set the minimum length to 8 or your desired value

// CreateUser creates a new user
func CreateUser(c *fiber.Ctx) error {
	var user models.User

	// Parse request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Validate email and password using the config
	if !utils.ValidateEmail(user.Email) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid email format")
	}

	// Check password validation result
	if err := utils.ValidatePassword(user.Password, userPasswordConfig); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid password: " + err.Error())
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

// GetUser retrieves a user by ID
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
		// Check password validation result
		if err := utils.ValidatePassword(updates.Password, userPasswordConfig); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid password: " + err.Error())
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

func Login(c *fiber.Ctx) error {

	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Validate email format
	if !utils.ValidateEmail(req.Email) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid email format")
	}

	// Lookup requested User in the database
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Database error")
	}

	// Compare sent password against saved user pass hash
	if !utils.ComparePassword(existingUser.Password, req.Password) {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
	}

	// Generate JWT token with default expiration (24 hours)
	token, err := middlewares.GenerateJWT(existingUser.ID, existingUser.Email, time.Hour*24) // Default expiration of 24 hours
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error generating token")
	}

	// Set the token in a cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt_token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24), // Cookie expires in 24 hours
		HTTPOnly: true,                           // Security: makes the cookie accessible only via HTTP (not JavaScript)
		Secure:   false,                          // Set to true if using HTTPS
	})

	// Send generated token back to the client
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
	})
}
