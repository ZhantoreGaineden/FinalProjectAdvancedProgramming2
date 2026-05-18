package entity

import "time"

type User struct {
	ID           string
	FullName     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
