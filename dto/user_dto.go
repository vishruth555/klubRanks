package dto

type SignupRequest struct {
	Username string `json:"username" binding:"required" example:"john"`
	AvatarID string `json:"avatar_id" binding:"required"`
	Password string `json:"password" binding:"required" example:"strongpassword"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"john"`
	AvatarID string `json:"avatar_id"`
	Password string `json:"password" binding:"required" example:"strongpassword"`
}

type MessageResponse struct {
	Message string `json:"message" example:"user created successfully"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"could not parse data"`
}

type LoginResponse struct {
	Message string `json:"message" example:"login successful"`
	Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
