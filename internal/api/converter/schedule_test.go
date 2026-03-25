package converter_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/avito-internships/test-backend-1-mmms914/internal/api/converter"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
)

func TestScheduleToObject(t *testing.T) {
	scheduleID := uuid.New()
	roomID := uuid.New()
	daysOfWeek := []time.Weekday{time.Monday, time.Wednesday, time.Friday}
	startTime := 9 * time.Hour
	endTime := 18 * time.Hour

	schedule := domain.NewSchedule(domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
		ID:         scheduleID,
		RoomID:     roomID,
		DaysOfWeek: daysOfWeek,
		StartTime:  startTime,
		EndTime:    endTime,
	}))

	result := converter.ScheduleToObject(schedule)

	assert.NotNil(t, result)
	assert.Equal(t, scheduleID, result.ID)
	assert.Equal(t, roomID, result.RoomID)
	assert.Equal(t, []int{1, 3, 5}, result.DaysOfWeek)
	assert.Equal(t, "09:00", result.StartTime)
	assert.Equal(t, "18:00", result.EndTime)
}

func TestScheduleToResponse(t *testing.T) {
	scheduleID := uuid.New()
	roomID := uuid.New()
	daysOfWeek := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday}
	startTime := 10 * time.Hour
	endTime := 19 * time.Hour

	schedule := domain.NewSchedule(domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
		ID:         scheduleID,
		RoomID:     roomID,
		DaysOfWeek: daysOfWeek,
		StartTime:  startTime,
		EndTime:    endTime,
	}))

	result := converter.ScheduleToResponse(schedule)

	assert.NotNil(t, result)
	assert.NotNil(t, result.Schedule)
	assert.Equal(t, scheduleID, result.Schedule.ID)
	assert.Equal(t, roomID, result.Schedule.RoomID)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result.Schedule.DaysOfWeek)
	assert.Equal(t, "10:00", result.Schedule.StartTime)
	assert.Equal(t, "19:00", result.Schedule.EndTime)
}
