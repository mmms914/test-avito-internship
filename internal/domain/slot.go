package domain

import (
	"github.com/google/uuid"
	"time"
)

type Slot struct {
	ID        uuid.UUID
	RoomID    uuid.UUID
	StartTime time.Time
	EndTime   time.Time
}
