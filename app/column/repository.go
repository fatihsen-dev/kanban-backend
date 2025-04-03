package column

import (
	"context"

	models "github.com/fatihsen-dev/kanban-backend/domain"
)

type Repository interface {
	CreateColumn(ctx context.Context, column *models.Column) error
	GetColumn(ctx context.Context, id string) (*models.Column, error)
}
