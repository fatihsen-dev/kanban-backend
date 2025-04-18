package responses

type TeamResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	ProjectID string `json:"project_id"`
	CreatedAt string `json:"created_at"`
}

type UpdateTeamResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name,omitempty"`
	Role      string `json:"role,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
}

type DeleteTeamResponse struct {
	ID string `json:"id"`
}

type TeamWithMembersResponse struct {
	ID        string               `json:"id"`
	Name      string               `json:"name"`
	Role      string               `json:"role"`
	ProjectID string               `json:"project_id"`
	CreatedAt string               `json:"created_at"`
	Members   []TeamMemberResponse `json:"members"`
}
