package postgres

import (
	"context"
	"fmt"

	"github.com/fatihsen-dev/kanban-backend/domain"
	"github.com/jackc/pgx/v5"
)

func (pg *postgres) InsertItem(ctx context.Context, item *domain.Item) error {
	query := `INSERT INTO items (title) VALUES (@title) RETURNING id, title, created_at`

	args := pgx.NamedArgs{
		"title": item.Title,
	}

	err := pg.db.QueryRow(ctx, query, args).Scan(&item.ID, &item.Title, &item.CreatedAt)
	if err != nil {
		return fmt.Errorf("unable to insert and retrieve row: %w", err)
	}

	return nil
}
