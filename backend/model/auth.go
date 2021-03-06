package model

type NewUserForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

type LoginForm struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type GoogleLoginForm struct {
	State string `json:"state" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type LoginResponse struct {
	AuthToken    string      `json:"auth_token"`
	RefreshToken string      `json:"refresh_token"`
	User         *PublicUser `json:"user"`
}

type RefreshAuthTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshAuthTokenResponse struct {
	AuthToken string `json:"auth_token"`
}
