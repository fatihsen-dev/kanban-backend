package requests

type InvitationCreateRequest struct {
	InviterID  string   `json:"inviter_id" validate:"required,uuid4"`
	InviteeIDs []string `json:"invitee_ids" validate:"required,dive,uuid4"`
	ProjectID  string   `json:"project_id" validate:"required,uuid4"`
	Message    *string  `json:"message" validate:"omitempty,min=3,max=100"`
}

type InvitationUpdateStatusRequest struct {
	ID     string `json:"id" validate:"required,uuid4"`
	Status string `json:"status" validate:"required,oneof=accepted rejected"`
	UserID string `json:"user_id" validate:"required,uuid4"`
}
