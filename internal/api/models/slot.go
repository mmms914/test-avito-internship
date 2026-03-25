package models

import (
	"time"

	"github.com/google/uuid"
)

type SlotsResponse struct {
	Slots []*SlotObject `json:"slots"`
}

type SlotObject struct {
	ID        uuid.UUID `json:"id"`
	RoomID    uuid.UUID `json:"roomId"`
	StartTime time.Time `json:"start"`
	EndTime   time.Time `json:"end"`
}
