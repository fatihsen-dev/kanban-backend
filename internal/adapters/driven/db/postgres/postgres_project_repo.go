package db

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driven"
)

type PostgresProjectRepository struct {
	PostgresRepository
}

func NewPostgresProjectRepo(baseRepo *PostgresRepository) ports.ProjectRepository {
	return &PostgresProjectRepository{PostgresRepository: *baseRepo}
}

func (r *PostgresProjectRepository) Save(ctx context.Context, project *domain.Project) error {
	query := `INSERT INTO projects (name) VALUES ($1) RETURNING id, name, created_at`
	err := r.DB.QueryRowContext(ctx, query, project.Name).Scan(&project.ID, &project.Name, &project.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresProjectRepository) GetByID(ctx context.Context, id string) (*domain.Project, error) {
	query := `SELECT id, name FROM projects WHERE id = $1`
	var project domain.Project
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&project.ID, &project.Name)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *PostgresProjectRepository) GetAll(ctx context.Context) ([]*domain.Project, error) {
	query := `SELECT id, name, created_at FROM projects`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*domain.Project
	for rows.Next() {
		var project domain.Project
		err := rows.Scan(&project.ID, &project.Name, &project.CreatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, &project)
	}

	return projects, nil
}
