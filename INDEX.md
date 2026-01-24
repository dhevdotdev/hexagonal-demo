# Hexagonal Architecture Demo - Book Library

A complete, production-ready example of **Hexagonal Architecture** (Ports & Adapters pattern) implemented in Go with PostgreSQL.

## 📚 What's This?

This project demonstrates how to build maintainable, testable, and flexible applications using the Hexagonal Architecture pattern. It's a fully functional book library system with:

- ✅ **HTTP REST API** - Create, read, update, delete books
- ✅ **CLI Tool** - Command-line interface for the same operations
- ✅ **PostgreSQL Database** - Real persistence layer
- ✅ **Clean Architecture** - Business logic independent of infrastructure
- ✅ **Comprehensive Documentation** - Learn by example

## 🎯 Quick Start

### Prerequisites
- PostgreSQL 18 (installed with user `dhev` / password `dhev`)
- Go 1.25.6

### Run the Project

```bash
# 1. Clone or navigate to the project
cd C:\dev\hexagonal-demo

# 2. Set up database and build (Windows)
setup.bat

# 3. Try the CLI
bin\cli.exe create "Clean Code" "Robert C. Martin" "978-0132350884"
bin\cli.exe list

# 4. Or start the HTTP API
bin\api.exe
# Then visit: http://localhost:8080/books
```

## 📖 Documentation

Read in this order to understand the architecture:

1. **[README.md](README.md)** - Architecture overview with Mermaid diagrams
2. **[ARCHITECTURE.md](ARCHITECTURE.md)** - Deep dive into Ports & Adapters
3. **[QUICKSTART.md](QUICKSTART.md)** - How to run and use the apps
4. **[PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)** - Complete project summary

## 🏗️ Architecture Overview

```
┌─────────────────┐        ┌─────────────────┐
│   HTTP Client   │        │    CLI User     │
└────────┬────────┘        └────────┬────────┘
         │                          │
         v                          v
┌─────────────────┐        ┌─────────────────┐
│  HTTP Adapter   │        │   CLI Adapter   │
└────────┬────────┘        └────────┬────────┘
         │                          │
         └──────────┬───────────────┘
                    v
           ┌─────────────────┐
           │  BookService    │ ← Inbound Port (Interface)
           │   (Interface)   │
           └────────┬────────┘
                    v
           ┌─────────────────┐
           │ Application     │
           │ Service Layer   │ → Implements use cases
           └────────┬────────┘
                    v
           ┌─────────────────┐
           │ Domain Layer    │ → Pure business logic
           │  (Book Entity)  │
           └────────┬────────┘
                    v
           ┌─────────────────┐
           │ BookRepository  │ ← Outbound Port (Interface)
           │   (Interface)   │
           └────────┬────────┘
                    v
           ┌─────────────────┐
           │ PostgreSQL      │
           │    Adapter      │ → Database implementation
           └────────┬────────┘
                    v
                ┌───────┐
                │  DB   │
                └───────┘
```

## 🗂️ Project Structure

```
hexagonal-demo/
├── internal/
│   ├── domain/              # 🟢 CORE: Pure business logic
│   ├── ports/               # 🟡 BOUNDARY: Interface contracts
│   ├── application/         # 🔵 USE CASES: Orchestration
│   └── adapters/            # 🔴 INFRASTRUCTURE: Implementations
│       ├── http/           # Inbound: REST API
│       ├── cli/            # Inbound: CLI commands
│       └── postgres/       # Outbound: Database
├── cmd/
│   ├── api/                # HTTP server entry point
│   └── cli/                # CLI tool entry point
└── migrations/             # Database schema
```

## 🎓 Key Concepts Demonstrated

### 1. Dependency Inversion
- Core business logic doesn't know about HTTP or databases
- Adapters depend on ports (interfaces), not vice versa
- Dependencies point **inward** toward the domain

### 2. Ports & Adapters
- **Inbound Ports**: Define what your app can DO (use cases)
- **Outbound Ports**: Define what your app NEEDS (dependencies)
- **Adapters**: Concrete implementations of ports

### 3. Testability
- Domain layer: Pure Go, no mocks needed
- Application layer: Mock the repository interface
- Adapters: Integration tests with real infrastructure

### 4. Flexibility
Want to switch databases? Create new adapter, change one line in `main.go`
Want to add gRPC? Create new adapter, business logic untouched

## 🧪 Example Usage

### CLI
```bash
# Create books
bin\cli.exe create "Design Patterns" "Gang of Four" "978-0201633610"

# List all
bin\cli.exe list

# Get specific book
bin\cli.exe get {book-id}
```

### HTTP API
```bash
# Start server
bin\api.exe

# Create book (PowerShell)
curl -X POST http://localhost:8080/books `
  -H "Content-Type: application/json" `
  -d '{"title":"Clean Architecture","author":"Robert Martin","isbn":"978-0134494166"}'

# List books
curl http://localhost:8080/books
```

## 🔍 What Makes This "Hexagonal"?

The **hexagon** represents your application core. Each "side" is a port where adapters can plug in:

- **Left side** (driving/inbound): HTTP, CLI, gRPC, GraphQL
- **Right side** (driven/outbound): Database, APIs, Message queues, Cache

The shape doesn't matter - it's the **principle**: 
> "Business logic is independent of the delivery mechanism and infrastructure"

## 📝 Technologies Used

- **Go 1.25.6** - Language
- **PostgreSQL 18** - Database
- **Gorilla Mux** - HTTP routing
- **lib/pq** - PostgreSQL driver
- **google/uuid** - UUID generation

## 🎯 Learning Path

1. **Start simple**: Read `internal/domain/book.go` - pure business logic
2. **Understand contracts**: Check `internal/ports/` - interfaces
3. **See orchestration**: Review `internal/application/book_service.go`
4. **Follow the wiring**: Look at `cmd/api/main.go` - dependency injection
5. **Trace a request**: HTTP → Handler → Service → Domain → Repository → DB

## 💡 Why Use This Pattern?

**Traditional Layered Architecture:**
```
[UI] → [Business Logic with DB calls] → [Database]
```
Problem: Business logic is tightly coupled to infrastructure

**Hexagonal Architecture:**
```
[Adapters] → [Ports] → [Core Business Logic] ← [Ports] ← [Adapters]
```
Solution: Business logic is independent and testable

## 🚀 Extend This Project

Try these exercises to deepen your understanding:

1. **Add a field**: Add `PublishedYear` to books
2. **New adapter**: Create an in-memory repository for testing
3. **New use case**: Implement "Search books by author"
4. **New interface**: Add a gRPC adapter alongside HTTP
5. **Add caching**: Create a Redis adapter that wraps the repository

## 📊 Code Metrics

- **19 files** committed
- **2,342 lines** of code and documentation
- **0 dependencies** from domain on infrastructure
- **100%** business logic reusability across adapters

## 🎓 Resources

- [Original paper by Alistair Cockburn](https://alistair.cockburn.us/hexagonal-architecture/)
- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## 📄 License

This is a learning project - use it however you like!

## 🙋 Questions?

Read the detailed documentation files:
- Architecture deep dive: `ARCHITECTURE.md`
- Usage guide: `QUICKSTART.md`
- Complete summary: `PROJECT_SUMMARY.md`

---

**Built to demonstrate clean, maintainable architecture in Go** 🚀

Made with ❤️ for learning hexagonal architecture
