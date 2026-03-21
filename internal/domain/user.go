package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	id        uuid.UUID
	email     string
	role      Role
	createdAt time.Time
}

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

func (u *User) ID() uuid.UUID {
	return u.id
}
func (u *User) Email() string {
	return u.email
}
func (u *User) Role() Role {
	return u.role
}
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) IsAdmin() bool {
	return u.role == AdminRole
}

type UserAuth struct {
	ID   uuid.UUID
	Role Role
}

func GetCredentialsFromContext(ctx context.Context) (*UserAuth, error) {
	userID, ok := ctx.Value("userID").(uuid.UUID)
	if !ok {
		return nil, errors.New("userID not found in context")
	}

	userRole, ok := ctx.Value("role").(Role)
	if !ok {
		return nil, errors.New("role not found in context")
	}
	return &UserAuth{
		ID:   userID,
		Role: userRole,
	}, nil
}
