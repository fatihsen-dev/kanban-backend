package responses

type InvitationResponse struct {
	ID        string          `json:"id"`
	Inviter   UserResponse    `json:"inviter"`
	Invitee   UserResponse    `json:"invitee"`
	Project   ProjectResponse `json:"project"`
	Message   *string         `json:"message"`
	Status    string          `json:"status"`
	CreatedAt string          `json:"created_at"`
}
