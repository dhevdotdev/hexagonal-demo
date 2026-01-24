package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"hexagonal-demo/internal/adapters/cli"
	"hexagonal-demo/internal/adapters/postgres"
	"hexagonal-demo/internal/application"

	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Configuration
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "bookstore")

	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// DEPENDENCY INJECTION - Same wiring as the HTTP API

	// 1. Create the PostgreSQL adapter (implements outbound port)
	bookRepo := postgres.NewBookRepository(db)

	// 2. Create the application service (implements inbound port)
	bookService := application.NewBookService(bookRepo)

	// 3. Create the CLI adapter (uses inbound port)
	bookCLI := cli.NewBookCLI(bookService)

	// Parse command and execute
	command := os.Args[1]

	switch command {
	case "create":
		if len(os.Args) != 5 {
			fmt.Println("Usage: cli create <title> <author> <isbn>")
			os.Exit(1)
		}
		err = bookCLI.CreateBook(os.Args[2], os.Args[3], os.Args[4])

	case "get":
		if len(os.Args) != 3 {
			fmt.Println("Usage: cli get <id>")
			os.Exit(1)
		}
		err = bookCLI.GetBook(os.Args[2])

	case "list":
		err = bookCLI.ListBooks()

	case "update":
		if len(os.Args) != 6 {
			fmt.Println("Usage: cli update <id> <title> <author> <isbn>")
			os.Exit(1)
		}
		err = bookCLI.UpdateBook(os.Args[2], os.Args[3], os.Args[4], os.Args[5])

	case "delete":
		if len(os.Args) != 3 {
			fmt.Println("Usage: cli delete <id>")
			os.Exit(1)
		}
		err = bookCLI.DeleteBook(os.Args[2])

	default:
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func printUsage() {
	fmt.Println("Book Library CLI - Hexagonal Architecture Demo")
	fmt.Println("\nUsage:")
	fmt.Println("  cli create <title> <author> <isbn>")
	fmt.Println("  cli get <id>")
	fmt.Println("  cli list")
	fmt.Println("  cli update <id> <title> <author> <isbn>")
	fmt.Println("  cli delete <id>")
	fmt.Println("\nExamples:")
	fmt.Println("  cli create \"Clean Code\" \"Robert C. Martin\" \"978-0132350884\"")
	fmt.Println("  cli list")
	fmt.Println("  cli get abc-123-def")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
