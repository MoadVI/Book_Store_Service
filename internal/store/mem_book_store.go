package store

import (
	"Book-Store/internal/models"
	"context"
	"errors"
	"slices"
	"sort"
	"strings"
)

func (s *MemStore) CreateBook(ctx context.Context, book models.Book) (models.Book, error) {
	select {
	case <-ctx.Done():
		return models.Book{}, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	author, ok := s.Authors[book.Author.ID]
	if !ok {
		return models.Book{}, errors.New("author not found")
	}

	book.Author.FirstName = author.FirstName
	book.Author.LastName = author.LastName
	book.Author.Bio = author.Bio

	maxID := -1
	for id := range s.Books {
		if id > maxID {
			maxID = id
		}
	}

	book.ID = maxID + 1
	s.Books[book.ID] = book

	if err := s.SaveToFile(); err != nil {
		return models.Book{}, err
	}

	return book, nil
}

func (s *MemStore) GetBook(ctx context.Context, id int) (models.Book, error) {
	select {
	case <-ctx.Done():
		return models.Book{}, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	book, exists := s.Books[id]
	if !exists {
		return models.Book{}, errors.New("book not found")
	}
	return book, nil
}

func (s *MemStore) UpdateBook(ctx context.Context, id int, book models.Book) (models.Book, error) {
	select {
	case <-ctx.Done():
		return models.Book{}, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Books[id]; !exists {
		return models.Book{}, errors.New("book not found")
	}

	author, ok := s.Authors[book.Author.ID]
	if !ok {
		return models.Book{}, errors.New("author not found")
	}

	book.Author.FirstName = author.FirstName
	book.Author.LastName = author.LastName
	book.Author.Bio = author.Bio

	book.ID = id
	s.Books[id] = book

	if err := s.SaveToFile(); err != nil {
		return models.Book{}, err
	}

	return book, nil
}

func (s *MemStore) DeleteBook(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Books[id]; !exists {
		return errors.New("book not found")
	}

	delete(s.Books, id)

	if err := s.SaveToFile(); err != nil {
		return err
	}

	return nil
}

func (s *MemStore) SearchBooks(ctx context.Context, criteria models.SearchCriteria) ([]models.Book, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]models.Book, 0)
	for _, b := range s.Books {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if criteria.Title != "" && !strings.Contains(strings.ToLower(b.Title), strings.ToLower(criteria.Title)) {
			continue
		}

		if criteria.Author != "" {
			searchWords := strings.Fields(strings.ToLower(criteria.Author))
			first := strings.ToLower(b.Author.FirstName)
			last := strings.ToLower(b.Author.LastName)
			matched := false
			for _, word := range searchWords {
				if strings.Contains(first, word) || strings.Contains(last, word) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		if criteria.Genre != "" && !slices.Contains(b.Genres, criteria.Genre) {
			continue
		}

		if criteria.MinPrice != nil && b.Price < *criteria.MinPrice {
			continue
		}

		if criteria.MaxPrice != nil && b.Price > *criteria.MaxPrice {
			continue
		}

		results = append(results, b)
	}

	if criteria.SortBy != "" {
		switch strings.ToLower(criteria.SortBy) {
		case "title":
			sort.Slice(results, func(i, j int) bool {
				if strings.ToLower(criteria.SortOrder) == "desc" {
					return results[i].Title > results[j].Title
				}
				return results[i].Title < results[j].Title
			})
		case "price":
			sort.Slice(results, func(i, j int) bool {
				if strings.ToLower(criteria.SortOrder) == "desc" {
					return results[i].Price > results[j].Price
				}
				return results[i].Price < results[j].Price
			})
		}
	}

	return results, nil
}

func (s *MemStore) BookExists(id int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.Books[id]
	return exists
}
