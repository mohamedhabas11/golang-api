package models

import "gorm.io/gorm"

// Shop represents the business entity with its inventories and staff.
type Shop struct {
	gorm.Model
	Name        string         `json:"name"`
	Email       string         `json:"email" gorm:"unique"`
	OwnerID     uint           `json:"owner_id"`    // Reference to the shop owner.
	Owner       ShopOwner      `json:"owner"`       // One-to-one relation.
	Employees   []ShopEmployee `json:"employees"`   // One-to-many relation.
	Inventories []Inventory    `json:"inventories"` // One-to-many relation.
}

// ShopOwner represents a user who can manage the shop (and its employees/inventories).
type ShopOwner struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	// Additional owner-specific fields can be added here.
}

// ShopEmployee represents a user working within the shop’s scope.
type ShopEmployee struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	ShopID   uint   `json:"shop_id"` // Foreign key to the Shop.
}

// Customer represents a simple user who browses and buys items.
type Customer struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	// Customer-specific fields, like shipping address, can be added.
}

// Inventory represents a shop's collection of items.
type Inventory struct {
	gorm.Model
	ShopID        uint   `json:"shop_id"` // Foreign key to Shop.
	InventoryName string `json:"inventory_name"`
	Items         []Item `json:"items"` // One-to-many relationship.
}

// Item represents a product in an inventory.
type Item struct {
	gorm.Model
	InventoryID uint   `json:"inventory_id"` // Foreign key to Inventory.
	Name        string `json:"name"`
	Quantity    int    `json:"quantity"`
}
