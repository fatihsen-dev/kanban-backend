package ports

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
)

type ProjectService interface {
	CreateProject(ctx context.Context, project *domain.Project) error
	GetProjectByID(ctx context.Context, id string) (*domain.Project, error)
	GetUserProjects(ctx context.Context, userID string) ([]*domain.Project, error)
	GetProjectWithDetails(ctx context.Context, projectID string) (*domain.Project, []*domain.Column, map[string][]*domain.Task, []*domain.Team, []*domain.ProjectMember, []*domain.User, error)
}
