package models

import (
	"time"

	"github.com/google/uuid"
)

type SlotResponse struct {
	ID        *uuid.UUID `json:"id"`
	RoomID    *uuid.UUID `json:"roomId"`
	StartTime *time.Time `json:"startTime"`
	EndTime   *time.Time `json:"endTime"`
}
