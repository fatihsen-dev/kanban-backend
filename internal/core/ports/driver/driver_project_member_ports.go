package ports

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
)

type ProjectMemberService interface {
	CreateProjectMember(ctx context.Context, projectMember *domain.ProjectMember) error
	GetProjectMembersByProjectID(ctx context.Context, projectID string, query *string) ([]*domain.ProjectMember, []*domain.User, error)
	DeleteProjectMemberByID(ctx context.Context, id string) error
	GetByUserIDAndProjectID(ctx context.Context, userID, projectID string) (*domain.ProjectMember, error)
	UpdateProjectMember(ctx context.Context, projectMember *domain.ProjectMember) error
}
