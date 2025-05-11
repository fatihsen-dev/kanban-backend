package responses

type TaskResponse struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Content   *string `json:"content"`
	ProjectID string  `json:"project_id"`
	ColumnID  string  `json:"column_id"`
	CreatedAt string  `json:"created_at"`
}

type TaskUpdateResponse struct {
	ID       string  `json:"id"`
	Title    string  `json:"title,omitempty"`
	Content  *string `json:"content,omitempty"`
	ColumnID string  `json:"column_id,omitempty"`
}

type TaskDeleteResponse struct {
	ID string `json:"id"`
}
