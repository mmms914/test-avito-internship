//go:build integration
// +build integration

package integration_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mmms914/test-avito-internship/internal/domain"
	"github.com/mmms914/test-avito-internship/internal/errs"
)

func (s *IntegrationTestSuite) TestScheduleRepository_Create() {
	room := s.createTestRoom()

	scheduleID := uuid.New()
	daysOfWeek := []time.Weekday{time.Monday, time.Wednesday, time.Friday}
	startTime := 9 * time.Hour
	endTime := 18 * time.Hour

	schedule := domain.NewSchedule(domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
		ID:         scheduleID,
		RoomID:     room.ID(),
		DaysOfWeek: daysOfWeek,
		StartTime:  startTime,
		EndTime:    endTime,
	}))

	err := s.scheduleRepo.Create(s.ctx, schedule)
	assert.NoError(s.T(), err)
}

func (s *IntegrationTestSuite) TestScheduleRepository_GetForRoom() {
	room := s.createTestRoom()

	scheduleID := uuid.New()
	daysOfWeek := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday}
	startTime := 10 * time.Hour
	endTime := 19 * time.Hour

	schedule := domain.NewSchedule(domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
		ID:         scheduleID,
		RoomID:     room.ID(),
		DaysOfWeek: daysOfWeek,
		StartTime:  startTime,
		EndTime:    endTime,
	}))

	err := s.scheduleRepo.Create(s.ctx, schedule)
	require.NoError(s.T(), err)

	found, err := s.scheduleRepo.GetForRoom(s.ctx, room.ID())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), scheduleID, found.ID())
	assert.Equal(s.T(), room.ID(), found.RoomID())
	assert.Equal(s.T(), daysOfWeek, found.DaysOfWeek())
	assert.Equal(s.T(), startTime, found.StartTime())
	assert.Equal(s.T(), endTime, found.EndTime())
}

func (s *IntegrationTestSuite) TestScheduleRepository_GetForRoom_NotFound() {
	roomID := uuid.New()
	_, err := s.scheduleRepo.GetForRoom(s.ctx, roomID)
	assert.ErrorIs(s.T(), err, errs.ErrScheduleNotExists)
}

func (s *IntegrationTestSuite) TestScheduleRepository_ExistsForRoom() {
	room := s.createTestRoom()

	exists, err := s.scheduleRepo.ExistsForRoom(s.ctx, room.ID())
	assert.NoError(s.T(), err)
	assert.False(s.T(), exists)

	schedule := domain.NewSchedule(domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
		ID:         uuid.New(),
		RoomID:     room.ID(),
		DaysOfWeek: []time.Weekday{time.Monday},
		StartTime:  9 * time.Hour,
		EndTime:    17 * time.Hour,
	}))

	err = s.scheduleRepo.Create(s.ctx, schedule)
	require.NoError(s.T(), err)

	exists, err = s.scheduleRepo.ExistsForRoom(s.ctx, room.ID())
	assert.NoError(s.T(), err)
	assert.True(s.T(), exists)
}

func (s *IntegrationTestSuite) TestScheduleRepository_GetForRoom_WithFullWeekSchedule() {
	room := s.createTestRoom()

	daysOfWeek := []time.Weekday{
		time.Monday, time.Tuesday, time.Wednesday,
		time.Thursday, time.Friday, time.Saturday, time.Sunday,
	}
	startTime := 8 * time.Hour
	endTime := 20 * time.Hour

	schedule := domain.NewSchedule(domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
		ID:         uuid.New(),
		RoomID:     room.ID(),
		DaysOfWeek: daysOfWeek,
		StartTime:  startTime,
		EndTime:    endTime,
	}))

	err := s.scheduleRepo.Create(s.ctx, schedule)
	require.NoError(s.T(), err)

	found, err := s.scheduleRepo.GetForRoom(s.ctx, room.ID())
	assert.NoError(s.T(), err)
	assert.Len(s.T(), found.DaysOfWeek(), 7)
	assert.Equal(s.T(), startTime, found.StartTime())
	assert.Equal(s.T(), endTime, found.EndTime())
}

func (s *IntegrationTestSuite) createTestRoom() *domain.Room {
	roomID := uuid.New()
	room := domain.NewRoom(domain.WithRoomRestoreSpecs(&domain.RoomRestoreSpecs{
		ID:        roomID,
		Name:      "Test Room for Schedule",
		CreatedAt: time.Now().UTC(),
	}))

	err := s.roomRepo.Create(s.ctx, room)
	require.NoError(s.T(), err)

	return room
}
