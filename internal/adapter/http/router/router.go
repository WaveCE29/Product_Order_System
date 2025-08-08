package router

import (
	"github.com/WaveCE29/product_order_system/internal/adapter/http/handler"
	"github.com/WaveCE29/product_order_system/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func SetupRoutes(app *fiber.App, h *handler.Handler, logger logger.Logger) {
	// Middleware
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "*",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "Product Order System",
		})
	})

	// API routes
	api := app.Group("/api/v1")

	// Product routes
	products := api.Group("/products")
	products.Post("/", h.CreateProduct)
	products.Get("/", h.GetAllProducts)
	products.Get("/:id", h.GetProduct)

	// Order routes
	orders := api.Group("/orders")
	orders.Post("/", h.CreateOrder)

	// Legacy routes (without /api/v1 prefix for compatibility)
	app.Post("/products", h.CreateProduct)
	app.Get("/products", h.GetAllProducts)
	app.Get("/products/:id", h.GetProduct)
	app.Post("/orders", h.CreateOrder)

	logger.Info("Routes configured successfully")
}
