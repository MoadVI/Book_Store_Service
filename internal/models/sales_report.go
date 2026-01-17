package models

import "time"

type BookSales struct {
	Book     Book `json:"book"`
	Quantity int  `json:"quantity"`
}

type SalesReport struct {
	Timestamp      time.Time   `json:"timestamp"`
	TotalRevenue   float64     `json:"total_revenue"`
	TotalOrders    int         `json:"total_orders"`
	TopSellingBook []BookSales `json:"top_selling_books"`
}
