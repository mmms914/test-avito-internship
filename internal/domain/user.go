package domain

import (
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
