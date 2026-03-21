package domain

import (
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
	"github.com/google/uuid"
	"time"
)

type Schedule struct {
	ID         uuid.UUID
	RoomID     uuid.UUID
	DaysOfWeek []Day
	StartTime  time.Time
	EndTime    time.Time
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
