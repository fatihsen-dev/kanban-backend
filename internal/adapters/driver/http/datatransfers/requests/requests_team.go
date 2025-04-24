package requests

type CreateTeamRequest struct {
	Name      string `json:"name" validate:"required,min=3,max=26"`
	Role      string `json:"role" validate:"required,oneof=admin write read"`
	ProjectID string `json:"project_id" validate:"required,uuid4"`
}

type UpdateTeamRequest struct {
	Role      *string `json:"role,omitempty" validate:"oneof=admin write read"`
	Name      *string `json:"name,omitempty" validate:"min=3,max=26"`
	ProjectID *string `json:"project_id,omitempty" validate:"uuid4"`
}
