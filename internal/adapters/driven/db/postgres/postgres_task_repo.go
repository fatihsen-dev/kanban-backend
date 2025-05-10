package db

import (
	"context"
	"fmt"
	"strings"

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
	queryBase := `INSERT INTO tasks (title, column_id, project_id`
	valuesBase := `VALUES ($1, $2, $3`
	returningBase := `RETURNING id, title, content, column_id, project_id`
	args := []interface{}{task.Title, task.ColumnID, task.ProjectID}
	paramIndex := 4

	if task.Content != nil {
		queryBase += `, content`
		valuesBase += fmt.Sprintf(`, $%d`, paramIndex)
		returningBase += `, content`
		args = append(args, *task.Content)
		paramIndex++
	}

	queryBase += `) `
	valuesBase += `) `
	returningBase += `, created_at`

	query := queryBase + valuesBase + returningBase

	scanArgs := []interface{}{&task.ID, &task.Title, &task.Content, &task.ColumnID, &task.ProjectID}
	scanArgs = append(scanArgs, &task.CreatedAt)

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(scanArgs...)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresTaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	query := `SELECT id, title, content, column_id, project_id, created_at FROM tasks WHERE id = $1`
	var task domain.Task
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&task.ID, &task.Title, &task.Content, &task.ColumnID, &task.ProjectID, &task.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *PostgresTaskRepository) GetTasksByColumnIDs(ctx context.Context, columnIDs []string) ([]*domain.Task, error) {
	query := `SELECT id, title, content, column_id, project_id, created_at FROM tasks WHERE column_id = ANY($1) ORDER BY created_at ASC`
	rows, err := r.DB.QueryContext(ctx, query, pq.Array(columnIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		var task domain.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Content, &task.ColumnID, &task.ProjectID, &task.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}

func (r *PostgresTaskRepository) GetAll(ctx context.Context) ([]*domain.Task, error) {
	query := `SELECT id, title, content, column_id, project_id, created_at FROM tasks ORDER BY created_at ASC`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		var task domain.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Content, &task.ColumnID, &task.ProjectID, &task.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}

func (r *PostgresTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	queryBase := "UPDATE tasks SET "
	queryWhere := " WHERE id = $%d"

	setClauses := []string{}
	args := []interface{}{}
	paramIndex := 1

	if task.Title != "" {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", paramIndex))
		args = append(args, task.Title)
		paramIndex++
	}

	if task.Content != nil {
		setClauses = append(setClauses, fmt.Sprintf("content = $%d", paramIndex))
		args = append(args, task.Content)
		paramIndex++
	}

	if task.ColumnID != "" {
		setClauses = append(setClauses, fmt.Sprintf("column_id = $%d", paramIndex))
		args = append(args, task.ColumnID)
		paramIndex++
	}

	if len(setClauses) == 0 {
		return nil
	}

	querySet := strings.Join(setClauses, ", ")

	args = append(args, task.ID)

	finalQuery := queryBase + querySet + fmt.Sprintf(queryWhere, paramIndex)

	_, err := r.DB.ExecContext(ctx, finalQuery, args...)
	if err != nil {
		return fmt.Errorf("task update failed: %w", err)
	}

	return nil
}

func (r *PostgresTaskRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
