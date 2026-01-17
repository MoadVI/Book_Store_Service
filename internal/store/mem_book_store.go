package store

import (
	"Book-Store/internal/models"
	"errors"
)

func (s *MemStore) CreateBook(book models.Book) (models.Book, error) {
	s.mu.Lock()

	maxID := -1
	for id := range s.Books {
		if id > maxID {
			maxID = id
		}
	}

	book.ID = maxID + 1
	s.Books[book.ID] = book

	s.mu.Unlock()

	_ = s.SaveToFile()

	return book, nil
}

func (s *MemStore) GetBook(id int) (models.Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	book, exists := s.Books[id]
	if !exists {
		return models.Book{}, errors.New("book not found")
	}

	return book, nil
}

func (s *MemStore) UpdateBook(id int, book models.Book) (models.Book, error) {
	s.mu.Lock()

	if _, exists := s.Books[id]; !exists {
		s.mu.Unlock()
		return models.Book{}, errors.New("book not found")
	}

	book.ID = id
	s.Books[id] = book

	s.mu.Unlock()

	_ = s.SaveToFile()

	return book, nil
}

func (s *MemStore) SearchBooks(criteria models.SearchCriteria) ([]models.Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]models.Book, 0)
	for _, b := range s.Books {
		if criteria.Title != "" && b.Title != criteria.Title {
			continue
		}
		results = append(results, b)
	}

	return results, nil
}

func (s *MemStore) DeleteBook(id int) error {
	s.mu.Lock()

	if _, exists := s.Books[id]; !exists {
		s.mu.Unlock()
		return errors.New("book not found")
	}

	delete(s.Books, id)

	s.mu.Unlock()

	_ = s.SaveToFile()

	return nil
}
