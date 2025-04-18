package domain

import "time"

type Team struct {
	ID        string
	Name      string
	Role      string // Admin, Write, Read
	ProjectID string
	CreatedAt time.Time
}
