package domain

import "time"

type Column struct {
	ID        string
	Title     string
	ProjectID string
	CreatedAt time.Time
}
