package reports

import (
	"Book-Store/internal/models"
	"Book-Store/internal/store"
	"context"
	"time"
)

func GenerateSalesReport(ctx context.Context, orderStore store.OrderStore) (*models.SalesReport, error) {
	orders, err := orderStore.ListOrders(ctx)
	if err != nil {
		return nil, err
	}

	report := &models.SalesReport{
		Timestamp:      time.Now(),
		TopSellingBook: make([]models.BookSales, 0),
	}

	bookSalesMap := make(map[int]*models.BookSales)

	for _, order := range orders {
		report.TotalOrders++

		if order.Status == "completed" {
			report.TotalRevenue += order.TotalPrice

			for _, item := range order.Items {
				if bs, exists := bookSalesMap[item.Book.ID]; exists {
					bs.Quantity += item.Quantity
				} else {
					bookSalesMap[item.Book.ID] = &models.BookSales{
						Book:     item.Book,
						Quantity: item.Quantity,
					}
				}
			}
		}
	}

	for _, bs := range bookSalesMap {
		report.TopSellingBook = append(report.TopSellingBook, *bs)
	}

	sortTopSellingBooks(report.TopSellingBook)

	if len(report.TopSellingBook) > 10 {
		report.TopSellingBook = report.TopSellingBook[:10]
	}

	return report, nil
}

func sortTopSellingBooks(books []models.BookSales) {
	for i := 0; i < len(books)-1; i++ {
		for j := i + 1; j < len(books); j++ {
			if books[j].Quantity > books[i].Quantity {
				books[i], books[j] = books[j], books[i]
			}
		}
	}
}
