package domain

import (
	"time"

	"github.com/google/uuid"
)

type Slot struct {
	id        uuid.UUID
	roomID    uuid.UUID
	startTime time.Time
	endTime   time.Time
}

func (s *Slot) ID() uuid.UUID {
	return s.id
}
func (s *Slot) RoomID() uuid.UUID {
	return s.roomID
}
func (s *Slot) StartTime() time.Time {
	return s.startTime
}
func (s *Slot) EndTime() time.Time {
	return s.endTime
}
