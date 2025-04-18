package ports

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
)

type TeamRepository interface {
	Save(ctx context.Context, team *domain.Team) error
	GetByID(ctx context.Context, id string) (*domain.Team, error)
	GetTeamsByProjectID(ctx context.Context, projectID string) ([]*domain.Team, error)
	DeleteByID(ctx context.Context, id string) error
	Update(ctx context.Context, team *domain.Team) error
}
