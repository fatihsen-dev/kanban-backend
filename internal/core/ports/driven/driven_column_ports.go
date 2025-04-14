package ports

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
)

type ColumnRepository interface {
	Save(ctx context.Context, column *domain.Column) error
	GetByID(ctx context.Context, id string) (*domain.Column, error)
	GetAll(ctx context.Context) ([]*domain.Column, error)
	GetColumnsByProjectID(ctx context.Context, projectID string) ([]*domain.Column, error)
	Update(ctx context.Context, column *domain.Column) error
	Delete(ctx context.Context, id string) error
}
