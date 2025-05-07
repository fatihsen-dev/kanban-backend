package service

import (
	"context"

	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driven"
)

type ProjectService struct {
	projectRepo       ports.ProjectRepository
	projectMemberRepo ports.ProjectMemberRepository
	teamRepo          ports.TeamRepository
	columnRepo        ports.ColumnRepository
	taskRepo          ports.TaskRepository
	userRepo          ports.UserRepository
}

func NewProjectService(projectRepo ports.ProjectRepository, columnRepo ports.ColumnRepository, taskRepo ports.TaskRepository, teamRepo ports.TeamRepository, projectMemberRepo ports.ProjectMemberRepository, userRepo ports.UserRepository) *ProjectService {
	return &ProjectService{
		projectRepo:       projectRepo,
		columnRepo:        columnRepo,
		taskRepo:          taskRepo,
		teamRepo:          teamRepo,
		projectMemberRepo: projectMemberRepo,
		userRepo:          userRepo,
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, project *domain.Project) error {
	err := s.projectRepo.Save(ctx, project)

	if err != nil {
		return err
	}

	projectMember := &domain.ProjectMember{
		ProjectID: project.ID,
		UserID:    project.OwnerID,
		Role:      domain.AccessOwnerRole,
	}

	columns := []*domain.Column{
		{
			ProjectID: project.ID,
			Name:      "To Do",
		},
		{
			ProjectID: project.ID,
			Name:      "In Progress",
		},
		{
			ProjectID: project.ID,
			Name:      "Done",
		},
	}

	for _, column := range columns {
		err = s.columnRepo.Save(ctx, column)
		if err != nil {
			return err
		}
	}

	err = s.projectMemberRepo.Save(ctx, projectMember)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProjectService) GetProjectByID(ctx context.Context, id string) (*domain.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) GetUserProjects(ctx context.Context, userID string) ([]*domain.Project, error) {
	projectMembers, err := s.projectMemberRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	projectIDs := make([]string, len(projectMembers))
	for i, projectMember := range projectMembers {
		projectIDs[i] = projectMember.ProjectID
	}

	return s.projectRepo.GetByIDs(ctx, projectIDs)
}

func (s *ProjectService) GetProjectWithDetails(ctx context.Context, projectID string) (*domain.Project, []*domain.Column, map[string][]*domain.Task, []*domain.Team, []*domain.ProjectMember, []*domain.User, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	teams, err := s.teamRepo.GetTeamsByProjectID(ctx, projectID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	projectMembers, err := s.projectMemberRepo.GetProjectMembersByProjectID(ctx, projectID, nil)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	userIDs := make([]string, len(projectMembers))
	for i, projectMember := range projectMembers {
		userIDs[i] = projectMember.UserID
	}

	users, err := s.userRepo.GetByIDs(ctx, userIDs)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	userMap := make(map[string]*domain.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	columns, err := s.columnRepo.GetColumnsByProjectID(ctx, projectID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	columnIDs := make([]string, len(columns))
	for i, column := range columns {
		columnIDs[i] = column.ID
	}

	tasks, err := s.taskRepo.GetTasksByColumnIDs(ctx, columnIDs)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	tasksByColumn := make(map[string][]*domain.Task)
	for _, task := range tasks {
		tasksByColumn[task.ColumnID] = append(tasksByColumn[task.ColumnID], task)
	}

	return project, columns, tasksByColumn, teams, projectMembers, users, nil
}
