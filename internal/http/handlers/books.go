package handlers

import (
	"Book-Store/internal/models"
	"Book-Store/internal/response"
	"Book-Store/internal/store"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type BookHandler struct {
	BookStore   store.BookStore
	AuthorStore store.AuthorStore
}

func (h *BookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")

	var id int
	var hasID bool

	if len(pathParts) > 1 && pathParts[1] != "" {
		parsedID, err := strconv.Atoi(pathParts[1])
		if err == nil {
			id = parsedID
			hasID = true
		}
	}

	switch r.Method {
	case http.MethodPost:
		h.createBook(w, r)
	case http.MethodGet:
		if hasID {
			h.getBookById(w, id)
		} else {
			h.searchBooks(w, r)
		}
	case http.MethodPut:
		if !hasID {
			response.RespondWithError(w, http.StatusBadRequest, "Missing book ID")
			return
		}
		h.updateBook(w, r, id)
	case http.MethodDelete:
		if !hasID {
			response.RespondWithError(w, http.StatusBadRequest, "Missing Book ID")
			return
		}
		h.deleteBook(w, id)

	default:
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *BookHandler) createBook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		fmt.Println("JSON Decode Error in createBook: ", err)
		response.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if !h.AuthorStore.AuthorExists(book.Author.ID) {
		log.Printf("Cannot create book: author %d not found", book.Author.ID)
		response.RespondWithError(w, http.StatusNotFound, "Author not found")
		return
	}
	createdBook, _ := h.BookStore.CreateBook(book)
	response.RespondWithJSON(w, http.StatusCreated, createdBook)

}

func (h *BookHandler) searchBooks(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	author := r.URL.Query().Get("author")
	genre := r.URL.Query().Get("genre")
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")

	var minPricePtr, maxPricePtr *float64
	if s := r.URL.Query().Get("min_price"); s != "" {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			minPricePtr = &f
		}
	}
	if s := r.URL.Query().Get("max_price"); s != "" {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			maxPricePtr = &f
		}
	}

	criteria := models.SearchCriteria{
		Title:     title,
		Author:    author,
		Genre:     genre,
		MinPrice:  minPricePtr,
		MaxPrice:  maxPricePtr,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	books, err := h.BookStore.SearchBooks(criteria)
	if err != nil {
		fmt.Println("Error finding Books:", err)
		response.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, books)
}

func (h *BookHandler) getBookById(w http.ResponseWriter, id int) {
	book, lookUpError := h.BookStore.GetBook(id)
	if lookUpError != nil {
		fmt.Printf("Error finding the book using id : %d\nError: %v\n", id, lookUpError)
		response.RespondWithError(w, http.StatusNotFound, "Book does not exist")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, book)

}

func (h *BookHandler) updateBook(w http.ResponseWriter, r *http.Request, id int) {
	defer r.Body.Close()

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		fmt.Println("JSON Decode Error in updateBook: ", err)
		response.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if !h.AuthorStore.AuthorExists(book.Author.ID) {
		log.Printf("Cannot create book: author %d not found", book.Author.ID)
		response.RespondWithError(w, http.StatusNotFound, "Author not found")
		return
	}

	updated_book, update_err := h.BookStore.UpdateBook(id, book)
	if update_err != nil {
		response.RespondWithError(w, http.StatusNotFound, "Book not found")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, updated_book)

}

func (h *BookHandler) deleteBook(w http.ResponseWriter, id int) {
	delete_err := h.BookStore.DeleteBook(id)
	if delete_err != nil {
		response.RespondWithError(w, http.StatusNotFound, "Book not found")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, "Book deleted successfully")
}
