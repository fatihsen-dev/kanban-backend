package requests

type CreateProjectMemberRequest struct {
	UserID    string `json:"user_id" validate:"required,uuid4"`
	ProjectID string `json:"project_id" validate:"required,uuid4"`
	TeamID    string `json:"team_id,omitempty" validate:"uuid4"`
}
