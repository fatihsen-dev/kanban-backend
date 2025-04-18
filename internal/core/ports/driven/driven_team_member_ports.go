package ports

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
)

type TeamMemberRepository interface {
	Save(ctx context.Context, teamMember *domain.TeamMember) error
	GetTeamMembersByTeamID(ctx context.Context, teamID string) ([]*domain.TeamMember, error)
	DeleteByID(ctx context.Context, id string) error
}
