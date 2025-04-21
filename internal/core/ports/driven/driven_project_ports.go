package ports

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
)

type ProjectRepository interface {
	Save(ctx context.Context, project *domain.Project) error
	GetByID(ctx context.Context, id string) (*domain.Project, error)
	GetUserProjects(ctx context.Context, userID string) ([]*domain.Project, error)
}
