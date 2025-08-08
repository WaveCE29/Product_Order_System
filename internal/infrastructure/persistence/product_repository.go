package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/WaveCE29/product_order_system/internal/domain/entity"
	"github.com/WaveCE29/product_order_system/internal/domain/repository"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) repository.ProductRepository {
	return &productRepository{db: db}
}

// Create implements repository.ProductRepository.
func (p *productRepository) Create(ctx context.Context, product *entity.Product) error {
	query := `
		INSERT INTO products (name, stock, created_at, updated_at) 
		VALUES (?, ?, ?, ?)
	`

	result, err := p.db.ExecContext(ctx, query,
		product.Name,
		product.Stock,
		product.CreatedAt,
		product.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	product.ID = int(id)
	return nil
}

// GetAll implements repository.ProductRepository.
func (p *productRepository) GetAll(ctx context.Context) ([]*entity.Product, error) {
	query := `SELECT * FROM products ORDER BY created_at DESC`

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}

	defer rows.Close()

	var products []*entity.Product

	for rows.Next() {
		var product entity.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Stock,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

// GetbyID implements repository.ProductRepository.
func (p *productRepository) GetbyID(ctx context.Context, id int) (*entity.Product, error) {
	query := `SELECT id, name, stock, created_at, updated_at FROM products WHERE id = ?`
	var product entity.Product
	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return &product, nil
}

// Update implements repository.ProductRepository.
func (p *productRepository) Update(ctx context.Context, product *entity.Product) error {
	query := `
		UPDATE products 
		SET name = ?, stock = ?, updated_at = ? 
		WHERE id = ?
	`
	product.UpdatedAt = time.Now()

	result, err := p.db.ExecContext(ctx, query,
		product.Name,
		product.Stock,
		product.UpdatedAt,
		product.ID)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found", product.ID)
	}

	return nil

}

// UpdateStock implements repository.ProductRepository.
func (p *productRepository) UpdateStock(ctx context.Context, productID int, newStock int) error {
	query := `
		UPDATE products 
		SET stock = ?, updated_at = ? 
		WHERE id = ?
	`
	result, err := p.db.ExecContext(ctx, query, newStock, time.Now(), productID)
	if err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found", productID)
	}

	return nil

}
