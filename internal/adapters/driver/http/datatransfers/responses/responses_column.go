package responses

type ColumnResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ProjectID string `json:"project_id"`
	CreatedAt string `json:"created_at"`
}

type ColumnWithTasksResponse struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	CreatedAt string         `json:"created_at"`
	Tasks     []TaskResponse `json:"tasks"`
}
