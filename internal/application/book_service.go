package application

import (
	"errors"
	"hexagonal-demo/internal/domain"
	"hexagonal-demo/internal/ports"
)

// Common application errors
var (
	ErrBookNotFound = errors.New("book not found")
)

// bookService implements the BookService port (inbound port)
// This is where the USE CASES live
// It orchestrates domain objects and calls the repository (outbound port)
type bookService struct {
	// We depend on the PORT (interface), not the concrete implementation
	// We don't know if this is PostgreSQL, MongoDB, or in-memory storage
	repo ports.BookRepository
}

// NewBookService creates a new BookService
// This is dependency injection - we inject the repository implementation
func NewBookService(repo ports.BookRepository) ports.BookService {
	return &bookService{
		repo: repo,
	}
}

// CreateBook implements the CreateBook use case
func (s *bookService) CreateBook(title, author, isbn string) (*domain.Book, error) {
	// Step 1: Use the domain factory to create a validated book
	// The domain layer handles all business validation
	book, err := domain.NewBook(title, author, isbn)
	if err != nil {
		return nil, err
	}

	// Step 2: Persist the book using the repository port
	// We don't care HOW it's saved - that's the adapter's job
	if err := s.repo.Save(book); err != nil {
		return nil, err
	}

	return book, nil
}

// GetBook implements the GetBook use case
func (s *bookService) GetBook(id string) (*domain.Book, error) {
	book, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if book == nil {
		return nil, ErrBookNotFound
	}

	return book, nil
}

// ListBooks implements the ListBooks use case
func (s *bookService) ListBooks() ([]*domain.Book, error) {
	return s.repo.FindAll()
}

// UpdateBook implements the UpdateBook use case
func (s *bookService) UpdateBook(id, title, author, isbn string) (*domain.Book, error) {
	// Step 1: Retrieve the existing book
	book, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if book == nil {
		return nil, ErrBookNotFound
	}

	// Step 2: Update the book using domain logic (with validation)
	if err := book.Update(title, author, isbn); err != nil {
		return nil, err
	}

	// Step 3: Persist the updated book
	if err := s.repo.Update(book); err != nil {
		return nil, err
	}

	return book, nil
}

// DeleteBook implements the DeleteBook use case
func (s *bookService) DeleteBook(id string) error {
	// Check if book exists first
	book, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if book == nil {
		return ErrBookNotFound
	}

	return s.repo.Delete(id)
}
