package domain

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

// Book represents a book in our library system
// This is a domain entity - it contains business logic and validation
// It has NO dependencies on infrastructure (no database, no HTTP, etc.)
type Book struct {
	ID        string
	Title     string
	Author    string
	ISBN      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Domain errors
var (
	ErrInvalidTitle  = errors.New("title cannot be empty")
	ErrInvalidAuthor = errors.New("author cannot be empty")
	ErrInvalidISBN   = errors.New("invalid ISBN format")
)

// ISBN regex pattern (simplified - matches ISBN-10 and ISBN-13 with optional hyphens)
// Pattern matches: 978-0132350884, 0132350884, 978-0-13-235088-4, etc.
var isbnPattern = regexp.MustCompile(`^(?:ISBN(?:-1[03])?:?\s*)?(?:97[89][- ]?)?[0-9][- ]?[0-9][- ]?[0-9][- ]?[0-9][- ]?[0-9][- ]?[0-9][- ]?[0-9][- ]?[0-9][- ]?[0-9][- ]?[0-9Xx](?:[- ]?[0-9Xx]){0,3}$`)

// NewBook creates a new Book with validation
// This is a factory function that enforces business rules
func NewBook(title, author, isbn string) (*Book, error) {
	if err := ValidateTitle(title); err != nil {
		return nil, err
	}

	if err := ValidateAuthor(author); err != nil {
		return nil, err
	}

	if err := ValidateISBN(isbn); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Book{
		ID:        uuid.New().String(),
		Title:     title,
		Author:    author,
		ISBN:      isbn,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Update updates the book's mutable fields with validation
func (b *Book) Update(title, author, isbn string) error {
	if err := ValidateTitle(title); err != nil {
		return err
	}

	if err := ValidateAuthor(author); err != nil {
		return err
	}

	if err := ValidateISBN(isbn); err != nil {
		return err
	}

	b.Title = title
	b.Author = author
	b.ISBN = isbn
	b.UpdatedAt = time.Now()

	return nil
}

// ValidateTitle checks if the title is valid
func ValidateTitle(title string) error {
	if title == "" {
		return ErrInvalidTitle
	}
	return nil
}

// ValidateAuthor checks if the author is valid
func ValidateAuthor(author string) error {
	if author == "" {
		return ErrInvalidAuthor
	}
	return nil
}

// ValidateISBN checks if the ISBN is in a valid format
func ValidateISBN(isbn string) error {
	if isbn == "" || !isbnPattern.MatchString(isbn) {
		return ErrInvalidISBN
	}
	return nil
}
