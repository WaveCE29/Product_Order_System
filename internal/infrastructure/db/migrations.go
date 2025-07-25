package database

import "database/sql"

func Migrate(db *sql.DB) error {
	createProductsTable := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		stock INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);`

	createOrdersTable := `
	CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		product_id INTEGER NOT NULL,
		user_id TEXT NOT NULL,
		quantity INTEGER NOT NULL,
		status TEXT NOT NULL,
		idempotency_key TEXT UNIQUE,
		created_at DATETIME NOT NULL,
		FOREIGN KEY (product_id) REFERENCES products (id)
	);`

	if _, err := db.Exec(createProductsTable); err != nil {
		return err
	}

	if _, err := db.Exec(createOrdersTable); err != nil {
		return err
	}

	return nil
}
