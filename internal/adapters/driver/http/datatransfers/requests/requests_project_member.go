package requests

type CreateProjectMemberRequest struct {
	UserID    string `json:"user_id" validate:"required,uuid4"`
	ProjectID string `json:"project_id" validate:"required,uuid4"`
	TeamID    string `json:"team_id,omitempty" validate:"uuid4"`
}

type UpdateProjectMemberRequest struct {
	ID     string  `json:"id" validate:"required,uuid4"`
	Role   *string `json:"role,omitempty" validate:"omitempty,oneof=admin write read"`
	TeamID *string `json:"team_id,omitempty" validate:"omitempty,uuid4|eq="`
}
