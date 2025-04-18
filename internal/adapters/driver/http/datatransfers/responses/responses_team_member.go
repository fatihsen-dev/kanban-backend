package responses

type TeamMemberResponse struct {
	ID        string `json:"id"`
	TeamID    string `json:"team_id"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

type DeleteTeamMemberResponse struct {
	ID string `json:"id"`
}
