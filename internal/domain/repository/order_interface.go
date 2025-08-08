package repository

import (
	"context"

	"github.com/WaveCE29/product_order_system/internal/domain/entity"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	GetByID(ctx context.Context, id int) (*entity.Order, error)
	GetByIdempotencyKey(ctx context.Context, key string) (*entity.Order, error)
	GetAll(ctx context.Context) ([]*entity.Order, error)
}
