package auth

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=40"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=12,max=40,containsuppercase,containslowercase,containsnumber,containsspecial"`
	Username string `json:"username" validate:"required,min=5,max=30,alphanumunicode"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}
