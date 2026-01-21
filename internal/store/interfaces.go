package store

import (
	"Book-Store/internal/models"
	"context"
	"time"
)

type BookStore interface {
	CreateBook(ctx context.Context, book models.Book) (models.Book, error)
	GetBook(ctx context.Context, id int) (models.Book, error)
	UpdateBook(ctx context.Context, id int, book models.Book) (models.Book, error)
	DeleteBook(ctx context.Context, id int) error
	SearchBooks(ctx context.Context, criteria models.SearchCriteria) ([]models.Book, error)
	BookExists(id int) bool
}

type AuthorStore interface {
	CreateAuthor(ctx context.Context, author models.Author) (models.Author, error)
	GetAuthor(ctx context.Context, id int) (models.Author, error)
	ListAuthors(ctx context.Context) ([]models.Author, error)
	UpdateAuthor(ctx context.Context, id int, author models.Author) (models.Author, error)
	DeleteAuthor(ctx context.Context, id int) error
	AuthorExists(id int) bool
}

type CustomerStore interface {
	CreateCustomer(ctx context.Context, customer models.Customer) (models.Customer, error)
	GetCustomer(ctx context.Context, id int) (models.Customer, error)
	UpdateCustomer(ctx context.Context, id int, customer models.Customer) (models.Customer, error)
	ListCustomers(ctx context.Context) ([]models.Customer, error)
	DeleteCustomer(ctx context.Context, id int) error
	CustomerExists(id int) bool
}

type OrderStore interface {
	CreateOrder(ctx context.Context, order models.Order) (models.Order, error)
	GetOrder(ctx context.Context, id int) (models.Order, error)
	ListOrders(ctx context.Context) ([]models.Order, error)
	SearchOrderByStatus(ctx context.Context, status string) ([]models.Order, error)
	CompleteOrder(ctx context.Context, id int) (bool, error)
	CancelOrder(ctx context.Context, id int) (bool, error)
	GetOrdersInTimeRange(ctx context.Context, start, end time.Time) ([]models.Order, error)
}

