package store

import (
	"Book-Store/internal/models"
	"encoding/json"
	"os"
	"sync"
)

type MemStore struct {
	mu      sync.RWMutex
	dbPath  string
	Books   map[int]models.Book   `json:"books"`
	Authors map[int]models.Author `json:"authors"`
	Orders  map[int]models.Order  `json:"orders"`
}

func NewMemStore() *MemStore {
	return &MemStore{
		Books:   make(map[int]models.Book),
		Authors: make(map[int]models.Author),
		Orders:  make(map[int]models.Order),
	}
}

func (s *MemStore) SaveToFile() error {
	s.mu.Lock()
	path := s.dbPath
	s.mu.Unlock()

	if path == "" {
		path = getDBPath()

	}

	data, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (s *MemStore) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			s.mu.Lock()
			s.dbPath = path
			s.mu.Unlock()
			return nil
		}
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.dbPath = path
	return json.Unmarshal(data, s)
}

func getDBPath() string {
	path := os.Getenv("DB_PATH")
	if path == "" {
		return "database.json"
	}
	return path
}
