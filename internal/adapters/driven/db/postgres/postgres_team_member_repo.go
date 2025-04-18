package db

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driven"
)

type PostgresTeamMemberRepository struct {
	PostgresRepository
}

func NewPostgresTeamMemberRepo(baseRepo *PostgresRepository) ports.TeamMemberRepository {
	return &PostgresTeamMemberRepository{PostgresRepository: *baseRepo}
}

func (r *PostgresTeamMemberRepository) Save(ctx context.Context, teamMember *domain.TeamMember) error {
	query := `INSERT INTO team_members (id, team_id, user_id, project_id, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.DB.ExecContext(ctx, query, teamMember.ID, teamMember.TeamID, teamMember.UserID, teamMember.ProjectID, teamMember.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresTeamMemberRepository) GetTeamMembersByTeamID(ctx context.Context, teamID string) ([]*domain.TeamMember, error) {
	query := `SELECT id, team_id, user_id, project_id, created_at FROM team_members WHERE team_id = $1`
	rows, err := r.DB.QueryContext(ctx, query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teamMembers []*domain.TeamMember
	for rows.Next() {
		var teamMember domain.TeamMember
		err := rows.Scan(&teamMember.ID, &teamMember.TeamID, &teamMember.UserID, &teamMember.ProjectID, &teamMember.CreatedAt)
		if err != nil {
			return nil, err
		}
		teamMembers = append(teamMembers, &teamMember)
	}
	return teamMembers, nil
}

func (r *PostgresTeamMemberRepository) DeleteByID(ctx context.Context, id string) error {
	query := `DELETE FROM team_members WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
