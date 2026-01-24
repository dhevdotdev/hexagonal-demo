package ports

import "hexagonal-demo/internal/domain"

// BookRepository is an OUTBOUND PORT (Driven Port)
// It defines what our application NEEDS from the persistence layer
// This interface is called BY the application to save/retrieve data
// The infrastructure layer (adapters) implements this interface
// The core domain has NO idea if this is PostgreSQL, MongoDB, or in-memory storage
type BookRepository interface {
	// Save persists a new book
	Save(book *domain.Book) error

	// FindByID retrieves a book by its ID
	FindByID(id string) (*domain.Book, error)

	// FindAll retrieves all books
	FindAll() ([]*domain.Book, error)

	// Update updates an existing book
	Update(book *domain.Book) error

	// Delete removes a book by its ID
	Delete(id string) error
}
