package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mohamedhabas11/golang-api/database"
	"github.com/mohamedhabas11/golang-api/models"
	"github.com/mohamedhabas11/golang-api/utils"
	"gorm.io/gorm"
)

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
	token, err := utils.GenerateJWT(existingUser.ID, existingUser.Email, time.Hour*24) // Default expiration of 24 hours
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
