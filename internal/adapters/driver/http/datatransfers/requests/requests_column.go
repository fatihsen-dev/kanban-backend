package requests

type ColumnCreateRequest struct {
	Name      string  `json:"name" validate:"required,min=3,max=26"`
	Color     *string `json:"color" validate:"omitempty,hexcolor"`
	ProjectID string  `json:"project_id" validate:"required,uuid4"`
}

type ColumnUpdateRequest struct {
	Name  *string `json:"name,omitempty" validate:"omitempty,min=3,max=26"`
	Color *string `json:"color,omitempty" validate:"omitempty,hexcolor|eq="`
}

type ColumnDeleteRequest struct {
	ColumnID  string `json:"column_id" validate:"required,uuid4"`
	ProjectID string `json:"project_id" validate:"required,uuid4"`
}
