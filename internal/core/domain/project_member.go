package domain

import "time"

type ProjectMember struct {
	ID        string
	TeamID    *string
	UserID    string
	ProjectID string
	Role      AccessRole
	CreatedAt time.Time
}
