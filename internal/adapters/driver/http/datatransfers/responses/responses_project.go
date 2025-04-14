package responses

type ProjectResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type ProjectWithDetailsResponse struct {
	ID        string                    `json:"id"`
	Name      string                    `json:"name"`
	CreatedAt string                    `json:"created_at"`
	Columns   []ColumnWithTasksResponse `json:"columns"`
}
