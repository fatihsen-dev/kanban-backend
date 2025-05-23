package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driven"
)

type PostgresProjectMemberRepository struct {
	PostgresRepository
}

func NewPostgresProjectMemberRepo(baseRepo *PostgresRepository) ports.ProjectMemberRepository {
	return &PostgresProjectMemberRepository{PostgresRepository: *baseRepo}
}

func (r *PostgresProjectMemberRepository) Save(ctx context.Context, projectMember *domain.ProjectMember) error {
	if projectMember.TeamID == nil {
		query := `INSERT INTO project_members (user_id, project_id, role) VALUES ($1, $2, $3) RETURNING id, user_id, project_id, role, created_at`
		err := r.DB.QueryRowContext(ctx, query, projectMember.UserID, projectMember.ProjectID, projectMember.Role).Scan(&projectMember.ID, &projectMember.UserID, &projectMember.ProjectID, &projectMember.Role, &projectMember.CreatedAt)
		if err != nil {
			return err
		}
	} else {
		query := `INSERT INTO project_members (user_id, project_id, team_id, role) VALUES ($1, $2, $3, $4) RETURNING id, user_id, project_id, team_id, role, created_at`
		err := r.DB.QueryRowContext(ctx, query, projectMember.UserID, projectMember.ProjectID, projectMember.TeamID, projectMember.Role).Scan(&projectMember.ID, &projectMember.UserID, &projectMember.ProjectID, &projectMember.TeamID, &projectMember.Role, &projectMember.CreatedAt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresProjectMemberRepository) GetProjectMembersByProjectID(ctx context.Context, projectID string, searchQuery *string) ([]*domain.ProjectMember, error) {
	var query string
	var args []interface{}

	if searchQuery != nil && *searchQuery != "" {
		query = `
			SELECT DISTINCT pm.id, pm.team_id, pm.user_id, pm.project_id, pm.role, pm.created_at 
			FROM project_members pm
			INNER JOIN users u ON pm.user_id = u.id
			WHERE pm.project_id = $1 
			AND (
				LOWER(u.email) LIKE LOWER($2)
			)`
		args = []interface{}{projectID, "%" + *searchQuery + "%"}
	} else {
		query = `SELECT id, team_id, user_id, project_id, role, created_at 
				FROM project_members 
				WHERE project_id = $1`
		args = []interface{}{projectID}
	}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projectMembers []*domain.ProjectMember
	for rows.Next() {
		var projectMember domain.ProjectMember
		err := rows.Scan(&projectMember.ID, &projectMember.TeamID, &projectMember.UserID, &projectMember.ProjectID, &projectMember.Role, &projectMember.CreatedAt)
		if err != nil {
			return nil, err
		}
		projectMembers = append(projectMembers, &projectMember)
	}
	return projectMembers, nil
}

func (r *PostgresProjectMemberRepository) DeleteByID(ctx context.Context, id string) error {
	query := `DELETE FROM project_members WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresProjectMemberRepository) GetByUserIDAndProjectID(ctx context.Context, userID, projectID string) (*domain.ProjectMember, error) {
	query := `SELECT id, team_id, user_id, project_id, role, created_at FROM project_members WHERE user_id = $1 AND project_id = $2`
	var projectMember domain.ProjectMember
	err := r.DB.QueryRowContext(ctx, query, userID, projectID).Scan(&projectMember.ID, &projectMember.TeamID, &projectMember.UserID, &projectMember.ProjectID, &projectMember.Role, &projectMember.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &projectMember, nil
}

func (r *PostgresProjectMemberRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.ProjectMember, error) {
	query := `SELECT id, team_id, user_id, project_id, role, created_at FROM project_members WHERE user_id = $1`
	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projectMembers []*domain.ProjectMember
	for rows.Next() {
		var projectMember domain.ProjectMember
		err := rows.Scan(&projectMember.ID, &projectMember.TeamID, &projectMember.UserID, &projectMember.ProjectID, &projectMember.Role, &projectMember.CreatedAt)
		if err != nil {
			return nil, err
		}
		projectMembers = append(projectMembers, &projectMember)
	}
	return projectMembers, nil
}

func (r *PostgresProjectMemberRepository) UpdateProjectMember(ctx context.Context, projectMember *domain.ProjectMember) error {
	queryBase := "UPDATE project_members SET "
	queryWhere := " WHERE id = $%d"

	setClauses := []string{}
	args := []interface{}{}
	paramIndex := 1

	if projectMember.Role != "" {
		setClauses = append(setClauses, fmt.Sprintf("role = $%d", paramIndex))
		args = append(args, projectMember.Role)
		paramIndex++
	}

	if projectMember.TeamID != nil {
		teamID := *projectMember.TeamID
		if teamID == "" {
			setClauses = append(setClauses, "team_id = NULL")
		} else {
			setClauses = append(setClauses, fmt.Sprintf("team_id = $%d", paramIndex))
			args = append(args, teamID)
			paramIndex++
		}
	}

	if len(setClauses) == 0 {
		return nil
	}

	querySet := strings.Join(setClauses, ", ")

	args = append(args, projectMember.ID)

	finalQuery := queryBase + querySet + fmt.Sprintf(queryWhere, paramIndex)

	_, err := r.DB.ExecContext(ctx, finalQuery, args...)
	if err != nil {
		return fmt.Errorf("project member update failed: %w", err)
	}

	return nil
}
