package input

import (
	"context"

	"github.com/WaveCE29/product_order_system/internal/domain/entity"
)

type ProductUseCase interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (*entity.Product, error)
	GetProduct(ctx context.Context, id int) (*entity.Product, error)
	GetAllProduct(ctx context.Context) ([]*entity.Product, error)
}

type CreateProductRequest struct {
	Name  string `json:"name" validate:"required"`
	Stock int    `json:"stock" validate:"required"`
}
