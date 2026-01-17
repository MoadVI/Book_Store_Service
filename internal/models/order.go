package models

import "time"

type OrderItem struct {
	Book     Book `json:"book"`
	Quantity int  `json:"quantity"`
}

type Order struct {
	ID         int         `json:"id"`
	Customer   Customer    `json:"customer"`
	Items      []OrderItem `json:"items"`
	TotalPrice float64     `json:"total_price"`
	CreatedAt  time.Time   `json:"created_at"`
	Status     string      `json:"status"`
}
