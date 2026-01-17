package router

import (
	"Book-Store/internal/http/handlers"
	"net/http"
)

func Router(
	bookHandler *handlers.BookHandler,
) {
	http.Handle("/books/", bookHandler)
}
