package entity

import "time"

type Order struct {
	ID         string
	UserID     string
	UserEmail  string
	TotalPrice float64
	Status     string
	CreatedAt  time.Time
	Items      []OrderItem
}

type OrderItem struct {
	ID      string
	OrderID string
	PetID   string
	Price   float64
}
