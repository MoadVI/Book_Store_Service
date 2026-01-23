package handlers

import (
	"Book-Store/internal/response"
	"Book-Store/internal/store"
	"fmt"
	"net/http"
)

type MetricsHandler struct {
	BookStore     store.BookStore
	AuthorStore   store.AuthorStore
	CustomerStore store.CustomerStore
	OrderStore    store.OrderStore
}

func (m *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	q := r.URL.Query()

	switch {
	case q.Get("total_customers") != "":
		m.totalCustomers(w)
	case q.Get("total_books") != "":
		m.totalBooks(w)
	case q.Get("books_per_genre") != "":
		genre := r.URL.Query().Get("genre")
		m.getBooksPerGenre(w, genre)
	case q.Get("out_of_stock_books") != "":
		m.outOfStockBooks(w)
	case q.Get("total_authors") != "":
		m.totalAuthors(w)
	case q.Get("books_per_author") != "":
		m.getBooksPerAuthor(w)
	default:
		http.Error(w, "Unknown metric", http.StatusBadRequest)

	}
}

func (m *MetricsHandler) totalCustomers(w http.ResponseWriter) {
	fmt.Fprintf(w, "Total Customers: %d\n", m.CustomerStore.CustomersCount())
}

func (m *MetricsHandler) totalBooks(w http.ResponseWriter) {
	fmt.Fprintf(w, "Total Books: %d\n", m.BookStore.BooksCount())
}

func (m *MetricsHandler) totalAuthors(w http.ResponseWriter) {
	fmt.Fprintf(w, "Total Authors: %d\n", m.AuthorStore.AuthorsCount())
}

func (m *MetricsHandler) outOfStockBooks(w http.ResponseWriter) {
	booksOutOfStock := m.BookStore.OutOfStock
	response.RespondWithJSON(w, http.StatusOK, booksOutOfStock())
}

func (m *MetricsHandler) getBooksPerGenre(w http.ResponseWriter, genre string) {
	books := m.BookStore.GetBooksPerGenre(genre)
	response.RespondWithJSON(w, http.StatusOK, books)
}

func (m *MetricsHandler) getBooksPerAuthor(w http.ResponseWriter) {
	fmt.Fprintln(w, "Books per author:")
	for authorID, count := range m.AuthorStore.BooksPerAuthor() {
		fmt.Fprintf(w, "Author %d: %d\n", authorID, count)
	}
}
