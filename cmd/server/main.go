package main

import (
	"log"

	"github.com/WaveCE29/product_order_system/internal/infrastructure/config"
	database "github.com/WaveCE29/product_order_system/internal/infrastructure/db"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Load configuration
	config := config.Load()
	if config == nil {
		log.Fatal("Failed to load configuration")
	}

	// Initialize database connection
	db, err := database.NewSQLiteConnection(config.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Fatal(app.Listen(":" + config.Port))
}
