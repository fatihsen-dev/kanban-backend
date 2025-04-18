package requests

type CreateTeamMemberRequest struct {
	UserID    string `json:"user_id" validate:"required,uuid4"`
	ProjectID string `json:"project_id" validate:"required,uuid4"`
	TeamID    string `json:"team_id" validate:"required,uuid4"`
}

type DeleteTeamMemberRequest struct {
	ProjectID string `json:"project_id" validate:"required,uuid4"`
	TeamID    string `json:"team_id" validate:"required,uuid4"`
	MemberID  string `json:"member_id" validate:"required,uuid4"`
}
