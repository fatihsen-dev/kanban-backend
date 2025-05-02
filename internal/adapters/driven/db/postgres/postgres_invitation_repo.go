package db

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driven"
)

type PostgresInvitationRepository struct {
	PostgresRepository
}

func NewPostgresInvitationRepo(baseRepo *PostgresRepository) ports.InvitationRepository {
	return &PostgresInvitationRepository{PostgresRepository: *baseRepo}
}

func (r *PostgresInvitationRepository) SaveInvitations(ctx context.Context, invitations []*domain.Invitation) error {

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, invitation := range invitations {
		if invitation.InviterID == invitation.InviteeID {
			continue
		}

		checkQuery := `SELECT COUNT(*) FROM invitations WHERE invitee_id = $1 AND project_id = $2 AND status = 'pending'`
		var count int
		err := tx.QueryRowContext(ctx, checkQuery, invitation.InviteeID, invitation.ProjectID).Scan(&count)
		if err != nil {
			return err
		}

		if count > 0 {
			continue
		}

		if invitation.Message == nil {
			query := `INSERT INTO invitations (inviter_id, invitee_id, project_id, status) VALUES ($1, $2, $3, $4) RETURNING id`
			err := tx.QueryRowContext(ctx, query, invitation.InviterID, invitation.InviteeID, invitation.ProjectID, invitation.Status).Scan(&invitation.ID)
			if err != nil {
				return err
			}
		} else {
			query := `INSERT INTO invitations (inviter_id, invitee_id, project_id, message, status) VALUES ($1, $2, $3, $4, $5) RETURNING id`
			err := tx.QueryRowContext(ctx, query, invitation.InviterID, invitation.InviteeID, invitation.ProjectID, invitation.Message, invitation.Status).Scan(&invitation.ID)
			if err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *PostgresInvitationRepository) GetInvitations(ctx context.Context, userID string) ([]*domain.Invitation, error) {
	query := `SELECT * FROM invitations WHERE invitee_id = $1 AND status = 'pending' ORDER BY created_at DESC`
	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invitations := []*domain.Invitation{}
	for rows.Next() {
		var invitation domain.Invitation
		err := rows.Scan(&invitation.ID, &invitation.InviterID, &invitation.InviteeID, &invitation.ProjectID, &invitation.Message, &invitation.Status, &invitation.CreatedAt)
		if err != nil {
			return nil, err
		}
		invitations = append(invitations, &invitation)
	}

	return invitations, nil
}

func (r *PostgresInvitationRepository) GetByID(ctx context.Context, id string) (*domain.Invitation, error) {
	query := `SELECT * FROM invitations WHERE id = $1`
	var invitation domain.Invitation
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&invitation.ID, &invitation.InviterID, &invitation.InviteeID, &invitation.ProjectID, &invitation.Message, &invitation.Status, &invitation.CreatedAt)
	if err != nil {

		return nil, err
	}

	return &invitation, nil
}

func (r *PostgresInvitationRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE invitations SET status = $1 WHERE id = $2`
	_, err := r.DB.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}

	return nil
}
