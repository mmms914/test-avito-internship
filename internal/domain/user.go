package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
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

type UserIDKeyType string
type UserRoleKeyType string

const (
	UserIDKey   UserIDKeyType   = "userID"
	UserRoleKey UserRoleKeyType = "userRole"
)

func GetCredentialsFromContext(ctx context.Context) (*UserAuth, error) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("userID not found in context - %w", errs.ErrUnauthorized)
	}

	userRole, ok := ctx.Value(UserRoleKey).(Role)
	if !ok {
		return nil, fmt.Errorf("role not found in context - %w", errs.ErrUnauthorized)
	}
	return &UserAuth{
		ID:   userID,
		Role: userRole,
	}, nil
}
