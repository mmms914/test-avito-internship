package models

import (
	"github.com/google/uuid"
	"time"
)

type SlotsResponse struct {
	Slots []*SlotObject `json:"slots"`
}

type SlotObject struct {
	ID        uuid.UUID `json:"id"`
	RoomID    uuid.UUID `json:"roomId"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}
