package responses

type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	IsAdmin   bool   `json:"is_admin"`
	CreatedAt string `json:"created_at"`
}

type UserLoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserRegisterResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserAuthResponse struct {
	User UserResponse `json:"user"`
}
