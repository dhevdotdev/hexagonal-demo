# Hexagonal Architecture Deep Dive

## Understanding Ports and Adapters

This document provides an in-depth explanation of the Hexagonal Architecture pattern as implemented in this project.

## The Hexagon Metaphor

Imagine your business logic as a **hexagon** sitting in the center. The number of sides doesn't matter (it could be a circle), but the hexagon shape emphasizes that:

1. **All sides are equal** - No side is more important than another
2. **Multiple connection points** - You can have many adapters connected to your core
3. **Symmetry** - Inbound and outbound are treated similarly

```
        [HTTP Adapter]
              |
              | (calls)
              v
        [BookService Port] <--- INBOUND PORT
              |
              | (implemented by)
              v
    ╔═══════════════════════════════╗
    ║                               ║
    ║   APPLICATION SERVICE         ║
    ║   (Business Logic)            ║
    ║   - CreateBook()              ║
    ║   - GetBook()                 ║
    ║   - etc.                      ║
    ║                               ║
    ╚═══════════════════════════════╝
              |
              | (calls)
              v
        [Repository Port] <--- OUTBOUND PORT
              |
              | (implemented by)
              v
        [PostgreSQL Adapter]
              |
              v
          [Database]
```

## Detailed Component Breakdown

### 1. Domain Layer (The Core)

**Location**: `internal/domain/`

**Purpose**: Contains pure business logic with ZERO external dependencies.

**Key Characteristics**:
- No imports from infrastructure packages
- No database, HTTP, or framework code
- Just business entities and rules
- Can be tested without any external systems

**In Our Example** (`book.go`):
```go
type Book struct {
    ID        string
    Title     string
    Author    string
    ISBN      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

func NewBook(title, author, isbn string) (*Book, error) {
    // Validation logic here
    // Business rules enforced here
}
```

**Why This Matters**:
- If we switch from PostgreSQL to MongoDB, this code doesn't change
- If we move from HTTP to gRPC, this code doesn't change
- We can test business logic in pure Go without any mocks

### 2. Ports (The Boundaries)

**Location**: `internal/ports/`

**Purpose**: Define contracts (interfaces) between core and outside world.

#### 2.1 Inbound Ports (Driving Ports)

**What**: Interfaces that define what your application CAN DO (use cases)

**Who calls them**: Adapters (HTTP handlers, CLI, etc.)

**Who implements them**: Application service layer

**In Our Example** (`service.go`):
```go
type BookService interface {
    CreateBook(title, author, isbn string) (*domain.Book, error)
    GetBook(id string) (*domain.Book, error)
    // ... other use cases
}
```

**Flow**:
```
HTTP Request → HTTP Adapter → BookService Port → Application Service
```

#### 2.2 Outbound Ports (Driven Ports)

**What**: Interfaces that define what your application NEEDS (dependencies)

**Who calls them**: Application service layer

**Who implements them**: Infrastructure adapters (database, external APIs, etc.)

**In Our Example** (`repository.go`):
```go
type BookRepository interface {
    Save(book *domain.Book) error
    FindByID(id string) (*domain.Book, error)
    // ... other persistence operations
}
```

**Flow**:
```
Application Service → Repository Port → PostgreSQL Adapter → Database
```

### 3. Application Layer

**Location**: `internal/application/`

**Purpose**: Orchestrates business logic and coordinates between domain and ports.

**Key Characteristics**:
- Implements inbound ports
- Depends on outbound ports (interfaces, not implementations!)
- Contains use case logic
- Thin layer - mostly coordination

**In Our Example** (`book_service.go`):
```go
type bookService struct {
    repo ports.BookRepository  // Depends on INTERFACE, not concrete type
}

func (s *bookService) CreateBook(title, author, isbn string) (*domain.Book, error) {
    // 1. Use domain factory
    book, err := domain.NewBook(title, author, isbn)
    if err != nil {
        return nil, err
    }
    
    // 2. Save via repository port
    if err := s.repo.Save(book); err != nil {
        return nil, err
    }
    
    return book, nil
}
```

**Notice**: The application service has NO IDEA it's talking to PostgreSQL. It could be MySQL, MongoDB, in-memory storage - it doesn't care!

### 4. Adapters (The Outside World)

**Location**: `internal/adapters/`

**Purpose**: Connect external systems to your core application.

#### 4.1 Inbound Adapters (Driving Adapters)

**Examples**: HTTP REST API, CLI, gRPC, GraphQL

**In Our Project**:
- **HTTP** (`adapters/http/handler.go`): Translates HTTP requests to service calls
- **CLI** (`adapters/cli/commands.go`): Translates command-line args to service calls

**HTTP Adapter Example**:
```go
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
    // 1. Parse HTTP request
    var req CreateBookRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    // 2. Call service (through port)
    book, err := h.service.CreateBook(req.Title, req.Author, req.ISBN)
    
    // 3. Return HTTP response
    respondWithJSON(w, http.StatusCreated, toBookResponse(book))
}
```

**Key Point**: HTTP adapter converts HTTP-specific data to/from domain models.

#### 4.2 Outbound Adapters (Driven Adapters)

**Examples**: Database clients, External APIs, Message queues, File systems

**In Our Project**:
- **PostgreSQL** (`adapters/postgres/book_repository.go`): Implements repository port

**PostgreSQL Adapter Example**:
```go
func (r *bookRepository) Save(book *domain.Book) error {
    query := `INSERT INTO books (...) VALUES (...)`
    _, err := r.db.Exec(query, book.ID, book.Title, ...)
    return err
}
```

**Key Point**: PostgreSQL adapter converts domain models to/from database rows.

## The Dependency Rule

**Critical Rule**: Dependencies point INWARD, never outward.

