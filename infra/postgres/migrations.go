package postgres

import (
	"context"
	"fmt"
)

func Migration(ctx context.Context, pg *postgres) {
	_, err := pg.db.Query(ctx, "CREATE TABLE items (id SERIAL PRIMARY KEY, title TEXT NOT NULL, created_at TIMESTAMP DEFAULT now());")

	if err != nil {
		panic(fmt.Sprint("Migration error: %w", err))
	}
}
