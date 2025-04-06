package ports

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
)

type ColumnService interface {
	CreateColumn(ctx context.Context, column *domain.Column) error
	GetColumnByID(ctx context.Context, id string) (*domain.Column, error)
	GetColumns(ctx context.Context) ([]*domain.Column, error)
}
