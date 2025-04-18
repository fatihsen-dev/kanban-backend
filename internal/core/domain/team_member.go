package domain

import "time"

type TeamMember struct {
	ID        string
	TeamID    string
	UserID    string
	ProjectID string
	CreatedAt time.Time
}
