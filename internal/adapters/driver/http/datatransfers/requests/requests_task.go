package requests

type TaskCreateRequest struct {
	Title     string `json:"title" validate:"required,min=3,max=26"`
	ColumnID  string `json:"column_id" validate:"required,uuid4"`
	ProjectID string `json:"project_id" validate:"required,uuid4"`
}

type TaskUpdateRequest struct {
	Title    *string `json:"title,omitempty" validate:"omitempty,min=3,max=26"`
	ColumnID *string `json:"column_id,omitempty" validate:"omitempty,uuid4"`
}
