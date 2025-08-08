# Product Order System

A simple Product Order System API built with Go, following Clean Architecture principles.

## Features

- **Product Management**: Create and retrieve products with stock management
- **Order Processing**: Create orders with automatic stock deduction
- **Idempotency**: Prevents duplicate orders using idempotency keys
- **Clean Architecture**: Separation of concerns with clear boundaries
- **SQLite Database**: Lightweight database for data persistence
- **Structured Logging**: JSON-structured logging with Zap
- **Graceful Shutdown**: Proper server shutdown handling

## Tech Stack

- **Language**: Go 1.21+
- **Web Framework**: Fiber v2
- **Database**: SQLite3
- **Logging**: Zap
- **Architecture**: Clean Architecture

## Project Structure

```
Product Order System/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── adapter/
│   │   └── http/
│   │       ├── handler/
│   │       │   └── handler.go      # HTTP handlers
│   │       └── router/
│   │           └── router.go       # Route definitions
│   ├── application/
│   │   ├── port/
│   │   │   └── input/
│   │   │       └── interfaces.go   # Use case interfaces
│   │   └── usecase/
│   │       ├── product_usecase.go  # Product business logic
│   │       └── order_usecase.go    # Order business logic
│   ├── domain/
│   │   ├── entity/
│   │   │   └── product.go          # Domain entities
│   │   └── repository/
│   │       └── interfaces.go       # Repository interfaces
│   └── infrastructure/
│       ├── config/
│       │   └── config.go           # Configuration management
│       ├── database/
│       │   └── database.go         # Database connection & migrations
│       └── persistence/
│           ├── product_repository.go
│           └── order_repository.go
├── pkg/
│   └── logger/
│       └── logger.go               # Logging utilities
├── .env                            # Environment variables
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## API Endpoints

### Products

#### Create Product

```http
POST /products
Content-Type: application/json

{
  "name": "Product Name",
  "stock": 100
}
```

#### Get All Products

```http
GET /products
```

#### Get Product by ID

```http
GET /products/:id
```

### Orders

#### Create Order

```http
POST /orders
Content-Type: application/json

{
  "product_id": 1,
  "user_id": "user123",
  "quantity": 2,
  "idempotency_key": "unique-key-123"
}
```

## Installation & Usage

### Prerequisites

- Go 1.21 or higher
- SQLite3

### Setup

1. Clone the repository:

```bash
git clone <repository-url>
cd Product Order System
```

2. Install dependencies:

```bash
go mod download
```

3. Configure environment variables:

```bash
cp .env.example .env
# Edit .env file as needed
```

4. Run the application:

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

### Build

To build the application:

```bash
go build -o bin/server cmd/server/main.go
```

## Configuration

The application can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `HOST` | Server host | `0.0.0.0` |
| `DATABASE_PATH` | SQLite database file path | `./ecommerce.db` |

## Database Schema

### Products Table

```sql
CREATE TABLE products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    stock INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
```

### Orders Table

```sql
CREATE TABLE orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    product_id INTEGER NOT NULL,
    user_id TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    status TEXT NOT NULL,
    idempotency_key TEXT UNIQUE,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products (id)
);
```

## Architecture

This project follows Clean Architecture principles:

- **Domain Layer**: Contains business entities and repository interfaces
- **Application Layer**: Contains use cases and business logic
- **Infrastructure Layer**: Contains external concerns (database, HTTP)
- **Adapter Layer**: Contains HTTP handlers and routing

## Features

### Stock Management

- Automatic stock deduction when orders are created
- Stock validation to prevent overselling
- Real-time stock tracking

### Idempotency

- Prevents duplicate order creation using idempotency keys
- Returns existing order if duplicate key is detected

### Error Handling

- Comprehensive error handling with appropriate HTTP status codes
- Structured error responses
- Request validation

### Logging

- Structured JSON logging
- Request/response logging
- Error tracking

## Testing

Run tests with:

```bash
go test ./...
```

## Health Check

The API includes a health check endpoint:

```http
GET /health
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License.
