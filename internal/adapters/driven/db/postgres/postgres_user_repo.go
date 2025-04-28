package db

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driven"
	"github.com/lib/pq"
)

type PostgresUserRepository struct {
	PostgresRepository
}

func NewPostgresUserRepo(baseRepo *PostgresRepository) ports.UserRepository {
	return &PostgresUserRepository{PostgresRepository: *baseRepo}
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (name, email, password_hash, is_admin) VALUES ($1, $2, $3, $4) RETURNING id, name, email, password_hash, is_admin, created_at`
	err := r.DB.QueryRowContext(ctx, query, user.Name, user.Email, user.PasswordHash, user.IsAdmin).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.IsAdmin, &user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, name, email, password_hash, is_admin, created_at FROM users WHERE id = $1`
	var user domain.User
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.IsAdmin, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, name, email, password_hash, is_admin, created_at FROM users WHERE email = $1`
	var user domain.User
	err := r.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.IsAdmin, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	query := `SELECT id, name, email, password_hash, is_admin, created_at FROM users`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.IsAdmin, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *PostgresUserRepository) GetByIDs(ctx context.Context, ids []string) ([]*domain.User, error) {
	query := `SELECT id, name, email, password_hash, is_admin, created_at FROM users WHERE id = ANY($1)`
	rows, err := r.DB.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.IsAdmin, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *PostgresUserRepository) GetUsersByQuery(ctx context.Context, queryString string) ([]*domain.User, error) {
	query := `SELECT id, name, email, password_hash, is_admin, created_at FROM users WHERE name ILIKE $1 OR email ILIKE $2`
	rows, err := r.DB.QueryContext(ctx, query, "%"+queryString+"%", "%"+queryString+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.IsAdmin, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
