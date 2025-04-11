package db

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driven"
	"github.com/lib/pq"
)

type PostgresTaskRepository struct {
	PostgresRepository
}

func NewPostgresTaskRepo(baseRepo *PostgresRepository) ports.TaskRepository {
	return &PostgresTaskRepository{PostgresRepository: *baseRepo}
}

func (r *PostgresTaskRepository) Save(ctx context.Context, task *domain.Task) error {
	query := `INSERT INTO tasks (title, column_id, project_id) VALUES ($1, $2, $3) RETURNING id, title, column_id, project_id, created_at`
	err := r.DB.QueryRowContext(ctx, query, task.Title, task.ColumnID, task.ProjectID).Scan(&task.ID, &task.Title, &task.ColumnID, &task.ProjectID, &task.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresTaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	query := `SELECT id, title, column_id, project_id, created_at FROM tasks WHERE id = $1`
	var task domain.Task
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&task.ID, &task.Title, &task.ColumnID, &task.ProjectID, &task.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *PostgresTaskRepository) GetTasksByColumnIDs(ctx context.Context, columnIDs []string) ([]*domain.Task, error) {
	query := `SELECT id, title, column_id, project_id, created_at FROM tasks WHERE column_id = ANY($1)`
	rows, err := r.DB.QueryContext(ctx, query, pq.Array(columnIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		var task domain.Task
		err := rows.Scan(&task.ID, &task.Title, &task.ColumnID, &task.ProjectID, &task.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}

func (r *PostgresTaskRepository) GetAll(ctx context.Context) ([]*domain.Task, error) {
	query := `SELECT id, title, column_id, project_id, created_at FROM tasks`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		var task domain.Task
		err := rows.Scan(&task.ID, &task.Title, &task.ColumnID, &task.ProjectID, &task.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}
