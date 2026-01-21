package store

import (
	"Book-Store/internal/models"
	"context"
	"errors"
	"time"
)

func (s *MemStore) CreateCustomer(ctx context.Context, customer models.Customer) (models.Customer, error) {
	select {
	case <-ctx.Done():
		return models.Customer{}, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

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

	if err := s.SaveToFile(); err != nil {
		return models.Customer{}, err
	}

	return customer, nil
}

func (s *MemStore) GetCustomer(ctx context.Context, id int) (models.Customer, error) {
	select {
	case <-ctx.Done():
		return models.Customer{}, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	customer, exists := s.Customers[id]
	if !exists {
		return models.Customer{}, errors.New("Customer not found")
	}

	return customer, nil
}

func (s *MemStore) UpdateCustomer(ctx context.Context, id int, customer models.Customer) (models.Customer, error) {
	select {
	case <-ctx.Done():
		return models.Customer{}, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Customers[id]; !exists {
		return models.Customer{}, errors.New("Customer not found")
	}

	customer.ID = id
	s.Customers[id] = customer

	if err := s.SaveToFile(); err != nil {
		return models.Customer{}, err
	}

	return customer, nil
}

func (s *MemStore) ListCustomers(ctx context.Context) ([]models.Customer, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	customers := make([]models.Customer, 0)
	for _, customer := range s.Customers {
		customers = append(customers, customer)
	}
	return customers, nil
}

func (s *MemStore) DeleteCustomer(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.Customers[id]
	if !exists {
		return errors.New("Customer not found")
	}

	delete(s.Customers, id)

	if err := s.SaveToFile(); err != nil {
		return err
	}

	return nil
}

func (s *MemStore) CustomerExists(id int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.Customers[id]
	return exists
}

