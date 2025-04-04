package item

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/domain"
)

type Repository interface {
	InsertItem(ctx context.Context, item *domain.Item) error
}
