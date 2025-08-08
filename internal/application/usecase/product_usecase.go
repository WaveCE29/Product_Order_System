package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/WaveCE29/product_order_system/internal/application/port/input"
	"github.com/WaveCE29/product_order_system/internal/domain/entity"
	"github.com/WaveCE29/product_order_system/internal/domain/repository"
	"github.com/WaveCE29/product_order_system/pkg/logger"
)

type productUseCase struct {
	productRepo repository.ProductRepository
	logger      logger.Logger
}

// GetAllProduct implements input.ProductUseCase.
func (p *productUseCase) GetAllProduct(ctx context.Context) ([]*entity.Product, error) {
	p.logger.Info("Getting all products")

	products, err := p.productRepo.GetAll(ctx)
	if err != nil {
		p.logger.Error("Failed to get products", "error", err)
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	p.logger.Info("Retrieved products", "count", len(products))
	return products, nil
}

// CreateProduct implements input.ProductUseCase.
func (p *productUseCase) CreateProduct(ctx context.Context, req input.CreateProductRequest) (*entity.Product, error) {
	p.logger.Info("Creating new product", "name", req.Name, "stock", req.Stock)

	product := &entity.Product{
		Name:      req.Name,
		Stock:     req.Stock,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := p.productRepo.Create(ctx, product); err != nil {
		p.logger.Error("Failed to create product", "error", err)
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	p.logger.Info("Product created successfully", "id", product.ID)

	return product, nil
}

// GetProduct implements input.ProductUseCase.
func (p *productUseCase) GetProduct(ctx context.Context, id int) (*entity.Product, error) {
	p.logger.Info("Getting product", "id", id)

	product, err := p.productRepo.GetbyID(ctx, id)
	if err != nil {
		p.logger.Error("Failed to get product", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

func NewProductUseCase(productRepo repository.ProductRepository, logger logger.Logger) input.ProductUseCase {
	return &productUseCase{
		productRepo: productRepo,
		logger:      logger,
	}

}
