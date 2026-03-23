package dto

import "github.com/avito-internships/test-backend-1-mmms914/internal/domain"

type UserCreateModel struct {
	Email    string
	Password string
	Role     domain.Role
}

type UserCredentials struct {
	Email    string
	Password string
}
