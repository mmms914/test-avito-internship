package domain

import (
	"slices"
	"time"

	"github.com/google/uuid"
)

const (
	SlotDuration = 30 * time.Minute
	DayDuration  = 24 * time.Hour
)

type Slot struct {
	id        uuid.UUID
	roomID    uuid.UUID
	startTime time.Time
	endTime   time.Time
}

type SlotInitSpecs struct {
	RoomID    uuid.UUID
	StartTime time.Time
	EndTime   time.Time
}
type SlotRestoreSpecs struct {
	ID        uuid.UUID
	RoomID    uuid.UUID
	StartTime time.Time
	EndTime   time.Time
}

func NewSlot(opts ...SlotOption) *Slot {
	s := &Slot{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func WithSlotInitSpecs(sp *SlotInitSpecs) SlotOption {
	return func(s *Slot) {
		s.id = uuid.New()
		s.roomID = sp.RoomID
		s.startTime = sp.StartTime
		s.endTime = sp.EndTime
	}
}

func WithSlotRestoreSpecs(sp *SlotRestoreSpecs) SlotOption {
	return func(s *Slot) {
		s.id = sp.ID
		s.roomID = sp.RoomID
		s.startTime = sp.StartTime
		s.endTime = sp.EndTime
	}
}

type SlotOption func(*Slot)

func (s *Slot) ID() uuid.UUID {
	return s.id
}
func (s *Slot) RoomID() uuid.UUID {
	return s.roomID
}
func (s *Slot) StartTime() time.Time {
	return s.startTime.UTC()
}
func (s *Slot) EndTime() time.Time {
	return s.endTime.UTC()
}

func (s *Slot) IsInPast() bool {
	return s.startTime.Before(time.Now())
}

func GenerateSlots(schedule *Schedule, date time.Time) []*Slot {
	slots := make([]*Slot, 0)

	if slices.Contains(schedule.DaysOfWeek(), date.Weekday()) {
		beginTime := date.UTC().Truncate(DayDuration).Add(schedule.startTime)
		slotsPerDay := int((schedule.endTime - schedule.startTime) / SlotDuration)

		for slotNumber := range slotsPerDay {
			slots = append(slots, &Slot{
				id:        uuid.New(),
				roomID:    schedule.roomID,
				startTime: beginTime.Add(time.Duration(slotNumber) * SlotDuration),
				endTime:   beginTime.Add(time.Duration(slotNumber) * SlotDuration).Add(SlotDuration),
			})
		}
	}

	return slots
}
