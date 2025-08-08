package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/WaveCE29/product_order_system/internal/adapter/http/handler"
	"github.com/WaveCE29/product_order_system/internal/adapter/http/router"
	"github.com/WaveCE29/product_order_system/internal/application/usecase"
	"github.com/WaveCE29/product_order_system/internal/infrastructure/config"
	database "github.com/WaveCE29/product_order_system/internal/infrastructure/db"
	"github.com/WaveCE29/product_order_system/internal/infrastructure/persistence"
	"github.com/WaveCE29/product_order_system/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	config := config.LoadConfig()
	if config == nil {
		log.Fatal("Failed to load configuration")
	}

	// Initialize logger
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	// Initialize database
	db, err := database.NewDatabase(config.Database.Path, logger)
	if err != nil {
		logger.Error("Failed to initialize database", "error", err)
		log.Fatal(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("Failed to close database", "error", err)
		}
	}()

	productRepo := persistence.NewProductRepository(db.DB)
	orderRepo := persistence.NewOrderRepository(db.DB)

	productUseCase := usecase.NewProductUseCase(productRepo, logger)
	orderUseCase := usecase.NewOrderUseCase(orderRepo, productRepo, logger)

	h := handler.NewHandler(productUseCase, orderUseCase, logger)

	app := fiber.New(fiber.Config{
		AppName: "Product Order System",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			logger.Error("Request error", "error", err, "path", c.Path(), "method", c.Method())

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	router.SetupRoutes(app, h, logger)

	go func() {
		address := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
		logger.Info("Server starting", "address", address)

		if err := app.Listen(address); err != nil {
			logger.Error("Failed to start server", "error", err)
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	logger.Info("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server shutdown completed")

}