```
Adapters (outer) → Ports (boundary) → Application (inner) → Domain (core)
```

### What This Means:

1. **Domain** depends on nothing
2. **Application** depends on Domain and Ports
3. **Adapters** depend on Ports and Application
4. **Ports** depend on Domain

### Why This Matters:

If you want to:
- Change HTTP framework (gin → echo) → Only modify HTTP adapter
- Switch databases (PostgreSQL → MongoDB) → Only modify database adapter
- Add new interface (gRPC) → Add new adapter, core unchanged
- Change business rules → Only modify domain/application

## Dependency Injection (DI)

**Where it happens**: In `main.go`

**What it does**: Wires everything together (assembles the hexagon)

```go
func main() {
    // 1. Create outbound adapter (infrastructure)
    db := connectToPostgreSQL()
    bookRepo := postgres.NewBookRepository(db)
    
    // 2. Create application service (inject repository)
    bookService := application.NewBookService(bookRepo)
    
    // 3. Create inbound adapter (inject service)
    bookHandler := http.NewBookHandler(bookService)
    
    // 4. Start server
    startHTTPServer(bookHandler)
}
```

**Notice the direction**:
```
main.go knows about everything
  ↓
HTTP adapter knows about ports.BookService interface
  ↓
Application service knows about ports.BookRepository interface
  ↓
Domain knows about nothing
```

## Testing Strategy

### Unit Testing Domain

```go
func TestNewBook(t *testing.T) {
    book, err := domain.NewBook("Clean Code", "Robert Martin", "978-0132350884")
    // No database, no HTTP, no mocks needed!
    assert.NoError(t, err)
    assert.NotEmpty(t, book.ID)
}
```

### Unit Testing Application Service

```go
func TestBookService_CreateBook(t *testing.T) {
    // Mock the repository port
    mockRepo := &MockBookRepository{}
    service := application.NewBookService(mockRepo)
    
    // Test the use case
    book, err := service.CreateBook("Title", "Author", "ISBN")
    
    assert.NoError(t, err)
    assert.Equal(t, 1, mockRepo.SaveCallCount)
}
```

### Integration Testing Adapters

```go
func TestPostgreSQLAdapter(t *testing.T) {
    // Use real database (testcontainers)
    db := setupTestDB()
    repo := postgres.NewBookRepository(db)
    
    book := &domain.Book{...}
    err := repo.Save(book)
    
    assert.NoError(t, err)
}
```

## Real-World Benefits

### Scenario 1: Adding gRPC Support

**Without Hexagonal Architecture**:
- Business logic mixed in HTTP handlers
- Need to rewrite everything for gRPC
- Risk of bugs in business logic

**With Hexagonal Architecture**:
```go
// 1. Create gRPC adapter (new file)
type BookGRPCService struct {
    service ports.BookService
}

func (g *BookGRPCService) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.BookResponse, error) {
    book, err := g.service.CreateBook(req.Title, req.Author, req.ISBN)
    // Convert to protobuf and return
}

// 2. In main.go, wire it up
grpcService := grpc.NewBookGRPCService(bookService)
startGRPCServer(grpcService)
```

**Result**: Added gRPC without touching core business logic!

### Scenario 2: Switch from PostgreSQL to MongoDB

**Change required**:
1. Create new adapter: `adapters/mongodb/book_repository.go`
2. Implement `ports.BookRepository` interface
3. In `main.go`, swap: `postgres.NewBookRepository(db)` → `mongodb.NewBookRepository(client)`

**Result**: Zero changes to business logic!

### Scenario 3: Add Caching

```go
// Create caching adapter that wraps repository
type CachedBookRepository struct {
    repo  ports.BookRepository
    cache Cache
}

func (c *CachedBookRepository) FindByID(id string) (*domain.Book, error) {
    // Check cache first
    if book := c.cache.Get(id); book != nil {
        return book, nil
    }
    
    // Call underlying repository
    book, err := c.repo.FindByID(id)
    if err == nil {
        c.cache.Set(id, book)
    }
    return book, err
}
```

**Result**: Added caching by wrapping existing adapter!

## Common Pitfalls

### ❌ Don't Do This:

**1. Domain importing database packages**:
```go
package domain

import "database/sql"  // WRONG!

type Book struct {
    db *sql.DB  // WRONG!
}
```

**2. HTTP logic in domain**:
```go
package domain

import "net/http"  // WRONG!

func (b *Book) HandleHTTP(w http.ResponseWriter) {  // WRONG!
}
```

**3. Direct adapter-to-adapter calls**:
```go
package http

import "myapp/adapters/postgres"  // WRONG!

func (h *Handler) CreateBook() {
    repo := postgres.NewBookRepository()  // WRONG!
    repo.Save(book)
}
```

### ✅ Do This Instead:

**1. Domain is pure**:
```go
package domain

type Book struct {
    ID     string
    Title  string
    // No infrastructure!
}
```

**2. Adapters stay at the edges**:
```go
package http

type BookHandler struct {
    service ports.BookService  // Depend on PORT, not concrete adapter
}
```

**3. Wire in main.go**:
```go
func main() {
    repo := postgres.NewBookRepository(db)
    service := application.NewBookService(repo)
    handler := http.NewBookHandler(service)
}
```

## Summary

**Hexagonal Architecture** = **Ports & Adapters**

- **Ports**: Interfaces defining contracts
  - Inbound: What your app can DO
  - Outbound: What your app NEEDS
  
- **Adapters**: Implementations connecting to external systems
  - Inbound: HTTP, CLI, gRPC (drive your app)
  - Outbound: Database, APIs (driven by your app)

- **Core Principle**: Business logic is independent of infrastructure

- **Key Rule**: Dependencies point INWARD

**Result**: Flexible, testable, maintainable applications!
