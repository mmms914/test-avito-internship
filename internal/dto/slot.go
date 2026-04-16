package dto

import (
	"time"

	"github.com/google/uuid"
)

type SlotFilter struct {
	RoomID uuid.UUID
	Date   time.Time
}

func (f *SlotFilter) GetStartOfDay() time.Time {
	return time.Date(f.Date.Year(), f.Date.Month(), f.Date.Day(), 0, 0, 0, 0, time.UTC)
}

func (f *SlotFilter) GetEndOfDay() time.Time {
	return f.GetStartOfDay().AddDate(0, 0, 1)
}
