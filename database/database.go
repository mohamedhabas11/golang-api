package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mohamedhabas11/golang-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBinstance struct {
	Db *gorm.DB
}

var DB DBinstance

// ConnectDB establishes a connection to the PostgreSQL database
func ConnectDB() {
	// Load enviorment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get DataBase connection info from enviorment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbSslmode := os.Getenv("DB_SSLMODE")
	dbTimezone := os.Getenv("DB_TIMEZONE")

	// Default values if messing
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbTimezone == "" {
		dbTimezone = "UTC"
	}
	if dbSslmode == "" {
		dbSslmode = "disable"
	}

	// Build the Data Source Name (DSN) for PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSslmode, dbTimezone,
	)

	// open a connection to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // enable detailed logging of queries
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set the global DB instance
	DB = DBinstance{Db: db}

	// Run database migrations (for automatic schema generation) // TODO: further investigate
	log.Println("Running migrations...")
	if err := db.AutoMigrate(&models.Customer{}, &models.Inventory{}, &models.Item{}); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	log.Println("Migrations completed.")
}
