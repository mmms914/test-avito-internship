package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
)

type User struct {
	id           uuid.UUID
	email        string
	passwordHash string
	role         Role
	createdAt    time.Time
}

type UserInitSpecs struct {
	Email        string
	PasswordHash string
	Role         Role
}

type UserRestoreSpecs struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
}

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

func (r Role) String() string {
	return string(r)
}

func NewUser(opts ...UserOption) *User {
	u := &User{}

	for _, opt := range opts {
		opt(u)
	}

	return u
}

func WithUserInitSpecs(sp *UserInitSpecs) UserOption {
	return func(u *User) {
		u.id = uuid.New()
		u.email = sp.Email
		u.passwordHash = sp.PasswordHash
		u.role = sp.Role
		u.createdAt = time.Now().UTC()
	}
}

func WithUserRestoreSpecs(sp *UserRestoreSpecs) UserOption {
	return func(u *User) {
		u.id = sp.ID
		u.email = sp.Email
		u.passwordHash = sp.PasswordHash
		u.role = sp.Role
		u.createdAt = sp.CreatedAt
	}
}

type UserOption func(*User)

func (u *User) ID() uuid.UUID {
	return u.id
}
func (u *User) Email() string {
	return u.email
}
func (u *User) PasswordHash() string {
	return u.passwordHash
}
func (u *User) Role() Role {
	return u.role
}
func (u *User) CreatedAt() time.Time {
	return u.createdAt
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
