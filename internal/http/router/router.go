package router

import (
	"Book-Store/internal/http/handlers"
	"net/http"
)

func Router(
	bookHandler *handlers.BookHandler,
	authorHandler *handlers.AuthorHandler,
) {
	http.Handle("/books/", bookHandler)
	http.Handle("/authors/", authorHandler)
}
