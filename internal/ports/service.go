package ports

import "hexagonal-demo/internal/domain"

// BookService is an INBOUND PORT (Driving Port)
// It defines what our application can DO - the use cases
// This interface is called BY adapters (HTTP, CLI, etc.)
// The application layer implements this interface
type BookService interface {
	// CreateBook creates a new book in the system
	CreateBook(title, author, isbn string) (*domain.Book, error)

	// GetBook retrieves a book by its ID
	GetBook(id string) (*domain.Book, error)

	// ListBooks retrieves all books in the system
	ListBooks() ([]*domain.Book, error)

	// UpdateBook updates an existing book
	UpdateBook(id, title, author, isbn string) (*domain.Book, error)

	// DeleteBook removes a book from the system
	DeleteBook(id string) error
}
