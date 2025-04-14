package db

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driven"
)

type PostgresColumnRepository struct {
	PostgresRepository
}

func NewPostgresColumnRepo(baseRepo *PostgresRepository) ports.ColumnRepository {
	return &PostgresColumnRepository{PostgresRepository: *baseRepo}
}

func (r *PostgresColumnRepository) Save(ctx context.Context, column *domain.Column) error {
	query := `INSERT INTO columns (name, project_id) VALUES ($1, $2) RETURNING id, name, project_id, created_at`
	err := r.DB.QueryRowContext(ctx, query, column.Name, column.ProjectID).Scan(&column.ID, &column.Name, &column.ProjectID, &column.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresColumnRepository) GetByID(ctx context.Context, id string) (*domain.Column, error) {
	query := `SELECT id, name, project_id, created_at FROM columns WHERE id = $1`
	var column domain.Column
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&column.ID, &column.Name, &column.ProjectID, &column.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &column, nil
}

func (r *PostgresColumnRepository) GetColumnsByProjectID(ctx context.Context, projectID string) ([]*domain.Column, error) {
	query := `SELECT id, name, project_id, created_at FROM columns WHERE project_id = $1`
	rows, err := r.DB.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []*domain.Column
	for rows.Next() {
		var column domain.Column
		err := rows.Scan(&column.ID, &column.Name, &column.ProjectID, &column.CreatedAt)
		if err != nil {
			return nil, err
		}
		columns = append(columns, &column)
	}
	return columns, nil
}

func (r *PostgresColumnRepository) GetAll(ctx context.Context) ([]*domain.Column, error) {
	query := `SELECT id, name, project_id, created_at FROM columns`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []*domain.Column
	for rows.Next() {
		var column domain.Column
		err := rows.Scan(&column.ID, &column.Name, &column.ProjectID, &column.CreatedAt)
		if err != nil {
			return nil, err
		}
		columns = append(columns, &column)
	}
	return columns, nil
}

func (r *PostgresColumnRepository) Update(ctx context.Context, column *domain.Column) error {
	query := `UPDATE columns SET name = $1 WHERE id = $2`
	_, err := r.DB.ExecContext(ctx, query, column.Name, column.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresColumnRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM columns WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
