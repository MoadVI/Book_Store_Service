package handlerstests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"Book-Store/internal/http/handlers"
	"Book-Store/internal/models"
)

type MockStore struct {
	CreateFunc func(models.Book) (models.Book, error)
	GetFunc    func(int) (models.Book, error)
	UpdateFunc func(int, models.Book) (models.Book, error)
	DeleteFunc func(int) error
	SearchFunc func(models.SearchCriteria) ([]models.Book, error)
}

func (m *MockStore) CreateBook(b models.Book) (models.Book, error) {
	return m.CreateFunc(b)
}

func (m *MockStore) GetBook(id int) (models.Book, error) {
	return m.GetFunc(id)
}

func (m *MockStore) UpdateBook(id int, b models.Book) (models.Book, error) {
	return m.UpdateFunc(id, b)
}

func (m *MockStore) DeleteBook(id int) error {
	return m.DeleteFunc(id)
}

func (m *MockStore) SearchBooks(c models.SearchCriteria) ([]models.Book, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(c)
	}
	return []models.Book{}, nil
}

func TestCreateBook(t *testing.T) {
	mock := &MockStore{
		CreateFunc: func(b models.Book) (models.Book, error) {
			b.ID = 1
			return b, nil
		},
	}

	handler := &handlers.BookHandler{Store: mock}

	payload := models.Book{
		Title: "The Go Programming Language",
		Author: models.Author{
			FirstName: "Alan",
			LastName:  "Donovan",
		},
	}

	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/books/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
}

func TestGetBookById_Found(t *testing.T) {
	mock := &MockStore{
		GetFunc: func(id int) (models.Book, error) {
			return models.Book{ID: id, Title: "Test Book"}, nil
		},
	}

	handler := &handlers.BookHandler{Store: mock}

	req := httptest.NewRequest(http.MethodGet, "/books/123", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestGetBookById_NotFound(t *testing.T) {
	mock := &MockStore{
		GetFunc: func(id int) (models.Book, error) {
			return models.Book{}, errors.New("not found")
		},
	}

	handler := &handlers.BookHandler{Store: mock}

	req := httptest.NewRequest(http.MethodGet, "/books/999", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestDeleteBook(t *testing.T) {
	mock := &MockStore{
		DeleteFunc: func(id int) error {
			return nil
		},
	}

	handler := &handlers.BookHandler{Store: mock}

	req := httptest.NewRequest(http.MethodDelete, "/books/1", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestSearchBooks(t *testing.T) {
	mock := &MockStore{
		SearchFunc: func(c models.SearchCriteria) ([]models.Book, error) {
			if c.Title == "Go" {
				return []models.Book{
					{ID: 1, Title: "Go"},
				}, nil
			}
			return []models.Book{}, nil
		},
	}

	handler := &handlers.BookHandler{Store: mock}

	req := httptest.NewRequest(http.MethodGet, "/books/?title=Go", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var results []models.Book
	if err := json.NewDecoder(rr.Body).Decode(&results); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(results) != 1 || results[0].Title != "Go" {
		t.Fatalf("unexpected search result: %+v", results)
	}
}
