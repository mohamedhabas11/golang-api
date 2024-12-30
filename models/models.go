package models

import "gorm.io/gorm"

// Customer represents a customer in the system, responsible for managing their own inventory.
type Customer struct {
	gorm.Model
	Name        string      `json:"name"`
	Email       string      `json:"email" gorm:"unique"` // Unique constraint for the customer email
	Inventories []Inventory `json:"inventories"`         // One-to-many relationship with Inventory
	Users       []User      `json:"users"`               // One-to-many relationship with Users (customers can have multiple users)
}

// User represents a user in the system, who can be associated with a customer.
type User struct {
	gorm.Model
	Name       string   `json:"name"`
	Email      string   `json:"email" gorm:"unique"` // Unique constraint for user email
	Password   string   `json:"password"`
	CustomerID uint     `json:"customer_id"` // Foreign Key to link the user to a customer
	Customer   Customer `json:"customer"`    // Relationship to Customer
}

// Inventory represents an inventory owned and managed by a customer.
type Inventory struct {
	gorm.Model
	CustomerID    uint     `json:"customer_id"` // Foreign Key to Customer
	InventoryName string   `json:"inventory_name"`
	Items         []Item   `json:"items"`    // One-to-many relationship with Item
	Customer      Customer `json:"customer"` // Relationship to Customer
}

// Item represents a unique item in an inventory.
type Item struct {
	gorm.Model
	InventoryID uint      `json:"inventory_id"` // Foreign Key to Inventory
	Name        string    `json:"name"`
	Quantity    int       `json:"quantity"`
	Inventory   Inventory `json:"inventory"` // Relationship to Inventory
}
