package store

import (
	"Book-Store/internal/models"
	"errors"
)

func (s *MemStore) CreateAuthor(author models.Author) (models.Author, error) {
	s.mu.Lock()

	maxID := -1
	for id := range s.Authors {
		if id > maxID {
			maxID = id
		}
	}

	author.ID = maxID + 1
	s.Authors[author.ID] = author

	s.mu.Unlock()
	if err := s.SaveToFile(); err != nil {
		return models.Author{}, err
	}

	return author, nil

}

func (s *MemStore) GetAuthor(id int) (models.Author, error) {
	s.mu.RLock()

	author, exists := s.Authors[id]
	if !exists {
		return models.Author{}, errors.New("Author not found")
	}
	return author, nil
}

func (s *MemStore) ListAuthors() ([]models.Author, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	authors := make([]models.Author, 0)
	for _, a := range s.Authors {
		authors = append(authors, a)
	}

	return authors, nil
}

func (s *MemStore) UpdateAuthor(id int, author models.Author) (models.Author, error) {
	s.mu.Lock()
	if _, exists := s.Authors[id]; !exists {
		s.mu.Unlock()
		return models.Author{}, errors.New("Author not found")
	}

	author.ID = id
	s.Authors[id] = author

	s.mu.Unlock()
	_ = s.SaveToFile()

	return author, nil

}

func (s *MemStore) DeleteAuthor(id int) error {
	s.mu.Lock()

	if _, exists := s.Authors[id]; !exists {
		s.mu.Unlock()
		return errors.New("Author not found")
	}

	delete(s.Authors, id)
	s.mu.Unlock()

	_ = s.SaveToFile()
	return nil

}
