package service

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driven"
)

type ColumnService struct {
	columnRepo ports.ColumnRepository
}

func NewColumnService(columnRepo ports.ColumnRepository) *ColumnService {
	return &ColumnService{columnRepo: columnRepo}
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
