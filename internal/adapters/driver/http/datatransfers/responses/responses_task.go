package responses

type TaskResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	ProjectID string `json:"project_id"`
	ColumnID  string `json:"column_id"`
	CreatedAt string `json:"created_at"`
}
