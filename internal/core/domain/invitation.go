package domain

import "time"

type InvitationStatus string

const (
	InvitationStatusPending  InvitationStatus = "pending"
	InvitationStatusAccepted InvitationStatus = "accepted"
	InvitationStatusRejected InvitationStatus = "rejected"
)

type Invitation struct {
	ID        string
	InviterID string
	InviteeID string
	ProjectID string
	Message   *string
	Status    InvitationStatus
	CreatedAt time.Time
}
