package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
)

func TestNewSchedule_WithInitSpecs(t *testing.T) {
	roomID := uuid.New()
	days := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday}
	startTime := time.Hour
	endTime := 23 * time.Hour

	schedule := domain.NewSchedule(
		domain.WithScheduleInitSpecs(&domain.ScheduleInitSpecs{
			RoomID:     roomID,
			DaysOfWeek: days,
			StartTime:  startTime,
			EndTime:    endTime,
		}),
	)

	assert.Equal(t, roomID, schedule.RoomID())
	assert.Equal(t, days, schedule.DaysOfWeek())
	assert.Equal(t, startTime, schedule.StartTime())
	assert.Equal(t, endTime, schedule.EndTime())
}

func TestNewSchedule_WithRestoreSpecs(t *testing.T) {
	id := uuid.New()
	roomID := uuid.New()
	days := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday}
	startTime := time.Hour
	endTime := 23 * time.Hour

	schedule := domain.NewSchedule(
		domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
			ID:         id,
			RoomID:     roomID,
			DaysOfWeek: days,
			StartTime:  startTime,
			EndTime:    endTime,
		}),
	)

	assert.Equal(t, id, schedule.ID())
	assert.Equal(t, roomID, schedule.RoomID())
	assert.Equal(t, days, schedule.DaysOfWeek())
	assert.Equal(t, startTime, schedule.StartTime())
	assert.Equal(t, endTime, schedule.EndTime())
}
