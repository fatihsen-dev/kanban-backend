package ports

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
)

type ColumnService interface {
	CreateColumn(ctx context.Context, column *domain.Column) error
	GetColumnByID(ctx context.Context, id string) (*domain.Column, error)
	GetColumnsByProjectID(ctx context.Context, projectID string) ([]*domain.Column, error)
	GetColumnWithDetails(ctx context.Context, columnID string) (*domain.Column, []*domain.Task, error)
	UpdateColumn(ctx context.Context, column *domain.Column) error
	DeleteColumn(ctx context.Context, id string) error
}
