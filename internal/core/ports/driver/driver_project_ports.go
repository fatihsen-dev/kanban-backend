package ports

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
)

type ProjectService interface {
	CreateProject(ctx context.Context, project *domain.Project) error
	GetProjectByID(ctx context.Context, id string) (*domain.Project, error)
	GetProjects(ctx context.Context) ([]*domain.Project, error)
}
