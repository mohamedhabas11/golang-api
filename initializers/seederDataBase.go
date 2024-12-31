package initializers

import (
	"encoding/json"
	"log"
	"os"

	"github.com/mohamedhabas11/golang-api/database"
	"github.com/mohamedhabas11/golang-api/models"
)

func SeedDatabase(seedFile string) {
	// Read seed data from JSON file
	file, err := os.ReadFile(seedFile)
	if err != nil {
		log.Fatal("Error reading seed file: ", err)
	}

	var seedData struct {
		Customers []models.Customer `json:"customers"`
	}
	if err := json.Unmarshal(file, &seedData); err != nil {
		log.Fatal("Error unmarshalling seed data: ", err)
	}

	// Insert all customers, their users, inventories, and items
	for _, customer := range seedData.Customers {
		// Create customer and related fields (users, inventories, items)
		if err := database.DB.Create(&customer).Error; err != nil {
			log.Printf("Error creating customer %s: %v", customer.Name, err)
			continue
		}
	}
}
