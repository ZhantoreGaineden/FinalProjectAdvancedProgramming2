package entity

import "time"

type Pet struct {
	ID        string
	Name      string
	Category  string
	Breed     string
	Age       int32
	Price     float64
	Status    string
	CreatedAt time.Time
}
