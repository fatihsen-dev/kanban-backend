package requests

type ColumnCreateRequest struct {
	Name      string `json:"name" validate:"required,min=3,max=26"`
	ProjectID string `json:"project_id" validate:"required,uuid4"`
}

type ColumnUpdateRequest struct {
	Name *string `json:"name,omitempty" validate:"omitempty,min=3,max=26"`
}

type ColumnDeleteRequest struct {
	ColumnID  string `json:"column_id" validate:"required,uuid4"`
	ProjectID string `json:"project_id" validate:"required,uuid4"`
}
