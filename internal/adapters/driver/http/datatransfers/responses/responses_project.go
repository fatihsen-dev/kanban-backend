package responses

type ProjectResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	OwnerID   string `json:"owner_id"`
	CreatedAt string `json:"created_at"`
}

type ProjectWithDetailsResponse struct {
	ID        string                          `json:"id"`
	Name      string                          `json:"name"`
	OwnerID   string                          `json:"owner_id"`
	CreatedAt string                          `json:"created_at"`
	Columns   []ColumnWithDetailsResponse     `json:"columns"`
	Teams     []TeamResponse                  `json:"teams"`
	Members   []ProjectMemberWithUserResponse `json:"members"`
}
