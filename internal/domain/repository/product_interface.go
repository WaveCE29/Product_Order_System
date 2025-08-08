package repository

import (
	"context"

	"github.com/WaveCE29/product_order_system/internal/domain/entity"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	GetbyID(ctx context.Context, id int) (*entity.Product, error)
	GetAll(ctx context.Context) ([]*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
	UpdateStock(ctx context.Context, productID int, newStock int) error
}
