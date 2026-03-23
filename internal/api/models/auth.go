package models

type DummyLoginRequest struct {
	Role *string `json:"role" validate:"required,oneof=admin user"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    *string `json:"role" validate:"required"`
	Password *string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Email    *string `json:"email" validate:"required"`
	Password *string `json:"password" validate:"required"`
	Role     *string `json:"role" validate:"required,oneof=admin user"`
}

type RegisterResponse struct {
	User *UserObject `json:"user"`
}
