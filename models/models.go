package models

import "gorm.io/gorm"

// Customer represents a customer in the system responsible for managing thier own inventory
type Costumer struct {
	gorm.Model
	Name        string      `json:"name"`
	Inventories []Inventory `json:"inventories"` // One-to-many relationship with Inventory
}

// Inventory represent an inventory owned and managed by a customer
type Inventory struct {
	gorm.Model
	CustomerID    uint   `json:"customer_id"` // Foreign Key to Customer
	InventoryName string `json:"inventory_name"`
	Items         []Item `json:"items"` // One-to-many relationship with Item
}

// Item represents a unique item in an inventory
type Item struct {
	gorm.Model
	InventoryID uint   `json:"inventory_id"` // Foreign Key to Inventory
	Name        string `json:"name"`
	Quantity    int    `json:"quantity"`
}
