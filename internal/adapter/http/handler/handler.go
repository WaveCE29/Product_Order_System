package handler

import (
	"strconv"

	"github.com/WaveCE29/product_order_system/internal/application/port/input"
	"github.com/WaveCE29/product_order_system/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	productUseCase input.ProductUseCase
	orderUseCase   input.OrderUseCase
	logger         logger.Logger
}

func NewHandler(productUseCase input.ProductUseCase, orderUseCase input.OrderUseCase, logger logger.Logger) *Handler {
	return &Handler{
		productUseCase: productUseCase,
		orderUseCase:   orderUseCase,
		logger:         logger,
	}
}

// Product handlers
func (h *Handler) CreateProduct(c *fiber.Ctx) error {
	var req input.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product name is required",
		})
	}

	if req.Stock < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Stock must be non-negative",
		})
	}

	product, err := h.productUseCase.CreateProduct(c.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create product", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Product created successfully",
		"data":    product,
	})
}

func (h *Handler) GetProduct(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	product, err := h.productUseCase.GetProduct(c.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get product", "id", id, "error", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Product retrieved successfully",
		"data":    product,
	})
}

func (h *Handler) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.productUseCase.GetAllProduct(c.Context())
	if err != nil {
		h.logger.Error("Failed to get products", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve products",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Products retrieved successfully",
		"data":    products,
		"count":   len(products),
	})
}

// Order handlers
func (h *Handler) CreateOrder(c *fiber.Ctx) error {
	var req input.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Generate idempotency key if not provided
	if req.IdempotencyKey == "" {
		req.IdempotencyKey = uuid.New().String()
	}

	// Validate request
	if req.ProductID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Valid product ID is required",
		})
	}

	if req.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	if req.Quantity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Quantity must be greater than 0",
		})
	}

	order, err := h.orderUseCase.CreateOrder(c.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create order", "error", err)

		// Check for specific error types
		errMsg := err.Error()
		if contains(errMsg, "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": errMsg,
			})
		}
		if contains(errMsg, "insufficient stock") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": errMsg,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create order",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Order created successfully",
		"data":    order,
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
