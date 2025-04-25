package domain

import "time"

type AccessRole string

const (
	AccessOwnerRole AccessRole = "owner"
	AccessAdminRole AccessRole = "admin"
	AccessWriteRole AccessRole = "write"
	AccessReadRole  AccessRole = "read"
)

type Team struct {
	ID        string
	Name      string
	Role      AccessRole
	ProjectID string
	CreatedAt time.Time
}
