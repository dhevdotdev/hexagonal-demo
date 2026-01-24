package postgres

import (
	"database/sql"
	"hexagonal-demo/internal/domain"
	"hexagonal-demo/internal/ports"
	"time"

	_ "github.com/lib/pq"
)

// bookRepository implements the BookRepository port (outbound port)
// This is an ADAPTER that connects our application to PostgreSQL
// It translates between domain objects and database rows
type bookRepository struct {
	db *sql.DB
}

// NewBookRepository creates a new PostgreSQL book repository
func NewBookRepository(db *sql.DB) ports.BookRepository {
	return &bookRepository{
		db: db,
	}
}

// Save persists a book to PostgreSQL
func (r *bookRepository) Save(book *domain.Book) error {
	query := `
		INSERT INTO books (id, title, author, isbn, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(
		query,
		book.ID,
		book.Title,
		book.Author,
		book.ISBN,
		book.CreatedAt,
		book.UpdatedAt,
	)

	return err
}

// FindByID retrieves a book by ID from PostgreSQL
func (r *bookRepository) FindByID(id string) (*domain.Book, error) {
	query := `
		SELECT id, title, author, isbn, created_at, updated_at
		FROM books
		WHERE id = $1
	`

	var book domain.Book
	err := r.db.QueryRow(query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.ISBN,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &book, nil
}

// FindAll retrieves all books from PostgreSQL
func (r *bookRepository) FindAll() ([]*domain.Book, error) {
	query := `
		SELECT id, title, author, isbn, created_at, updated_at
		FROM books
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*domain.Book
	for rows.Next() {
		var book domain.Book
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	return books, rows.Err()
}

// Update updates an existing book in PostgreSQL
func (r *bookRepository) Update(book *domain.Book) error {
	query := `
		UPDATE books
		SET title = $1, author = $2, isbn = $3, updated_at = $4
		WHERE id = $5
	`

	// Update the timestamp
	book.UpdatedAt = time.Now()

	_, err := r.db.Exec(
		query,
		book.Title,
		book.Author,
		book.ISBN,
		book.UpdatedAt,
		book.ID,
	)

	return err
}

// Delete removes a book from PostgreSQL
func (r *bookRepository) Delete(id string) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
