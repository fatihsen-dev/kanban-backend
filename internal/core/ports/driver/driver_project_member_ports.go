package ports

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
)

type ProjectMemberService interface {
	CreateProjectMember(ctx context.Context, projectMember *domain.ProjectMember) error
	GetProjectMembersByProjectID(ctx context.Context, projectID string) ([]*domain.ProjectMember, error)
	DeleteProjectMemberByID(ctx context.Context, id string) error
}
