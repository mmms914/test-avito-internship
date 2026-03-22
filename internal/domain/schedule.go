package domain

import (
	"time"

	"github.com/google/uuid"

	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
)

type Schedule struct {
	id         uuid.UUID
	roomID     uuid.UUID
	daysOfWeek []Day
	startTime  time.Time
	endTime    time.Time
}

type Day int

const (
	Monday Day = iota + 1
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func (d Day) Int() int {
	return int(d)
}

func DayFromInt(d int) (Day, error) {
	switch d {
	case Monday.Int():
		return Monday, nil
	case Tuesday.Int():
		return Tuesday, nil
	case Wednesday.Int():
		return Wednesday, nil
	case Thursday.Int():
		return Thursday, nil
	case Friday.Int():
		return Friday, nil
	case Saturday.Int():
		return Saturday, nil
	case Sunday.Int():
		return Sunday, nil
	default:
		return -1, errs.ErrDayIsInvalid
	}
}

type ScheduleRestoreSpecs struct {
	ID         uuid.UUID
	RoomID     uuid.UUID
	DaysOfWeek []Day
	StartTime  time.Time
	EndTime    time.Time
}

type ScheduleOption func(*Schedule)

func NewSchedule(opts ...ScheduleOption) *Schedule {
	s := &Schedule{}

	for _, opt := range opts {
		opt(s)
	}

	return s
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
func (s *Schedule) DaysOfWeek() []Day {
	return s.daysOfWeek
}
func (s *Schedule) StartTime() time.Time {
	return s.startTime
}
func (s *Schedule) EndTime() time.Time {
	return s.endTime
}
