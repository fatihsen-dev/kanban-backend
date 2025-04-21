package domain

import "time"

type TeamRole string

const (
	TeamOwnerRole TeamRole = "owner"
	TeamAdminRole TeamRole = "admin"
	TeamWriteRole TeamRole = "write"
	TeamReadRole  TeamRole = "read"
)

type Team struct {
	ID        string
	Name      string
	Role      TeamRole
	ProjectID string
	CreatedAt time.Time
}
