package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/WaveCE29/product_order_system/internal/application/port/input"
	"github.com/WaveCE29/product_order_system/internal/domain/entity"
	"github.com/WaveCE29/product_order_system/internal/domain/repository"
	"github.com/WaveCE29/product_order_system/pkg/logger"
)

type orderUseCase struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	logger      logger.Logger
}

// CreateOrder implements input.OrderUseCase.
func (o *orderUseCase) CreateOrder(ctx context.Context, req input.CreateOrderRequest) (*entity.Order, error) {
	o.logger.Info("Creating new order",
		"product_id", req.ProductID,
		"user_id", req.UserID,
		"quantity", req.Quantity,
		"idempotency_key", req.IdempotencyKey)

	// Check for existing order with same idempotency key
	existingOrder, err := o.orderRepo.GetByIdempotencyKey(ctx, req.IdempotencyKey)
	if err != nil && err != sql.ErrNoRows {
		o.logger.Error("Failed to check idempotency key", "error", err)
		return nil, fmt.Errorf("failed to check idempotency key: %w", err)
	}
	if existingOrder != nil {
		o.logger.Info("Order already exists with idempotency key", "order_id", existingOrder.ID)
		return existingOrder, nil
	}

	// Get product to check stock
	product, err := o.productRepo.GetbyID(ctx, req.ProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			o.logger.Error("Product not found", "product_id", req.ProductID)
			return nil, fmt.Errorf("product with id %d not found", req.ProductID)
		}
		o.logger.Error("Failed to get product", "product_id", req.ProductID, "error", err)
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Check if enough stock available
	if product.Stock < req.Quantity {
		o.logger.Warn("Insufficient stock",
			"product_id", req.ProductID,
			"available", product.Stock,
			"requested", req.Quantity)
		return nil, fmt.Errorf("insufficient stock: available %d, requested %d", product.Stock, req.Quantity)
	}

	// Create order
	order := &entity.Order{
		ProductID:      req.ProductID,
		UserID:         req.UserID,
		Quantity:       req.Quantity,
		Status:         entity.OrderStatusPending,
		IdempotencyKey: req.IdempotencyKey,
		CreatedAt:      time.Now(),
	}

	if err := o.orderRepo.Create(ctx, order); err != nil {
		o.logger.Error("Failed to create order", "error", err)
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Update product stock
	newStock := product.Stock - req.Quantity
	if err := o.productRepo.UpdateStock(ctx, req.ProductID, newStock); err != nil {
		o.logger.Error("Failed to update product stock", "product_id", req.ProductID, "error", err)
		return nil, fmt.Errorf("failed to update product stock: %w", err)
	}

	o.logger.Info("Order created successfully",
		"order_id", order.ID,
		"product_id", req.ProductID,
		"new_stock", newStock)

	return order, nil

}

func NewOrderUseCase(orderRepo repository.OrderRepository, productRepo repository.ProductRepository, logger logger.Logger) input.OrderUseCase {
	return &orderUseCase{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		logger:      logger,
	}

}
