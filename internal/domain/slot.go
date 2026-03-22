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

func NewSlot(opts ...SlotOption) *Slot {
	s := &Slot{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func WithSlotRestoreSpecs(sp *SlotRestoreSpecs) SlotOption {
	return func(s *Slot) {
		s.id = sp.ID
		s.roomID = sp.RoomID
		s.startTime = sp.StartTime
		s.endTime = sp.EndTime
	}
}

type SlotRestoreSpecs struct {
	ID        uuid.UUID
	RoomID    uuid.UUID
	StartTime time.Time
	EndTime   time.Time
}

type SlotOption func(*Slot)

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

func (s *Slot) IsInPast() bool {
	return s.startTime.Before(time.Now())
}
