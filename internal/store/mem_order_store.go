package store

import (
	"Book-Store/internal/models"
	"context"
	"errors"
	"time"
)

func (s *MemStore) CreateOrder(ctx context.Context, order models.Order) (models.Order, error) {
	select {
	case <-ctx.Done():
		return models.Order{}, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Customers[order.Customer.ID]; !exists {
		return models.Order{}, errors.New("customer not found")
	}

	var totalPrice float64
	for i, item := range order.Items {
		select {
		case <-ctx.Done():
			return models.Order{}, ctx.Err()
		default:
		}

		book, exists := s.Books[item.Book.ID]
		if !exists {
			return models.Order{}, errors.New("book not found in order")
		}
		if book.Stock < item.Quantity {
			return models.Order{}, errors.New("insufficient stock")
		}
		book.Stock -= item.Quantity
		s.Books[book.ID] = book
		order.Items[i].Book = book
		totalPrice += book.Price * float64(item.Quantity)
	}

	maxID := -1
	for id := range s.Orders {
		if id > maxID {
			maxID = id
		}
	}

	order.ID = maxID + 1
	order.Status = "created"
	order.TotalPrice = totalPrice
	order.CreatedAt = time.Now()
	s.Orders[order.ID] = order

	if err := s.SaveToFile(); err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (s *MemStore) GetOrder(ctx context.Context, id int) (models.Order, error) {
	select {
	case <-ctx.Done():
		return models.Order{}, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	order, exists := s.Orders[id]
	if !exists {
		return models.Order{}, errors.New("order not found")
	}
	return order, nil
}

func (s *MemStore) CompleteOrder(ctx context.Context, id int) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.Orders[id]
	if !exists {
		return false, errors.New("order not found")
	}
	if order.Status != "created" {
		return false, errors.New("order cannot be completed")
	}
	order.Status = "completed"
	s.Orders[id] = order
	return true, s.SaveToFile()
}

func (s *MemStore) CancelOrder(ctx context.Context, id int) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.Orders[id]
	if !exists {
		return false, errors.New("order not found")
	}
	if order.Status != "created" {
		return false, errors.New("order cannot be cancelled")
	}
	for _, item := range order.Items {
		if book, exists := s.Books[item.Book.ID]; exists {
			book.Stock += item.Quantity
			s.Books[book.ID] = book
		}
	}
	order.Status = "cancelled"
	s.Orders[id] = order
	return true, s.SaveToFile()
}

func (s *MemStore) SearchOrderByStatus(ctx context.Context, status string) ([]models.Order, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	orders := make([]models.Order, 0)
	for _, order := range s.Orders {
		if order.Status == status {
			orders = append(orders, order)
		}
	}
	return orders, nil
}

func (s *MemStore) ListOrders(ctx context.Context) ([]models.Order, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	orders := make([]models.Order, 0)
	for _, order := range s.Orders {
		orders = append(orders, order)
	}
	return orders, nil
}

func (s *MemStore) GetOrdersInTimeRange(ctx context.Context, start, end time.Time) ([]models.Order, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	orders := make([]models.Order, 0)
	startUnix := start.Unix()
	endUnix := end.Unix()

	for _, order := range s.Orders {
		orderUnix := order.CreatedAt.Unix()
		if orderUnix >= startUnix && orderUnix <= endUnix {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

