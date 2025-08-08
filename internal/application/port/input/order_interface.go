package input

import (
	"context"

	"github.com/WaveCE29/product_order_system/internal/domain/entity"
)

type OrderUseCase interface {
	CreateOrder(ctx context.Context, req CreateOrderRequest) (*entity.Order, error)
}

type CreateOrderRequest struct {
	ProductID      int    `json:"product_id" validate:"required"`
	UserID         string `json:"user_id" validate:"required"`
	Quantity       int    `json:"quantity" validate:"required"`
	IdempotencyKey string `json:"idempotency_key" validate:"required"`
}
