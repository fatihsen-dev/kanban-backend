package responses

type ProjectMemberResponse struct {
	ID        string  `json:"id"`
	ProjectID string  `json:"project_id"`
	UserID    string  `json:"user_id"`
	Role      string  `json:"role"`
	TeamID    *string `json:"team_id"`
	CreatedAt string  `json:"created_at"`
}

type ProjectMemberWithUserResponse struct {
	ID        string       `json:"id"`
	ProjectID string       `json:"project_id"`
	UserID    string       `json:"user_id"`
	Role      string       `json:"role"`
	TeamID    *string      `json:"team_id"`
	CreatedAt string       `json:"created_at"`
	User      UserResponse `json:"user"`
}

type DeleteProjectMemberResponse struct {
	ID string `json:"id"`
}

type UpdateProjectMemberResponse struct {
	ID     string  `json:"id"`
	Role   *string `json:"role,omitempty"`
	TeamID *string `json:"team_id,omitempty"`
}

type OnlineProjectMembersResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
