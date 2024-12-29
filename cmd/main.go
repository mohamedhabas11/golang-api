package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/mohamedhabas11/golang-api/database"
	"github.com/mohamedhabas11/golang-api/routes"
)

func main() {
	// Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize the database connection
	database.ConnectDB()

	// Create a new Fiber app
	app := fiber.New()

	// Set up the routes
	routes.SetupRoutes(app)

	// Get the port from the environment variable or use default port 3000
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000" // Default port if APP_PORT is not set
	}

	// Start the server
	log.Printf("Server started on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
