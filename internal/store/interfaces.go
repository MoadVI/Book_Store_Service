package store

import (
	"Book-Store/internal/models"
	"time"
)

type BookStore interface {
	CreateBook(book models.Book) (models.Book, error)
	GetBook(id int) (models.Book, error)
	UpdateBook(id int, book models.Book) (models.Book, error)
	DeleteBook(id int) error
	SearchBooks(criteria models.SearchCriteria) ([]models.Book, error)
}

type AuthorStore interface {
	CreateAuthor(author models.Author) (models.Author, error)
	GetAuthor(id int) (models.Author, error)
	ListAuthors() ([]models.Author, error)
	UpdateAuthor(id int, author models.Author) (models.Author, error)
	DeleteAuthor(id int) error
}

type CustomerStore interface {
	CreateCustomer(customer models.Customer) (models.Customer, error)
	GetCustomer(id int) (models.Customer, error)
	UpdateCustomer(id int, customer models.Customer) (models.Customer, error)
	ListCustomers() ([]models.Customer, error)
}

type OrderStore interface {
	CreateOrder(order models.Order) (models.Order, error)
	GetOrder(id int) (models.Order, error)
	ListOrders() ([]models.Order, error)

	GetOrdersInTimeRange(start, end time.Time) ([]models.Order, error)
}
