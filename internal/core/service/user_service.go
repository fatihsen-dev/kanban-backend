package service

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driven"
)

type UserService struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) error {
	err := s.userRepo.Save(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) GetUsersByIDs(ctx context.Context, ids []string) ([]*domain.User, error) {
	users, err := s.userRepo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]*domain.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	usersResponse := make([]*domain.User, len(ids))
	for i, id := range ids {
		usersResponse[i] = userMap[id]
	}
	return usersResponse, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *UserService) GetUsers(ctx context.Context) ([]*domain.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *UserService) GetUsersByQuery(ctx context.Context, query string) ([]*domain.User, error) {
	return s.userRepo.GetUsersByQuery(ctx, query)
}
