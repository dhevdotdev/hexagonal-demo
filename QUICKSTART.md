# Quick Start Guide

## Prerequisites

- PostgreSQL 18 installed and running
- Go 1.25.6 installed
- Database user `dhev` with password `dhev`

## Setup (First Time)

### Option 1: Automated Setup (Windows)

```bash
setup.bat
```

This will:
1. Create the `bookstore` database
2. Run migrations
3. Build the executables

### Option 2: Manual Setup

1. Create the database:
```bash
psql -U dhev -h localhost -c "CREATE DATABASE bookstore;"
```

2. Run migrations:
```bash
psql -U dhev -h localhost -d bookstore -f migrations\001_create_books_table.sql
```

3. Build the applications:
```bash
go build -o bin\api.exe cmd\api\main.go
go build -o bin\cli.exe cmd\cli\main.go
```

## Running the Applications

### HTTP API Server

Start the server:
```bash
bin\api.exe
```

The API will be available at `http://localhost:8080`

#### API Endpoints

**Create a book:**
```bash
curl -X POST http://localhost:8080/books ^
  -H "Content-Type: application/json" ^
  -d "{\"title\":\"Clean Code\",\"author\":\"Robert C. Martin\",\"isbn\":\"978-0132350884\"}"
```

**List all books:**
```bash
curl http://localhost:8080/books
```

**Get a specific book:**
```bash
curl http://localhost:8080/books/{id}
```

**Update a book:**
```bash
curl -X PUT http://localhost:8080/books/{id} ^
  -H "Content-Type: application/json" ^
  -d "{\"title\":\"Updated Title\",\"author\":\"Robert C. Martin\",\"isbn\":\"978-0132350884\"}"
```

**Delete a book:**
```bash
curl -X DELETE http://localhost:8080/books/{id}
```

### CLI Tool

**Create a book:**
```bash
bin\cli.exe create "Clean Code" "Robert C. Martin" "978-0132350884"
```

**List all books:**
```bash
bin\cli.exe list
```

**Get a specific book:**
```bash
bin\cli.exe get {id}
```

**Update a book:**
```bash
bin\cli.exe update {id} "New Title" "Robert C. Martin" "978-0132350884"
```

**Delete a book:**
```bash
bin\cli.exe delete {id}
```

## Example Workflow

```bash
# 1. Create some books
bin\cli.exe create "Clean Code" "Robert C. Martin" "978-0132350884"
bin\cli.exe create "The Pragmatic Programmer" "Andy Hunt" "978-0135957059"
bin\cli.exe create "Design Patterns" "Gang of Four" "978-0201633610"

# 2. List all books
bin\cli.exe list

# 3. Get a specific book (copy ID from list output)
bin\cli.exe get 67462a14-26b7-44ec-9dac-af367ea174eb

# 4. Start the HTTP API in another terminal
bin\api.exe

# 5. In another terminal, test the HTTP API
curl http://localhost:8080/books
```

## Configuration

Environment variables (optional):

- `DB_HOST` - Database host (default: `localhost`)
- `DB_PORT` - Database port (default: `5432`)
- `DB_USER` - Database user (default: `dhev`)
- `DB_PASSWORD` - Database password (default: `dhev`)
- `DB_NAME` - Database name (default: `bookstore`)
- `SERVER_PORT` - HTTP server port (default: `8080`)

Example:
```bash
set DB_HOST=192.168.1.100
set SERVER_PORT=3000
bin\api.exe
```

## Understanding the Code

Read the documentation in order:

1. **README.md** - Overview and architecture diagrams
2. **ARCHITECTURE.md** - Deep dive into hexagonal architecture
3. **Code tour**:
   - `internal/domain/book.go` - Start here: Pure business logic
   - `internal/ports/` - Interfaces (contracts)
   - `internal/application/book_service.go` - Use case implementation
   - `internal/adapters/` - Infrastructure implementations
   - `cmd/api/main.go` or `cmd/cli/main.go` - See how it all wires together

## Troubleshooting

**Cannot connect to database:**
- Check PostgreSQL is running: `psql -U dhev -h localhost`
- Verify database exists: `psql -U dhev -h localhost -c "\l"`
- Check pgpass.conf has correct credentials

**Build errors:**
- Run `go mod tidy` to clean up dependencies
- Ensure Go 1.21+ is installed

**Port already in use:**
- Change the port: `set SERVER_PORT=3000`
- Or stop other services using port 8080

## Next Steps

1. Try adding a new field to `Book` (e.g., `PublishYear`)
2. Implement a new adapter (e.g., in-memory repository for testing)
3. Add a new use case (e.g., search books by author)
4. Add another inbound adapter (e.g., WebSocket)

The beauty of hexagonal architecture: You can extend without breaking existing code!
