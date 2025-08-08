package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/WaveCE29/product_order_system/internal/domain/entity"
	"github.com/WaveCE29/product_order_system/internal/domain/repository"
)

type orderRepository struct {
	db *sql.DB
}

// Create implements repository.OrderRepository.
func (o *orderRepository) Create(ctx context.Context, order *entity.Order) error {
	query := `
		INSERT INTO orders (product_id, user_id, quantity, status, idempotency_key, created_at) 
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := o.db.ExecContext(ctx, query,
		order.ProductID,
		order.UserID,
		order.Quantity,
		order.Status,
		order.IdempotencyKey,
		order.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	order.ID = int(id)
	return nil
}

// GetAll implements repository.OrderRepository.
func (o *orderRepository) GetAll(ctx context.Context) ([]*entity.Order, error) {
	query := `
		SELECT id, product_id, user_id, quantity, status, idempotency_key, created_at 
		FROM orders 
		ORDER BY created_at DESC
	`

	rows, err := o.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	var orders []*entity.Order
	for rows.Next() {
		var order entity.Order
		err := rows.Scan(
			&order.ID,
			&order.ProductID,
			&order.UserID,
			&order.Quantity,
			&order.Status,
			&order.IdempotencyKey,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating orders: %w", err)
	}

	return orders, nil
}

// GetByID implements repository.OrderRepository.
func (o *orderRepository) GetByID(ctx context.Context, id int) (*entity.Order, error) {
	query := `
		SELECT id, product_id, user_id, quantity, status, idempotency_key, created_at 
		FROM orders 
		WHERE id = ?
	`

	var order entity.Order
	err := o.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.ProductID,
		&order.UserID,
		&order.Quantity,
		&order.Status,
		&order.IdempotencyKey,
		&order.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order, nil
}

// GetByIdempotencyKey implements repository.OrderRepository.
func (o *orderRepository) GetByIdempotencyKey(ctx context.Context, key string) (*entity.Order, error) {
	query := `
		SELECT id, product_id, user_id, quantity, status, idempotency_key, created_at 
		FROM orders 
		WHERE idempotency_key = ?
	`

	var order entity.Order
	err := o.db.QueryRowContext(ctx, query, key).Scan(
		&order.ID,
		&order.ProductID,
		&order.UserID,
		&order.Quantity,
		&order.Status,
		&order.IdempotencyKey,
		&order.CreatedAt,
	)
	if err != nil {
		return nil, err // Return sql.ErrNoRows as is for idempotency check
	}

	return &order, nil
}

func NewOrderRepository(db *sql.DB) repository.OrderRepository {
	return &orderRepository{db: db}
}
