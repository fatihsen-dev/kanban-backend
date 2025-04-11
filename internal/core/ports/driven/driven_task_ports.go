package ports

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
)

type TaskRepository interface {
	Save(ctx context.Context, task *domain.Task) error
	GetByID(ctx context.Context, id string) (*domain.Task, error)
	GetAll(ctx context.Context) ([]*domain.Task, error)
	GetTasksByColumnIDs(ctx context.Context, columnIDs []string) ([]*domain.Task, error)
}
