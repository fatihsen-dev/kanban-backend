package domain

import "time"

type ProjectMember struct {
	ID        string
	TeamID    string
	UserID    string
	ProjectID string
	Role      string // Owner, Admin, Write, Read
	CreatedAt time.Time
}
