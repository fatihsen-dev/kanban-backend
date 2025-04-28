package ports

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByIDs(ctx context.Context, ids []string) ([]*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	GetUsersByQuery(ctx context.Context, query string) ([]*domain.User, error)
}
