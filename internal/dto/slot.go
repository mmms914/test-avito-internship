package dto

import (
	"time"

	"github.com/google/uuid"
)

type SlotFilter struct {
	RoomID uuid.UUID
	Date   time.Time
}
