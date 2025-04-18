package service

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driven"
)

type TeamService struct {
	teamRepo          ports.TeamRepository
	projectMemberRepo ports.ProjectMemberRepository
	userRepo          ports.UserRepository
}

func NewTeamService(teamRepo ports.TeamRepository, projectMemberRepo ports.ProjectMemberRepository, userRepo ports.UserRepository) *TeamService {
	return &TeamService{teamRepo: teamRepo, projectMemberRepo: projectMemberRepo, userRepo: userRepo}
}

func (s *TeamService) CreateTeam(ctx context.Context, team *domain.Team) error {
	return s.teamRepo.Save(ctx, team)
}

func (s *TeamService) GetTeamWithMembersByID(ctx context.Context, id string) (*domain.Team, []*domain.ProjectMember, error) {
	team, err := s.teamRepo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	projectMembers, err := s.projectMemberRepo.GetProjectMembersByProjectID(ctx, team.ProjectID)
	if err != nil {
		return nil, nil, err
	}

	return team, projectMembers, nil
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

func (s *TeamService) CreateTeamMember(ctx context.Context, projectMember *domain.ProjectMember) error {
	return s.projectMemberRepo.Save(ctx, projectMember)
}

func (s *TeamService) DeleteTeamMemberByID(ctx context.Context, memberID string) error {
	return s.projectMemberRepo.DeleteByID(ctx, memberID)
}
