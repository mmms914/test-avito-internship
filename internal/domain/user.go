package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Role      Role
	CreatedAt time.Time
}

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

func (u *User) IsAdmin() bool {
	return u.Role == AdminRole
}
