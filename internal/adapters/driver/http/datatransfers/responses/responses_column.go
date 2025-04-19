package responses

type ColumnResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ProjectID string `json:"project_id"`
	CreatedAt string `json:"created_at"`
}

type ColumnUpdateResponse struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type ColumnDeleteResponse struct {
	ID string `json:"id"`
}

type ColumnWithDetailsResponse struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	CreatedAt string         `json:"created_at"`
	Tasks     []TaskResponse `json:"tasks"`
}
