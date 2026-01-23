package store

import (
	"Book-Store/internal/models"
	"context"
	"errors"
)

func (s *MemStore) CreateAuthor(ctx context.Context, author models.Author) (models.Author, error) {
	select {
	case <-ctx.Done():
		return models.Author{}, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	maxID := -1
	for id := range s.Authors {
		if id > maxID {
			maxID = id
		}
	}

	author.ID = maxID + 1
	s.Authors[author.ID] = author

	if err := s.SaveToFile(); err != nil {
		return models.Author{}, err
	}

	return author, nil
}

func (s *MemStore) GetAuthor(ctx context.Context, id int) (models.Author, error) {
	select {
	case <-ctx.Done():
		return models.Author{}, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	author, exists := s.Authors[id]
	if !exists {
		return models.Author{}, errors.New("Author not found")
	}
	return author, nil
}

func (s *MemStore) ListAuthors(ctx context.Context) ([]models.Author, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	authors := make([]models.Author, 0)
	for _, a := range s.Authors {
		authors = append(authors, a)
	}
	return authors, nil
}

func (s *MemStore) UpdateAuthor(ctx context.Context, id int, author models.Author) (models.Author, error) {
	select {
	case <-ctx.Done():
		return models.Author{}, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Authors[id]; !exists {
		return models.Author{}, errors.New("Author not found")
	}

	author.ID = id
	s.Authors[id] = author

	if err := s.SaveToFile(); err != nil {
		return models.Author{}, err
	}

	return author, nil
}

func (s *MemStore) DeleteAuthor(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Authors[id]; !exists {
		return errors.New("Author not found")
	}

	delete(s.Authors, id)

	if err := s.SaveToFile(); err != nil {
		return err
	}

	return nil
}

func (s *MemStore) AuthorExists(id int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.Authors[id]
	return exists
}

func (s *MemStore) AuthorsCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.Authors)
}

func (s *MemStore) BooksPerAuthor() map[int]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[int]int)
	for _, book := range s.Books {
		result[book.Author.ID]++
	}
	return result
}
