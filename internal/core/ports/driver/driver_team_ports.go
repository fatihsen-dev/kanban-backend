package ports

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
)

type TeamService interface {
	CreateTeam(ctx context.Context, team *domain.Team) error
	GetTeamsByProjectID(ctx context.Context, projectID string) ([]*domain.Team, error)
	UpdateTeam(ctx context.Context, team *domain.Team) error
	DeleteTeamByID(ctx context.Context, id string) error
	GetTeamByID(ctx context.Context, id string) (*domain.Team, error)
	AddTeamMembers(ctx context.Context, teamID string, memberIDs []string) ([]string, error)
}
