package store

import (
	"Book-Store/internal/models"
	"errors"
	"time"
)

func (s *MemStore) CreateCustomer(customer models.Customer) (models.Customer, error) {
	s.mu.Lock()

	maxID := -1
	for id := range s.Customers {
		if id > maxID {
			maxID = id
		}
	}

	customer.ID = maxID + 1
	if customer.CreatedAt.IsZero() {
		customer.CreatedAt = time.Now()
	}
	s.Customers[customer.ID] = customer

	s.mu.Unlock()
	if err := s.SaveToFile(); err != nil {
		return models.Customer{}, err
	}

	return customer, nil
}

func (s *MemStore) GetCustomer(id int) (models.Customer, error) {
	s.mu.RLock()
	customer, exists := s.Customers[id]
	s.mu.RUnlock()
	if !exists {
		return models.Customer{}, errors.New("Customer not found")
	}

	return customer, nil

}

func (s *MemStore) UpdateCustomer(id int, customer models.Customer) (models.Customer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Customers[id]; !exists {
		return models.Customer{}, errors.New("Customre not found")
	}

	customer.ID = id
	s.Customers[id] = customer

	if err := s.SaveToFile(); err != nil {
		return models.Customer{}, err
	}

	return customer, nil
}

func (s *MemStore) ListCustomers() ([]models.Customer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	customers := make([]models.Customer, 0)
	for _, customer := range s.Customers {
		customers = append(customers, customer)
	}

	return customers, nil

}

func (s *MemStore) DeleteCustomer(id int) error {
	s.mu.Lock()

	_, exists := s.Customers[id]
	if !exists {
		s.mu.Unlock()
		return errors.New("Customer not found")
	}

	delete(s.Customers, id)
	s.mu.Unlock()

	if err := s.SaveToFile(); err != nil {
		return err
	}

	return nil

}
