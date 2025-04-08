package domain

import "time"

type User struct {
	ID           string
	Name         string
	Email        string
	IsAdmin      bool
	PasswordHash string
	CreatedAt    time.Time
}
