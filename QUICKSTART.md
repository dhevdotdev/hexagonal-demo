# Quick Start Guide

## Setup

### Prerequisites
- PostgreSQL installed and running
- Go 1.21+ installed

### Initialize

```bash
# 1. Create database
psql -U postgres -h localhost -c "CREATE DATABASE bookstore;"

# 2. Run migrations
psql -U postgres -h localhost -d bookstore -f migrations/001_create_books_table.sql

# 3. Build
go build -o bin/api cmd/api/main.go
go build -o bin/cli cmd/cli/main.go
```

### Configuration (Optional)

Set environment variables to override defaults:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=bookstore
export SERVER_PORT=8080
```

## Usage

### CLI Commands

```bash
# Create
./bin/cli create "Clean Code" "Robert Martin" "978-0132350884"

# List all
./bin/cli list

# Get specific
./bin/cli get {book-id}

# Update
./bin/cli update {book-id} "New Title" "Author" "ISBN"

# Delete
./bin/cli delete {book-id}
```

### HTTP API

Start server:
```bash
./bin/api
# Server running on :8080
```

**Endpoints:**

| Method | Path | Description |
|--------|------|-------------|
| POST | /books | Create book |
| GET | /books | List all books |
| GET | /books/{id} | Get book |
| PUT | /books/{id} | Update book |
| DELETE | /books/{id} | Delete book |
| GET | /health | Health check |

**Examples:**

```bash
# Create
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title":"Clean Code","author":"Robert Martin","isbn":"978-0132350884"}'

# List
curl http://localhost:8080/books

# Get
curl http://localhost:8080/books/{id}

# Update
curl -X PUT http://localhost:8080/books/{id} \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated","author":"Author","isbn":"ISBN"}'

# Delete
curl -X DELETE http://localhost:8080/books/{id}
```

## Exploring the Code

### Start with the core
1. Read `internal/domain/book.go` - Pure business logic
2. Read `internal/ports/*.go` - Interface contracts
3. Read `internal/application/book_service.go` - Use case implementation

### See how it connects
4. Read `cmd/api/main.go` - Dependency injection
5. Read `internal/adapters/http/handler.go` - HTTP adapter
6. Read `internal/adapters/postgres/book_repository.go` - Database adapter

### Trace a request
Start at `adapters/http/handler.go:CreateBook()` and follow the flow:
1. HTTP adapter receives request
2. Calls `BookService.CreateBook()` (port)
3. Application service validates via domain
4. Application service calls `BookRepository.Save()` (port)
5. PostgreSQL adapter executes SQL
6. Result flows back through the layers

## Exercises

### Beginner
1. Add a `PublishedYear` field to Book
2. Add validation for the new field
3. Update migrations, repository, and handlers

### Intermediate
1. Create an in-memory repository for testing
2. Implement "Find books by author" use case
3. Add input validation in HTTP adapter

### Advanced
1. Add a caching layer (decorator pattern)
2. Implement a gRPC adapter
3. Add event publishing when books are created
4. Implement soft deletes

## Troubleshooting

**Database connection failed:**
```bash
# Check PostgreSQL is running
psql -U postgres -h localhost

# Verify database exists
psql -U postgres -h localhost -l | grep bookstore
```

**Build errors:**
```bash
# Clean and rebuild
go mod tidy
go build ./...
```

**Port already in use:**
```bash
# Change the port
export SERVER_PORT=3000
./bin/api
```

## Next Steps

- Read [README.md](README.md) for architecture overview
- Read [ARCHITECTURE.md](ARCHITECTURE.md) for deep dive
- Experiment with adding new features
- Try swapping adapters (e.g., add MongoDB)
