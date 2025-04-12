package service

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driven"
)

type ColumnService struct {
	columnRepo ports.ColumnRepository
	taskRepo   ports.TaskRepository
}

func NewColumnService(columnRepo ports.ColumnRepository, taskRepo ports.TaskRepository) *ColumnService {
	return &ColumnService{columnRepo: columnRepo, taskRepo: taskRepo}
}

func (s *ColumnService) CreateColumn(ctx context.Context, column *domain.Column) error {
	err := s.columnRepo.Save(ctx, column)
	if err != nil {
		return err
	}
	return nil
}

func (s *ColumnService) GetColumnByID(ctx context.Context, id string) (*domain.Column, error) {
	return s.columnRepo.GetByID(ctx, id)
}

func (s *ColumnService) GetColumns(ctx context.Context) ([]*domain.Column, error) {
	return s.columnRepo.GetAll(ctx)
}

func (s *ColumnService) GetColumnWithTasks(ctx context.Context, columnID string) (*domain.Column, []*domain.Task, error) {
	column, err := s.columnRepo.GetByID(ctx, columnID)
	if err != nil {
		return nil, nil, err
	}

	tasks, err := s.taskRepo.GetTasksByColumnIDs(ctx, []string{columnID})
	if err != nil {
		return nil, nil, err
	}

	return column, tasks, nil
}

func (s *ColumnService) UpdateColumn(ctx context.Context, column *domain.Column) error {
	return s.columnRepo.Update(ctx, column)
}

func (s *ColumnService) DeleteColumn(ctx context.Context, id string) error {
	err := s.taskRepo.DeleteTasksByColumnID(ctx, id)
	if err != nil {
		return err
	}

	err = s.columnRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
