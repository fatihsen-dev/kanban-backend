package service

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driven"
)

type ProjectMemberService struct {
	projectMemberRepo ports.ProjectMemberRepository
	userRepo          ports.UserRepository
}

func NewProjectMemberService(projectMemberRepo ports.ProjectMemberRepository, userRepo ports.UserRepository) *ProjectMemberService {
	return &ProjectMemberService{projectMemberRepo: projectMemberRepo, userRepo: userRepo}
}

func (s *ProjectMemberService) CreateProjectMember(ctx context.Context, projectMember *domain.ProjectMember) error {
	return s.projectMemberRepo.Save(ctx, projectMember)
}

func (s *ProjectMemberService) DeleteProjectMemberByID(ctx context.Context, id string) error {
	return s.projectMemberRepo.DeleteByID(ctx, id)
}

func (s *ProjectMemberService) GetProjectMembersByProjectID(ctx context.Context, projectID string) ([]*domain.ProjectMember, error) {
	return s.projectMemberRepo.GetProjectMembersByProjectID(ctx, projectID)
}
