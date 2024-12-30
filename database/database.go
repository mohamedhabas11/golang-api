package database

import (
	"fmt"
	"log"
	"os"

	"github.com/mohamedhabas11/golang-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBinstance represents a database instance with a direct gorm.DB embedding
type DBinstance struct {
	*gorm.DB
}

// Global DB instance
var DB DBinstance

// ConnectDB establishes a connection to the PostgreSQL database
func ConnectDB() {
	// Get DataBase connection info from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbSslmode := os.Getenv("DB_SSLMODE")
	dbTimezone := os.Getenv("DB_TIMEZONE")

	// Default values if missing
	if dbHost == "" {
		dbHost = "db"
		log.Println("DB_HOST is not set. Using default value: db")
	}
	if dbPort == "" {
		dbPort = "5432"
		log.Println("DB_PORT is not set. Using default value: 5432")
	}
	if dbTimezone == "" {
		dbTimezone = "UTC"
		log.Println("DB_TIMEZONE is not set. Using default value: UTC")
	}
	if dbSslmode == "" {
		dbSslmode = "disable"
		log.Println("DB_SSLMODE is not set. Using default value: disable")
	}

	// Build the Data Source Name (DSN) for PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSslmode, dbTimezone,
	)

	// Open a connection to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Enable detailed logging of queries
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set the global DB instance
	DB = DBinstance{DB: db}

	// Run database migrations (for automatic schema generation)
	log.Println("Running migrations...")
	if err := db.AutoMigrate(&models.Customer{}, &models.User{}, &models.Inventory{}, &models.Item{}); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}
	log.Println("Migrations completed.")
}
