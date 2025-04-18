package ports

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
)

type TeamService interface {
	CreateTeam(ctx context.Context, team *domain.Team) error
	GetTeamWithMembersByID(ctx context.Context, id string) (*domain.Team, []*domain.TeamMember, error)
	GetTeamsByProjectID(ctx context.Context, projectID string) ([]*domain.Team, error)
	UpdateTeam(ctx context.Context, team *domain.Team) error
	DeleteTeamByID(ctx context.Context, id string) error
	CreateTeamMember(ctx context.Context, teamMember *domain.TeamMember) error
	DeleteTeamMemberByID(ctx context.Context, teamID string, memberID string) error
}
