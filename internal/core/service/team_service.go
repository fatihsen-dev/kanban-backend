package service

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driven"
)

type TeamService struct {
	teamRepo       ports.TeamRepository
	teamMemberRepo ports.TeamMemberRepository
	userRepo       ports.UserRepository
}

func NewTeamService(teamRepo ports.TeamRepository, teamMemberRepo ports.TeamMemberRepository, userRepo ports.UserRepository) *TeamService {
	return &TeamService{teamRepo: teamRepo, teamMemberRepo: teamMemberRepo, userRepo: userRepo}
}

func (s *TeamService) CreateTeam(ctx context.Context, team *domain.Team) error {
	return s.teamRepo.Save(ctx, team)
}

func (s *TeamService) GetTeamWithMembersByID(ctx context.Context, id string) (*domain.Team, []*domain.TeamMember, error) {
	team, err := s.teamRepo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	teamMembers, err := s.teamMemberRepo.GetTeamMembersByTeamID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	return team, teamMembers, nil
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

func (s *TeamService) CreateTeamMember(ctx context.Context, teamMember *domain.TeamMember) error {
	return s.teamMemberRepo.Save(ctx, teamMember)
}

func (s *TeamService) DeleteTeamMemberByID(ctx context.Context, memberID string) error {
	return s.teamMemberRepo.DeleteByID(ctx, memberID)
}
