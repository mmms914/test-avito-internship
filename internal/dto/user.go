package dto

import "github.com/mmms914/test-avito-internship/internal/domain"

type UserCreateModel struct {
	Email    string
	Password string
	Role     domain.Role
}

type UserCredentials struct {
	Email    string
	Password string
}
