package db

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driven"
	"github.com/lib/pq"
)

type PostgresProjectRepository struct {
	PostgresRepository
}

func NewPostgresProjectRepo(baseRepo *PostgresRepository) ports.ProjectRepository {
	return &PostgresProjectRepository{PostgresRepository: *baseRepo}
}

func (r *PostgresProjectRepository) Save(ctx context.Context, project *domain.Project) error {
	query := `INSERT INTO projects (name, owner_id) VALUES ($1, $2) RETURNING id, name, owner_id, created_at`
	err := r.DB.QueryRowContext(ctx, query, project.Name, project.OwnerID).Scan(&project.ID, &project.Name, &project.OwnerID, &project.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresProjectRepository) GetByID(ctx context.Context, id string) (*domain.Project, error) {
	query := `SELECT id, name, owner_id, created_at FROM projects WHERE id = $1`
	var project domain.Project
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&project.ID, &project.Name, &project.OwnerID, &project.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *PostgresProjectRepository) GetUserProjects(ctx context.Context, userID string) ([]*domain.Project, error) {
	query := `SELECT id, name, owner_id, created_at FROM projects WHERE owner_id = $1 ORDER BY created_at ASC`
	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*domain.Project
	for rows.Next() {
		var project domain.Project
		err := rows.Scan(&project.ID, &project.Name, &project.OwnerID, &project.CreatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, &project)
	}

	return projects, nil
}

func (r *PostgresProjectRepository) GetByIDs(ctx context.Context, ids []string) ([]*domain.Project, error) {
	query := `SELECT id, name, owner_id, created_at FROM projects WHERE id = ANY($1)`
	rows, err := r.DB.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*domain.Project
	for rows.Next() {
		var project domain.Project
		err := rows.Scan(&project.ID, &project.Name, &project.OwnerID, &project.CreatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, &project)
	}
	return projects, nil
}
