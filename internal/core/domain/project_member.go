package domain

import "time"

type ProjectMemberRole string

const (
	ProjectOwnerRole ProjectMemberRole = "owner"
	ProjectAdminRole ProjectMemberRole = "admin"
	ProjectWriteRole ProjectMemberRole = "write"
	ProjectReadRole  ProjectMemberRole = "read"
)

type ProjectMember struct {
	ID        string
	TeamID    *string
	UserID    string
	ProjectID string
	Role      ProjectMemberRole
	CreatedAt time.Time
}
