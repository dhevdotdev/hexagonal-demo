package http

import (
	"encoding/json"
	"hexagonal-demo/internal/domain"
	"hexagonal-demo/internal/ports"
	"net/http"

	"github.com/gorilla/mux"
)

// BookHandler is an INBOUND ADAPTER for HTTP
// It translates HTTP requests into application service calls
// It depends on the BookService PORT, not the concrete implementation
type BookHandler struct {
	service ports.BookService
}

// NewBookHandler creates a new HTTP handler
func NewBookHandler(service ports.BookService) *BookHandler {
	return &BookHandler{
		service: service,
	}
}

// Request/Response DTOs (Data Transfer Objects)
// These are HTTP-specific and don't leak into the domain

type CreateBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	ISBN   string `json:"isbn"`
}

type UpdateBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	ISBN   string `json:"isbn"`
}

type BookResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	ISBN      string `json:"isbn"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// RegisterRoutes registers all book-related HTTP routes
func (h *BookHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/books", h.CreateBook).Methods("POST")
	router.HandleFunc("/books", h.ListBooks).Methods("GET")
	router.HandleFunc("/books/{id}", h.GetBook).Methods("GET")
	router.HandleFunc("/books/{id}", h.UpdateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", h.DeleteBook).Methods("DELETE")
}

// CreateBook handles POST /books
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Call the application service (through the port)
	book, err := h.service.CreateBook(req.Title, req.Author, req.ISBN)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, toBookResponse(book))
}

// GetBook handles GET /books/{id}
func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	book, err := h.service.GetBook(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, toBookResponse(book))
}

// ListBooks handles GET /books
func (h *BookHandler) ListBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.service.ListBooks()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := make([]BookResponse, len(books))
	for i, book := range books {
		response[i] = toBookResponse(book)
	}

	respondWithJSON(w, http.StatusOK, response)
}

// UpdateBook handles PUT /books/{id}
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	book, err := h.service.UpdateBook(id, req.Title, req.Author, req.ISBN)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, toBookResponse(book))
}

// DeleteBook handles DELETE /books/{id}
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.DeleteBook(id); err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper functions

func toBookResponse(book *domain.Book) BookResponse {
	return BookResponse{
		ID:        book.ID,
		Title:     book.Title,
		Author:    book.Author,
		ISBN:      book.ISBN,
		CreatedAt: book.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: book.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{Error: message})
}
