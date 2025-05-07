package ports

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
)

type ProjectMemberRepository interface {
	Save(ctx context.Context, projectMember *domain.ProjectMember) error
	GetProjectMembersByProjectID(ctx context.Context, projectID string) ([]*domain.ProjectMember, error)
	DeleteByID(ctx context.Context, id string) error
	GetByUserIDAndProjectID(ctx context.Context, userID, projectID string) (*domain.ProjectMember, error)
	GetByUserID(ctx context.Context, userID string) ([]*domain.ProjectMember, error)
	UpdateProjectMember(ctx context.Context, projectMember *domain.ProjectMember) error
}
