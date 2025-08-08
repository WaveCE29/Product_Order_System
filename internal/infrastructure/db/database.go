package database

import (
	"database/sql"
	"fmt"

	"github.com/WaveCE29/product_order_system/pkg/logger"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB     *sql.DB
	logger logger.Logger
}

func NewDatabase(dbPath string, logger logger.Logger) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Error("Failed to open database", "error", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{
		DB:     db,
		logger: logger,
	}

	if err := database.migrate(); err != nil {
		logger.Error("Failed to migrate database", "error", err)
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	logger.Info("Database connection established and migrations applied successfully")
	return database, nil
}

func (d *Database) migrate() error {
	d.logger.Info("Running database migrations")

	queries := []string{
		`CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			name TEXT NOT NULL,
			stock INTEGER NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			product_id INTEGER NOT NULL,            
			user_id TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			status TEXT NOT NULL,
			idempotency_key TEXT UNIQUE,
			created_at DATETIME NOT NULL,
			FOREIGN KEY (product_id) REFERENCES products (id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_orders_idempotency_key ON orders(idempotency_key)`,
		`CREATE INDEX IF NOT EXISTS idx_orders_product_id ON orders(product_id)`,
		`CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id)`,
	}

	for _, query := range queries {
		if _, err := d.DB.Exec(query); err != nil {
			return fmt.Errorf("failed to execute migration query: %w", err)
		}
	}

	d.logger.Info("Database migrations completed successfully")
	return nil
}

func (d *Database) Close() error {
	d.logger.Info("Closing database connection")
	return d.DB.Close()
}
