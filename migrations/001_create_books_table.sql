-- Create books table
CREATE TABLE IF NOT EXISTS books (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    isbn VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on ISBN for faster lookups
CREATE INDEX IF NOT EXISTS idx_books_isbn ON books(isbn);

-- Create index on author for faster searches
CREATE INDEX IF NOT EXISTS idx_books_author ON books(author);
