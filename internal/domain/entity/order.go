package entity

import "time"

type Order struct {
	ID             int       `json:"id" db:"id"`
	ProductID      int       `json:"product_id" db:"product_id"`
	UserID         string    `json:"user_id" db:"user_id"`
	Quantity       int       `json:"quantity" db:"quantity"`
	Status         string    `json:"status" db:"status"`
	IdempotencyKey string    `json:"idempotency_key,omitempty" db:"idempotency_key"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

const (
	OrderStatusPending   = "pending"
	OrderStatusCompleted = "completed"
	OrderStatusCancelled = "cancelled"
)
