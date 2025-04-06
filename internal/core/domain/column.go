package domain

import "time"

type Column struct {
	ID        string
	Name      string
	ProjectID string
	CreatedAt time.Time
}
