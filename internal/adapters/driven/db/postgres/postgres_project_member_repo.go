package db

import (
	"context"

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
		query := `INSERT INTO project_members (user_id, project_id, role) VALUES ($1, $2, $3)`
		_, err := r.DB.ExecContext(ctx, query, projectMember.UserID, projectMember.ProjectID, projectMember.Role)
		if err != nil {
			return err
		}
	} else {
		query := `INSERT INTO project_members (user_id, project_id, team_id, role) VALUES ($1, $2, $3, $4)`
		_, err := r.DB.ExecContext(ctx, query, projectMember.UserID, projectMember.ProjectID, projectMember.TeamID, projectMember.Role)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresProjectMemberRepository) GetProjectMembersByProjectID(ctx context.Context, projectID string) ([]*domain.ProjectMember, error) {
	query := `SELECT id, team_id, user_id, project_id, role, created_at FROM project_members WHERE project_id = $1`
	rows, err := r.DB.QueryContext(ctx, query, projectID)
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
