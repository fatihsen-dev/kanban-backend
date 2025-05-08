package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driven"
)

type PostgresTeamRepository struct {
	PostgresRepository
}

func NewPostgresTeamRepo(baseRepo *PostgresRepository) ports.TeamRepository {
	return &PostgresTeamRepository{PostgresRepository: *baseRepo}
}

func (r *PostgresTeamRepository) Save(ctx context.Context, team *domain.Team) error {
	query := `INSERT INTO teams (name, role, project_id) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := r.DB.QueryRowContext(ctx, query, team.Name, team.Role, team.ProjectID).Scan(&team.ID, &team.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresTeamRepository) GetByID(ctx context.Context, id string) (*domain.Team, error) {
	query := `SELECT id, name, role, project_id, created_at FROM teams WHERE id = $1`
	var team domain.Team
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&team.ID, &team.Name, &team.Role, &team.ProjectID, &team.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *PostgresTeamRepository) GetTeamsByProjectID(ctx context.Context, projectID string) ([]*domain.Team, error) {
	query := `SELECT id, name, role, project_id, created_at FROM teams WHERE project_id = $1 ORDER BY created_at ASC`
	rows, err := r.DB.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []*domain.Team
	for rows.Next() {
		var team domain.Team
		err := rows.Scan(&team.ID, &team.Name, &team.Role, &team.ProjectID, &team.CreatedAt)
		if err != nil {
			return nil, err
		}
		teams = append(teams, &team)
	}
	return teams, nil
}

func (r *PostgresTeamRepository) Update(ctx context.Context, team *domain.Team) error {
	queryBase := "UPDATE teams SET "
	queryWhere := " WHERE id = $%d"

	setClauses := []string{}
	args := []interface{}{}
	paramIndex := 1

	if team.Name != "" {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", paramIndex))
		args = append(args, team.Name)
		paramIndex++
	}

	if team.Role != "" {
		setClauses = append(setClauses, fmt.Sprintf("role = $%d", paramIndex))
		args = append(args, team.Role)
		paramIndex++
	}

	if len(setClauses) == 0 {
		return nil
	}

	querySet := strings.Join(setClauses, ", ")

	args = append(args, team.ID)

	finalQuery := queryBase + querySet + fmt.Sprintf(queryWhere, paramIndex)

	_, err := r.DB.ExecContext(ctx, finalQuery, args...)
	if err != nil {
		return fmt.Errorf("team update failed: %w", err)
	}

	return nil
}

func (r *PostgresTeamRepository) DeleteByID(ctx context.Context, id string) error {
	query := `DELETE FROM teams WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
