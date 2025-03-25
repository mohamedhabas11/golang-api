package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/mohamedhabas11/golang-api/database"
	"github.com/mohamedhabas11/golang-api/middlewares"
	"github.com/mohamedhabas11/golang-api/models"
	"github.com/mohamedhabas11/golang-api/utils"
)

var (
	passwordValidationLength int
)

// CreateCustomer registers a new customer.
func CreateCustomer(c *fiber.Ctx) error {
	var customer models.Customer

	// set default passwordValidationLength
	if passwordValidationLength == 0 {
		passwordValidationLength = 8
	}

	// Parse request body
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Validate email format
	if !utils.ValidateEmail(customer.Email) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid email format")
	}

	// Validate password (assuming a minimum length of 8)
	if err := utils.ValidatePassword(customer.Password, utils.NewPasswordValidationConfig(passwordValidationLength)); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid password: " + err.Error())
	}

	// Check if customer already exists
	var existing models.Customer
	if err := database.DB.Where("email = ?", customer.Email).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).SendString("Customer with this email already exists")
	} else if err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).SendString("Database error")
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(customer.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error hashing password")
	}
	customer.Password = hashedPassword

	// Save the new customer
	if err := database.DB.Create(&customer).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(customer)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginCustomer authenticates a customer and issues a JWT token.
func LoginCustomer(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Validate email format
	if !utils.ValidateEmail(req.Email) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid email format")
	}

	// Lookup customer in the database
	var customer models.Customer
	if err := database.DB.Where("email = ?", req.Email).First(&customer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Database error")
	}

	// Compare provided password with hashed password
	if !utils.ComparePassword(customer.Password, req.Password) {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
	}

	// Generate a JWT token (expires in 24 hours)
	token, err := middlewares.GenerateJWT(customer.ID, customer.Email, time.Hour*24)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error generating token")
	}

	// Set the token in a secure HTTP-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt_token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Secure:   false, // change to true if using HTTPS
	})

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
	})
}

// GetCustomers retrieves all customers with selected fields.
func GetCustomers(c *fiber.Ctx) error {
	// Define a lightweight struct to return only the necessary fields.
	var customers []struct {
		ID        uint      `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
	}

	// Query only selected fields from the Customer model.
	if err := database.DB.Model(&models.Customer{}).
		Select("id, created_at, updated_at, name, email").
		Scan(&customers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(customers)
}
