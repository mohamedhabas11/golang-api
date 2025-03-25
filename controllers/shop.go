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

// CreateShopRequest represents the JSON payload for creating a shop and its owner.
type CreateShopRequest struct {
	Name  string           `json:"name"`
	Email string           `json:"email"`
	Owner models.ShopOwner `json:"owner"`
	// You can extend this struct to accept initial employees or inventories if needed.
}

// CreateShop creates a new shop along with its owner.
func CreateShop(c *fiber.Ctx) error {
	var req CreateShopRequest

	// Parse the incoming request body.
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Validate shop email.
	if !utils.ValidateEmail(req.Email) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid shop email format")
	}

	// Validate owner email.
	if !utils.ValidateEmail(req.Owner.Email) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid owner email format")
	}

	// Validate owner's password.
	if err := utils.ValidatePassword(req.Owner.Password, utils.NewPasswordValidationConfig(8)); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid owner password: " + err.Error())
	}

	// Check if a shop with the same email already exists.
	var existingShop models.Shop
	if err := database.DB.Where("email = ?", req.Email).First(&existingShop).Error; err == nil {
		return c.Status(fiber.StatusConflict).SendString("Shop with this email already exists")
	}

	// Create or retrieve the shop owner.
	var owner models.ShopOwner
	if err := database.DB.Where("email = ?", req.Owner.Email).First(&owner).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Hash the owner's password.
			hashedPassword, err := utils.HashPassword(req.Owner.Password)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Error hashing owner password")
			}
			req.Owner.Password = hashedPassword

			// Create the shop owner.
			if err := database.DB.Create(&req.Owner).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Error creating shop owner")
			}
			owner = req.Owner
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString("Error checking shop owner")
		}
	}

	// Create the shop with the owner's ID.
	shop := models.Shop{
		Name:    req.Name,
		Email:   req.Email,
		OwnerID: owner.ID,
	}
	if err := database.DB.Create(&shop).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error creating shop")
	}

	// Optionally, load the owner details in the shop response.
	if err := database.DB.Preload("Owner").First(&shop, shop.ID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error loading created shop")
	}

	return c.Status(fiber.StatusCreated).JSON(shop)
}

// UpdateShop updates existing shop
func UpdateShop(c *fiber.Ctx) error {
	// Get the shop ID from the URL parameters
	id := c.Params("id")
	var shop models.Shop

	// Retrieve the shop record by ID.
	if err := database.DB.First(&shop, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Shop not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Define a payload for allowed updates.
	var updateData struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// Parse the request body.
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Update Name if provided and different.
	if updateData.Name != "" && updateData.Name != shop.Name {
		shop.Name = updateData.Name
	}

	// Update Email if provided and different.
	if updateData.Email != "" && updateData.Email != shop.Email {
		if !utils.ValidateEmail(updateData.Email) {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid email format")
		}
		// Check if the new email is already used by another shop.
		var existingShop models.Shop
		if err := database.DB.Where("email = ?", updateData.Email).First(&existingShop).Error; err == nil && existingShop.ID != shop.ID {
			return c.Status(fiber.StatusConflict).SendString("Shop with this email already exists")
		}

		shop.Email = updateData.Name
	}

	// save update to the database.
	if err := database.DB.Save(&shop).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// reload after update
	if err := database.DB.Preload("Owner").First(&shop, shop.ID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error preloading") // TODO: log on disk
	}

	return c.Status(fiber.StatusOK).JSON(shop)

}

// GetShops retrieves all shops with their owners.
func GetShops(c *fiber.Ctx) error {
	var shops []models.Shop

	// Preload the Owner relation.
	if err := database.DB.Preload("Owner").Find(&shops).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(shops)
}

// LoginShopOwner authenticates a shop owner and returns a JWT token.
func LoginShopOwner(c *fiber.Ctx) error {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Validate email format.
	if !utils.ValidateEmail(req.Email) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid email format")
	}

	// Lookup shop owner in the database.
	var owner models.ShopOwner
	if err := database.DB.Where("email = ?", req.Email).First(&owner).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Database error")
	}

	// Compare the provided password with the stored hashed password.
	if !utils.ComparePassword(owner.Password, req.Password) {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
	}

	// Generate a JWT token (here, valid for 24 hours).
	token, err := middlewares.GenerateJWT(owner.ID, owner.Email, time.Hour*24)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error generating token")
	}

	// Set the token in an HTTP-only cookie.
	c.Cookie(&fiber.Cookie{
		Name:     "jwt_token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Secure:   false, // Set to true in production if using HTTPS.
	})

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
	})
}
