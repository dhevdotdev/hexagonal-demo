package cli

import (
	"fmt"
	"hexagonal-demo/internal/ports"
)

// BookCLI is an INBOUND ADAPTER for command-line interface
// It translates CLI commands into application service calls
type BookCLI struct {
	service ports.BookService
}

// NewBookCLI creates a new CLI adapter
func NewBookCLI(service ports.BookService) *BookCLI {
	return &BookCLI{
		service: service,
	}
}

// CreateBook creates a book via CLI
func (cli *BookCLI) CreateBook(title, author, isbn string) error {
	book, err := cli.service.CreateBook(title, author, isbn)
	if err != nil {
		return fmt.Errorf("failed to create book: %w", err)
	}

	fmt.Printf("Book created successfully!\n")
	fmt.Printf("ID: %s\n", book.ID)
	fmt.Printf("Title: %s\n", book.Title)
	fmt.Printf("Author: %s\n", book.Author)
	fmt.Printf("ISBN: %s\n", book.ISBN)
	fmt.Printf("Created At: %s\n", book.CreatedAt.Format("2006-01-02 15:04:05"))

	return nil
}

// GetBook retrieves a book by ID via CLI
func (cli *BookCLI) GetBook(id string) error {
	book, err := cli.service.GetBook(id)
	if err != nil {
		return fmt.Errorf("failed to get book: %w", err)
	}

	fmt.Printf("Book Details:\n")
	fmt.Printf("ID: %s\n", book.ID)
	fmt.Printf("Title: %s\n", book.Title)
	fmt.Printf("Author: %s\n", book.Author)
	fmt.Printf("ISBN: %s\n", book.ISBN)
	fmt.Printf("Created At: %s\n", book.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated At: %s\n", book.UpdatedAt.Format("2006-01-02 15:04:05"))

	return nil
}

// ListBooks lists all books via CLI
func (cli *BookCLI) ListBooks() error {
	books, err := cli.service.ListBooks()
	if err != nil {
		return fmt.Errorf("failed to list books: %w", err)
	}

	if len(books) == 0 {
		fmt.Println("No books found.")
		return nil
	}

	fmt.Printf("Found %d book(s):\n\n", len(books))
	for i, book := range books {
		fmt.Printf("%d. %s by %s (ISBN: %s)\n", i+1, book.Title, book.Author, book.ISBN)
		fmt.Printf("   ID: %s\n", book.ID)
		fmt.Printf("   Created: %s\n\n", book.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	return nil
}

// UpdateBook updates a book via CLI
func (cli *BookCLI) UpdateBook(id, title, author, isbn string) error {
	book, err := cli.service.UpdateBook(id, title, author, isbn)
	if err != nil {
		return fmt.Errorf("failed to update book: %w", err)
	}

	fmt.Printf("Book updated successfully!\n")
	fmt.Printf("ID: %s\n", book.ID)
	fmt.Printf("Title: %s\n", book.Title)
	fmt.Printf("Author: %s\n", book.Author)
	fmt.Printf("ISBN: %s\n", book.ISBN)
	fmt.Printf("Updated At: %s\n", book.UpdatedAt.Format("2006-01-02 15:04:05"))

	return nil
}

// DeleteBook deletes a book via CLI
func (cli *BookCLI) DeleteBook(id string) error {
	if err := cli.service.DeleteBook(id); err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	fmt.Printf("Book deleted successfully! (ID: %s)\n", id)
	return nil
}
