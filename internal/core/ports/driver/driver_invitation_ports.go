package ports

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers/requests"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers/responses"
	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
)

type InvitationService interface {
	GetInvitationByID(ctx context.Context, id string) (*domain.Invitation, error)
	CreateInvitations(ctx context.Context, invitations []*domain.Invitation) ([]responses.InvitationResponse, error)
	GetInvitations(ctx context.Context, userID string) ([]responses.InvitationResponse, error)
	UpdateInvitationStatus(ctx context.Context, request requests.InvitationUpdateStatusRequest) (*domain.ProjectMember, *domain.User, error)
}
