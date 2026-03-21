package domain

import (
	"github.com/google/uuid"
	"time"
)

type Room struct {
	ID          uuid.UUID
	Name        string
	Description string
	Capacity    int
	CreatedAt   time.Time
}
