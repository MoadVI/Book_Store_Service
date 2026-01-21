package router

import (
	"Book-Store/internal/http/handlers"
	"net/http"
)

func Router(
	bookHandler *handlers.BookHandler,
	authorHandler *handlers.AuthorHandler,
	customerHandler *handlers.CustomerHandler,
	orderHandler *handlers.OrderHandler,
) {
	http.Handle("/books/", bookHandler)
	http.Handle("/authors/", authorHandler)
	http.Handle("/customers/", customerHandler)
	http.Handle("/orders/", orderHandler)
}
