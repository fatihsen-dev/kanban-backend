package ports

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
)

type TaskService interface {
	CreateTask(ctx context.Context, task *domain.Task) error
	GetTaskByID(ctx context.Context, id string) (*domain.Task, error)
	GetTasks(ctx context.Context) ([]*domain.Task, error)
}
