package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/mohamedhabas11/golang-api/database"
	"github.com/mohamedhabas11/golang-api/routes"
)

func main() {
	// Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

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

	// Start the server in a goroutine to handle graceful shutdown
	go func() {
		log.Printf("Server started on port %s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Set up signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit // Block until a signal is received
	log.Println("Shutting down server...")

	// Gracefully shutdown the Fiber app
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
