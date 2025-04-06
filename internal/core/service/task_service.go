package service

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driven"
)

type TaskService struct {
	taskRepo ports.TaskRepository
}

func NewTaskService(taskRepo ports.TaskRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

func (s *TaskService) CreateTask(ctx context.Context, task *domain.Task) error {
	err := s.taskRepo.Save(ctx, task)
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskService) GetTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	return s.taskRepo.GetByID(ctx, id)
}

func (s *TaskService) GetTasks(ctx context.Context) ([]*domain.Task, error) {
	return s.taskRepo.GetAll(ctx)
}
