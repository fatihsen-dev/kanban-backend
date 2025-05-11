package requests

type ProjectCreateRequest struct {
	Name string `json:"name" validate:"required,min=3,max=26,notblank"`
}
