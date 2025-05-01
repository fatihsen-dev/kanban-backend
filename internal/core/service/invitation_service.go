package service

import (
	"context"
	"errors"
	"time"

	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers/requests"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers/responses"
	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driven"
)

type InvitationService struct {
	invitationRepo    ports.InvitationRepository
	userRepo          ports.UserRepository
	projectRepo       ports.ProjectRepository
	projectMemberRepo ports.ProjectMemberRepository
}

func NewInvitationService(invitationRepo ports.InvitationRepository, userRepo ports.UserRepository, projectRepo ports.ProjectRepository, projectMemberRepo ports.ProjectMemberRepository) *InvitationService {
	return &InvitationService{
		invitationRepo:    invitationRepo,
		userRepo:          userRepo,
		projectRepo:       projectRepo,
		projectMemberRepo: projectMemberRepo,
	}
}

func (s *InvitationService) buildInvitationResponse(invitation *domain.Invitation, invitee *domain.User, inviter *domain.User, project *domain.Project) responses.InvitationResponse {
	return responses.InvitationResponse{
		ID: invitation.ID,
		Inviter: responses.UserResponse{
			ID:        inviter.ID,
			Name:      inviter.Name,
			Email:     inviter.Email,
			IsAdmin:   inviter.IsAdmin,
			CreatedAt: inviter.CreatedAt.Format(time.RFC3339),
		},
		Invitee: responses.UserResponse{
			ID:        invitee.ID,
			Name:      invitee.Name,
			Email:     invitee.Email,
			IsAdmin:   invitee.IsAdmin,
			CreatedAt: invitee.CreatedAt.Format(time.RFC3339),
		},
		Project: responses.ProjectResponse{
			ID:        project.ID,
			Name:      project.Name,
			OwnerID:   project.OwnerID,
			CreatedAt: project.CreatedAt.Format(time.RFC3339),
		},
		Message:   invitation.Message,
		Status:    string(invitation.Status),
		CreatedAt: invitation.CreatedAt.Format(time.RFC3339),
	}
}

func (s *InvitationService) CreateInvitations(ctx context.Context, invitations []*domain.Invitation) ([]responses.InvitationResponse, error) {
	err := s.invitationRepo.SaveInvitations(ctx, invitations)
	if err != nil {
		return nil, err
	}

	userIDs := make(map[string]struct{})
	projectIDs := make(map[string]struct{})

	for _, invitation := range invitations {
		if invitation.ID != "" {
			userIDs[invitation.InviteeID] = struct{}{}
			userIDs[invitation.InviterID] = struct{}{}
			projectIDs[invitation.ProjectID] = struct{}{}
		}
	}

	userIDList := make([]string, 0, len(userIDs))
	for id := range userIDs {
		userIDList = append(userIDList, id)
	}

	projectIDList := make([]string, 0, len(projectIDs))
	for id := range projectIDs {
		projectIDList = append(projectIDList, id)
	}

	users, err := s.userRepo.GetByIDs(ctx, userIDList)
	if err != nil {
		return nil, err
	}

	projects, err := s.projectRepo.GetByIDs(ctx, projectIDList)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]*domain.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	projectMap := make(map[string]*domain.Project)
	for _, project := range projects {
		projectMap[project.ID] = project
	}

	responseData := make([]responses.InvitationResponse, 0)
	for _, invitation := range invitations {
		if invitation.ID != "" {
			invitee := userMap[invitation.InviteeID]
			inviter := userMap[invitation.InviterID]
			project := projectMap[invitation.ProjectID]
			responseData = append(responseData, s.buildInvitationResponse(invitation, invitee, inviter, project))
		}
	}

	return responseData, nil
}

func (s *InvitationService) GetInvitations(ctx context.Context, userID string) ([]responses.InvitationResponse, error) {
	invitations, err := s.invitationRepo.GetInvitations(ctx, userID)
	if err != nil {
		return nil, err
	}

	userIDs := make(map[string]struct{})
	projectIDs := make(map[string]struct{})

	for _, invitation := range invitations {
		if invitation.ID != "" {
			userIDs[invitation.InviteeID] = struct{}{}
			userIDs[invitation.InviterID] = struct{}{}
			projectIDs[invitation.ProjectID] = struct{}{}
		}
	}

	userIDList := make([]string, 0, len(userIDs))
	for id := range userIDs {
		userIDList = append(userIDList, id)
	}

	projectIDList := make([]string, 0, len(projectIDs))
	for id := range projectIDs {
		projectIDList = append(projectIDList, id)
	}

	users, err := s.userRepo.GetByIDs(ctx, userIDList)
	if err != nil {
		return nil, err
	}

	projects, err := s.projectRepo.GetByIDs(ctx, projectIDList)
	if err != nil {
		return nil, err
	}

	userMap := make(map[string]*domain.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	projectMap := make(map[string]*domain.Project)
	for _, project := range projects {
		projectMap[project.ID] = project
	}

	responseData := make([]responses.InvitationResponse, 0)
	for _, invitation := range invitations {
		if invitation.ID != "" {
			invitee := userMap[invitation.InviteeID]
			inviter := userMap[invitation.InviterID]
			project := projectMap[invitation.ProjectID]
			responseData = append(responseData, s.buildInvitationResponse(invitation, invitee, inviter, project))
		}
	}

	return responseData, nil
}

func (s *InvitationService) GetInvitationByID(ctx context.Context, id string) (*domain.Invitation, error) {
	return s.invitationRepo.GetByID(ctx, id)
}

func (s *InvitationService) UpdateInvitationStatus(ctx context.Context, request requests.InvitationUpdateStatusRequest) error {
	invitation, err := s.GetInvitationByID(ctx, request.ID)
	if err != nil {
		return err
	}

	if invitation.InviteeID != request.UserID {
		return errors.New("you are not allowed to update this invitation")
	}

	if invitation.Status != domain.InvitationStatusPending {
		return errors.New("this invitation has already been accepted or rejected")
	}

	err = s.projectMemberRepo.Save(ctx, &domain.ProjectMember{
		UserID:    request.UserID,
		ProjectID: invitation.ProjectID,
		Role:      domain.AccessReadRole,
	})

	if err != nil {
		return err
	}

	err = s.invitationRepo.UpdateStatus(ctx, request.ID, request.Status)
	if err != nil {
		return err
	}

	return nil
}
