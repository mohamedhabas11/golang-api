package initializers

import (
	"encoding/json"
	"log"
	"os"

	"github.com/mohamedhabas11/golang-api/database"
	"github.com/mohamedhabas11/golang-api/models"
	"github.com/mohamedhabas11/golang-api/utils"
	"gorm.io/gorm"
)

func SeedDatabase(seedFile string) {
	// Read seed data from JSON file.
	file, err := os.ReadFile(seedFile)
	if err != nil {
		log.Fatal("Error reading seed file: ", err)
	}

	// Define a structure matching our JSON format.
	var seedData struct {
		Shops []struct {
			Name        string                `json:"name"`
			Email       string                `json:"email"`
			Owner       models.ShopOwner      `json:"owner"`
			Employees   []models.ShopEmployee `json:"employees"`
			Inventories []struct {
				InventoryName string        `json:"inventory_name"`
				Items         []models.Item `json:"items"`
			} `json:"inventories"`
		} `json:"shops"`
		Customers []models.Customer `json:"customers"`
	}

	if err := json.Unmarshal(file, &seedData); err != nil {
		log.Fatal("Error unmarshalling seed data: ", err)
	}

	// Seed Shops and their related data.
	for _, shopData := range seedData.Shops {
		// --- Seed Shop Owner ---
		// Always hash the owner's password.
		hashedOwnerPass, err := utils.HashPassword(shopData.Owner.Password)
		if err != nil {
			log.Printf("Error hashing password for shop owner %s: %v", shopData.Owner.Email, err)
			continue
		}
		shopData.Owner.Password = hashedOwnerPass

		// Use FirstOrCreate to ensure no duplicates
		var existingOwner models.ShopOwner
		result := database.DB.Where("email = ?", shopData.Owner.Email).First(&existingOwner)
		if result.Error == gorm.ErrRecordNotFound {
			// If owner doesn't exist, create
			if err := database.DB.Create(&shopData.Owner).Error; err != nil {
				log.Printf("Error seeding shop owner %s: %v", shopData.Owner.Email, err)
				continue
			}
		} else if result.Error != nil {
			log.Printf("Error checking shop owner %s: %v", shopData.Owner.Email, result.Error)
			continue
		} else {
			// If owner exists, update
			shopData.Owner.ID = existingOwner.ID
			if err := database.DB.Save(&shopData.Owner).Error; err != nil {
				log.Printf("Error updating shop owner %s: %v", shopData.Owner.Email, err)
				continue
			}
		}
		owner := shopData.Owner

		// --- Seed Shop ---
		shop := models.Shop{
			Name:    shopData.Name,
			Email:   shopData.Email,
			OwnerID: owner.ID,
		}
		var existingShop models.Shop
		result = database.DB.Where("email = ?", shopData.Email).First(&existingShop)
		if result.Error == gorm.ErrRecordNotFound {
			// If shop doesn't exist, create
			if err := database.DB.Create(&shop).Error; err != nil {
				log.Printf("Error seeding shop %s: %v", shopData.Name, err)
				continue
			}
		} else if result.Error != nil {
			log.Printf("Error checking shop %s: %v", shopData.Name, result.Error)
			continue
		} else {
			// If shop exists, update
			shop.ID = existingShop.ID
			if err := database.DB.Save(&shop).Error; err != nil {
				log.Printf("Error updating shop %s: %v", shopData.Name, err)
				continue
			}
		}

		// --- Seed Shop Employees ---
		for _, empData := range shopData.Employees {
			// Always hash the employee's password.
			hashedEmpPass, err := utils.HashPassword(empData.Password)
			if err != nil {
				log.Printf("Error hashing password for employee %s: %v", empData.Email, err)
				continue
			}
			empData.Password = hashedEmpPass
			empData.ShopID = shop.ID

			var existingEmployee models.ShopEmployee
			result := database.DB.Where("email = ? AND shop_id = ?", empData.Email, shop.ID).First(&existingEmployee)
			if result.Error == gorm.ErrRecordNotFound {
				// If employee doesn't exist, create
				if err := database.DB.Create(&empData).Error; err != nil {
					log.Printf("Error seeding employee %s: %v", empData.Email, err)
					continue
				}
			} else if result.Error != nil {
				log.Printf("Error checking employee %s: %v", empData.Email, result.Error)
				continue
			} else {
				// If employee exists, update
				empData.ID = existingEmployee.ID
				if err := database.DB.Save(&empData).Error; err != nil {
					log.Printf("Error updating employee %s: %v", empData.Email, err)
					continue
				}
			}
		}

		// --- Seed Inventories and Items ---
		for _, invData := range shopData.Inventories {
			// Upsert inventory record
			inv := models.Inventory{
				InventoryName: invData.InventoryName,
				ShopID:        shop.ID,
			}
			var existingInventory models.Inventory
			result := database.DB.Where("inventory_name = ? AND shop_id = ?", invData.InventoryName, shop.ID).First(&existingInventory)
			if result.Error == gorm.ErrRecordNotFound {
				// If inventory doesn't exist, create
				if err := database.DB.Create(&inv).Error; err != nil {
					log.Printf("Error seeding inventory %s: %v", invData.InventoryName, err)
					continue
				}
			} else if result.Error != nil {
				log.Printf("Error checking inventory %s: %v", invData.InventoryName, result.Error)
				continue
			} else {
				// If inventory exists, update
				inv.ID = existingInventory.ID
				if err := database.DB.Save(&inv).Error; err != nil {
					log.Printf("Error updating inventory %s: %v", invData.InventoryName, err)
					continue
				}
			}

			// Upsert each item for the inventory
			for _, itemData := range invData.Items {
				itemData.InventoryID = inv.ID
				var existingItem models.Item
				result := database.DB.Where("name = ? AND inventory_id = ?", itemData.Name, inv.ID).First(&existingItem)
				if result.Error == gorm.ErrRecordNotFound {
					// If item doesn't exist, create
					if err := database.DB.Create(&itemData).Error; err != nil {
						log.Printf("Error seeding item %s: %v", itemData.Name, err)
						continue
					}
				} else if result.Error != nil {
					log.Printf("Error checking item %s: %v", itemData.Name, result.Error)
					continue
				} else {
					// If item exists, update
					itemData.ID = existingItem.ID
					if err := database.DB.Save(&itemData).Error; err != nil {
						log.Printf("Error updating item %s: %v", itemData.Name, err)
						continue
					}
				}
			}
		}
	}

	// --- Seed Standalone Customers ---
	for _, custData := range seedData.Customers {
		// Always hash the customer's password.
		hashedCustPass, err := utils.HashPassword(custData.Password)
		if err != nil {
			log.Printf("Error hashing password for customer %s: %v", custData.Email, err)
			continue
		}
		custData.Password = hashedCustPass

		var existingCustomer models.Customer
		result := database.DB.Where("email = ?", custData.Email).First(&existingCustomer)
		if result.Error == gorm.ErrRecordNotFound {
			// If customer doesn't exist, create
			if err := database.DB.Create(&custData).Error; err != nil {
				log.Printf("Error seeding customer %s: %v", custData.Email, err)
				continue
			}
		} else if result.Error != nil {
			log.Printf("Error checking customer %s: %v", custData.Email, result.Error)
			continue
		} else {
			// If customer exists, update
			custData.ID = existingCustomer.ID
			if err := database.DB.Save(&custData).Error; err != nil {
				log.Printf("Error updating customer %s: %v", custData.Email, err)
				continue
			}
		}
	}
}
