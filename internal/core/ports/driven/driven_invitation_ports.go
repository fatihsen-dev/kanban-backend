package ports

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
)

type InvitationRepository interface {
	GetInvitations(ctx context.Context, userID string) ([]*domain.Invitation, error)
	SaveInvitations(ctx context.Context, invitations []*domain.Invitation) error
}
