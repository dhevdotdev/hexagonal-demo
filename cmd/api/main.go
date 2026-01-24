package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	httpAdapter "hexagonal-demo/internal/adapters/http"
	"hexagonal-demo/internal/adapters/postgres"
	"hexagonal-demo/internal/application"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Configuration
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "dhev")
	dbPassword := getEnv("DB_PASSWORD", "dhev")
	dbName := getEnv("DB_NAME", "bookstore")
	serverPort := getEnv("SERVER_PORT", "8080")

	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	// Connect to PostgreSQL (OUTBOUND ADAPTER)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL database")

	// DEPENDENCY INJECTION - Wiring everything together
	// This is where we assemble the hexagon!

	// 1. Create the PostgreSQL adapter (implements outbound port)
	bookRepo := postgres.NewBookRepository(db)

	// 2. Create the application service (implements inbound port)
	//    Inject the repository (dependency inversion!)
	bookService := application.NewBookService(bookRepo)

	// 3. Create the HTTP adapter (uses inbound port)
	bookHandler := httpAdapter.NewBookHandler(bookService)

	// Setup HTTP router
	router := mux.NewRouter()
	bookHandler.RegisterRoutes(router)

	// Add a health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Start server
	addr := fmt.Sprintf(":%s", serverPort)
	log.Printf("Starting HTTP server on %s", addr)
	log.Printf("API Endpoints:")
	log.Printf("  POST   /books         - Create a book")
	log.Printf("  GET    /books         - List all books")
	log.Printf("  GET    /books/{id}    - Get a book by ID")
	log.Printf("  PUT    /books/{id}    - Update a book")
	log.Printf("  DELETE /books/{id}    - Delete a book")
	log.Printf("  GET    /health        - Health check")

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
