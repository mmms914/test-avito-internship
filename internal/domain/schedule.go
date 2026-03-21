package domain

import (
	"time"

	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
	"github.com/google/uuid"
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
