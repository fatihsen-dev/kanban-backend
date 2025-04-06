package service

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driven"
)

type ProjectService struct {
	projectRepo ports.ProjectRepository
}

func NewProjectService(projectRepo ports.ProjectRepository) *ProjectService {
	return &ProjectService{projectRepo: projectRepo}
}

func (s *ProjectService) CreateProject(ctx context.Context, project *domain.Project) error {
	err := s.projectRepo.Save(ctx, project)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProjectService) GetProjectByID(ctx context.Context, id string) (*domain.Project, error) {
	return s.projectRepo.GetByID(ctx, id)
}

func (s *ProjectService) GetProjects(ctx context.Context) ([]*domain.Project, error) {
	return s.projectRepo.GetAll(ctx)
}
