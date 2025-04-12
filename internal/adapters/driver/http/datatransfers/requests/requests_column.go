package requests

type ColumnCreateRequest struct {
	Name      string `json:"name" validate:"required,min=3,max=26"`
	ProjectID string `json:"project_id" validate:"required"`
}

type ColumnUpdateRequest struct {
	Name string `json:"name" validate:"required,min=3,max=26"`
}
