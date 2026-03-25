package domain

import (
	"slices"
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	id         uuid.UUID
	roomID     uuid.UUID
	daysOfWeek []time.Weekday
	startTime  time.Duration
	endTime    time.Duration
}

type ScheduleInitSpecs struct {
	RoomID     uuid.UUID
	DaysOfWeek []time.Weekday
	StartTime  time.Duration
	EndTime    time.Duration
}

type ScheduleRestoreSpecs struct {
	ID         uuid.UUID
	RoomID     uuid.UUID
	DaysOfWeek []time.Weekday
	StartTime  time.Duration
	EndTime    time.Duration
}

type ScheduleOption func(*Schedule)

func NewSchedule(opts ...ScheduleOption) *Schedule {
	s := &Schedule{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func WithScheduleInitSpecs(sp *ScheduleInitSpecs) ScheduleOption {
	return func(s *Schedule) {
		s.id = uuid.New()
		s.roomID = sp.RoomID
		s.daysOfWeek = sp.DaysOfWeek
		s.startTime = sp.StartTime
		s.endTime = sp.EndTime
	}
}

func WithScheduleRestoreSpecs(sp *ScheduleRestoreSpecs) ScheduleOption {
	return func(s *Schedule) {
		s.id = sp.ID
		s.roomID = sp.RoomID
		s.daysOfWeek = sp.DaysOfWeek
		s.startTime = sp.StartTime
		s.endTime = sp.EndTime
	}
}

func (s *Schedule) ID() uuid.UUID {
	return s.id
}
func (s *Schedule) RoomID() uuid.UUID {
	return s.roomID
}
func (s *Schedule) DaysOfWeek() []time.Weekday {
	return s.daysOfWeek
}
func (s *Schedule) StartTime() time.Duration {
	return s.startTime
}
func (s *Schedule) EndTime() time.Duration {
	return s.endTime
}

func (s *Schedule) IsAppliedForDate(date time.Time) bool {
	return slices.Contains(s.DaysOfWeek(), date.Weekday())
}
