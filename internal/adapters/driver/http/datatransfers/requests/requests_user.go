package requests

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=36"`
}

type UserRegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=26"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=36"`
}
