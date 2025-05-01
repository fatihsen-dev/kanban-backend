package ports

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
)

type InvitationRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Invitation, error)
	GetInvitations(ctx context.Context, userID string) ([]*domain.Invitation, error)
	SaveInvitations(ctx context.Context, invitations []*domain.Invitation) error
	UpdateStatus(ctx context.Context, id string, status string) error
}
