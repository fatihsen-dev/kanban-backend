package service

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driven"
)

type TeamService struct {
	teamRepo ports.TeamRepository
}

func NewTeamService(teamRepo ports.TeamRepository) *TeamService {
	return &TeamService{teamRepo: teamRepo}
}

func (s *TeamService) CreateTeam(ctx context.Context, team *domain.Team) error {
	return s.teamRepo.Save(ctx, team)
}

func (s *TeamService) UpdateTeam(ctx context.Context, team *domain.Team) error {
	return s.teamRepo.Update(ctx, team)
}

func (s *TeamService) GetTeamsByProjectID(ctx context.Context, projectID string) ([]*domain.Team, error) {
	return s.teamRepo.GetTeamsByProjectID(ctx, projectID)
}

func (s *TeamService) DeleteTeamByID(ctx context.Context, id string) error {
	return s.teamRepo.DeleteByID(ctx, id)
}

func (s *TeamService) GetTeamByID(ctx context.Context, id string) (*domain.Team, error) {
	return s.teamRepo.GetByID(ctx, id)
}
