package responses

type ProjectMemberResponse struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	TeamID    string `json:"team_id"`
	CreatedAt string `json:"created_at"`
}

type DeleteProjectMemberResponse struct {
	ID string `json:"id"`
}
