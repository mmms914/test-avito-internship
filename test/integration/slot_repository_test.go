//go:build integration
// +build integration

package integration_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mmms914/test-avito-internship/internal/domain"
	"github.com/mmms914/test-avito-internship/internal/dto"
	"github.com/mmms914/test-avito-internship/internal/errs"
)

func (s *IntegrationTestSuite) TestSlotRepository_Create() {
	room := s.createTestRoom()
	s.createTestSchedule(room.ID())

	startTime := time.Now().UTC().Add(24 * time.Hour).Truncate(24 * time.Hour).Add(10 * time.Hour)

	slots := []*domain.Slot{
		domain.NewSlot(domain.WithSlotRestoreSpecs(&domain.SlotRestoreSpecs{
			ID:        uuid.New(),
			RoomID:    room.ID(),
			StartTime: startTime,
			EndTime:   startTime.Add(30 * time.Minute),
		})),
		domain.NewSlot(domain.WithSlotRestoreSpecs(&domain.SlotRestoreSpecs{
			ID:        uuid.New(),
			RoomID:    room.ID(),
			StartTime: startTime.Add(30 * time.Minute),
			EndTime:   startTime.Add(1 * time.Hour),
		})),
	}

	err := s.slotRepo.Create(s.ctx, slots)
	assert.NoError(s.T(), err)
}

func (s *IntegrationTestSuite) TestSlotRepository_CreateEmpty() {
	err := s.slotRepo.Create(s.ctx, []*domain.Slot{})
	assert.NoError(s.T(), err)
}

func (s *IntegrationTestSuite) TestSlotRepository_GetByID() {
	room := s.createTestRoom()
	s.createTestSchedule(room.ID())

	startTime := time.Now().UTC().Add(24 * time.Hour).Truncate(24 * time.Hour).Add(11 * time.Hour)
	slot := domain.NewSlot(domain.WithSlotRestoreSpecs(&domain.SlotRestoreSpecs{
		ID:        uuid.New(),
		RoomID:    room.ID(),
		StartTime: startTime,
		EndTime:   startTime.Add(30 * time.Minute),
	}))

	err := s.slotRepo.Create(s.ctx, []*domain.Slot{slot})
	require.NoError(s.T(), err)

	found, err := s.slotRepo.GetByID(s.ctx, slot.ID())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), slot.ID(), found.ID())
	assert.Equal(s.T(), room.ID(), found.RoomID())
	assert.Equal(s.T(), startTime, found.StartTime())
	assert.Equal(s.T(), startTime.Add(30*time.Minute), found.EndTime())
}

func (s *IntegrationTestSuite) TestSlotRepository_GetByID_NotFound() {
	slotID := uuid.New()
	_, err := s.slotRepo.GetByID(s.ctx, slotID)
	assert.ErrorIs(s.T(), err, errs.ErrSlotNotFound)
}

func (s *IntegrationTestSuite) TestSlotRepository_GetAllAvailable() {
	room := s.createTestRoom()
	s.createTestSchedule(room.ID())

	date := time.Now().UTC().Add(24 * time.Hour).Truncate(24 * time.Hour)
	startTime1 := date.Add(9 * time.Hour)
	startTime2 := date.Add(10 * time.Hour)
	startTime3 := date.Add(11 * time.Hour)

	slot1 := domain.NewSlot(domain.WithSlotRestoreSpecs(&domain.SlotRestoreSpecs{
		ID:        uuid.New(),
		RoomID:    room.ID(),
		StartTime: startTime1,
		EndTime:   startTime1.Add(30 * time.Minute),
	}))
	slot2 := domain.NewSlot(domain.WithSlotRestoreSpecs(&domain.SlotRestoreSpecs{
		ID:        uuid.New(),
		RoomID:    room.ID(),
		StartTime: startTime2,
		EndTime:   startTime2.Add(30 * time.Minute),
	}))
	slot3 := domain.NewSlot(domain.WithSlotRestoreSpecs(&domain.SlotRestoreSpecs{
		ID:        uuid.New(),
		RoomID:    room.ID(),
		StartTime: startTime3,
		EndTime:   startTime3.Add(30 * time.Minute),
	}))

	err := s.slotRepo.Create(s.ctx, []*domain.Slot{slot1, slot2, slot3})
	require.NoError(s.T(), err)

	user := s.createTestUser()

	booking := domain.NewBooking(domain.WithBookingRestoreSpecs(&domain.BookingRestoreSpecs{
		ID:             uuid.New(),
		SlotID:         slot2.ID(),
		UserID:         user.ID(),
		Status:         domain.ActiveBookingStatus,
		ConferenceLink: nil,
		CreatedAt:      time.Now().UTC(),
	}))
	err = s.bookingRepo.Create(s.ctx, booking)
	require.NoError(s.T(), err)

	filter := &dto.SlotFilter{
		RoomID: room.ID(),
		Date:   date,
	}

	slots, err := s.slotRepo.GetAllAvailable(s.ctx, filter)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), slots, 2)

	assert.Equal(s.T(), slot1.ID(), slots[0].ID())
	assert.Equal(s.T(), slot3.ID(), slots[1].ID())
}

func (s *IntegrationTestSuite) TestSlotRepository_GetAllAvailable_NoSlots() {
	room := s.createTestRoom()
	s.createTestSchedule(room.ID())

	date := time.Now().UTC().Add(24 * time.Hour).Truncate(24 * time.Hour)

	filter := &dto.SlotFilter{
		RoomID: room.ID(),
		Date:   date,
	}

	slots, err := s.slotRepo.GetAllAvailable(s.ctx, filter)
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), slots)
}

func (s *IntegrationTestSuite) TestSlotRepository_IsSlotsExist() {
	room := s.createTestRoom()
	s.createTestSchedule(room.ID())

	date := time.Now().UTC().Add(24 * time.Hour).Truncate(24 * time.Hour)
	startTime := date.Add(9 * time.Hour)

	filter := &dto.SlotFilter{
		RoomID: room.ID(),
		Date:   date,
	}

	exists, err := s.slotRepo.IsSlotsExist(s.ctx, filter)
	assert.NoError(s.T(), err)
	assert.False(s.T(), exists)

	slot := domain.NewSlot(domain.WithSlotRestoreSpecs(&domain.SlotRestoreSpecs{
		ID:        uuid.New(),
		RoomID:    room.ID(),
		StartTime: startTime,
		EndTime:   startTime.Add(30 * time.Minute),
	}))

	err = s.slotRepo.Create(s.ctx, []*domain.Slot{slot})
	require.NoError(s.T(), err)

	exists, err = s.slotRepo.IsSlotsExist(s.ctx, filter)
	assert.NoError(s.T(), err)
	assert.True(s.T(), exists)
}

func (s *IntegrationTestSuite) TestSlotRepository_GetAllAvailable_WithDifferentDays() {
	room := s.createTestRoom()
	s.createTestSchedule(room.ID())

	today := time.Now().UTC().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	dayAfter := tomorrow.Add(24 * time.Hour)

	slotToday := domain.NewSlot(domain.WithSlotRestoreSpecs(&domain.SlotRestoreSpecs{
		ID:        uuid.New(),
		RoomID:    room.ID(),
		StartTime: today.Add(10 * time.Hour),
		EndTime:   today.Add(10*time.Hour + 30*time.Minute),
	}))
	slotTomorrow := domain.NewSlot(domain.WithSlotRestoreSpecs(&domain.SlotRestoreSpecs{
		ID:        uuid.New(),
		RoomID:    room.ID(),
		StartTime: tomorrow.Add(10 * time.Hour),
		EndTime:   tomorrow.Add(10*time.Hour + 30*time.Minute),
	}))
	slotDayAfter := domain.NewSlot(domain.WithSlotRestoreSpecs(&domain.SlotRestoreSpecs{
		ID:        uuid.New(),
		RoomID:    room.ID(),
		StartTime: dayAfter.Add(10 * time.Hour),
		EndTime:   dayAfter.Add(10*time.Hour + 30*time.Minute),
	}))

	err := s.slotRepo.Create(s.ctx, []*domain.Slot{slotToday, slotTomorrow, slotDayAfter})
	require.NoError(s.T(), err)

	filter := &dto.SlotFilter{
		RoomID: room.ID(),
		Date:   tomorrow,
	}

	slots, err := s.slotRepo.GetAllAvailable(s.ctx, filter)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), slots, 1)
	assert.Equal(s.T(), slotTomorrow.ID(), slots[0].ID())
}

func (s *IntegrationTestSuite) TestSlotRepository_GetAllAvailable_WithTimeFilter() {
	room := s.createTestRoom()
	s.createTestSchedule(room.ID())

	date := time.Now().UTC().Add(24 * time.Hour).Truncate(24 * time.Hour)

	startTimeEarly := date.Add(8 * time.Hour)
	startTimeMid := date.Add(12 * time.Hour)
	startTimeLate := date.Add(16 * time.Hour)

	slotEarly := domain.NewSlot(domain.WithSlotRestoreSpecs(&domain.SlotRestoreSpecs{
		ID:        uuid.New(),
		RoomID:    room.ID(),
		StartTime: startTimeEarly,
		EndTime:   startTimeEarly.Add(30 * time.Minute),
	}))
	slotMid := domain.NewSlot(domain.WithSlotRestoreSpecs(&domain.SlotRestoreSpecs{
		ID:        uuid.New(),
		RoomID:    room.ID(),
		StartTime: startTimeMid,
		EndTime:   startTimeMid.Add(30 * time.Minute),
	}))
	slotLate := domain.NewSlot(domain.WithSlotRestoreSpecs(&domain.SlotRestoreSpecs{
		ID:        uuid.New(),
		RoomID:    room.ID(),
		StartTime: startTimeLate,
		EndTime:   startTimeLate.Add(30 * time.Minute),
	}))

	err := s.slotRepo.Create(s.ctx, []*domain.Slot{slotEarly, slotMid, slotLate})
	require.NoError(s.T(), err)

	filter := &dto.SlotFilter{
		RoomID: room.ID(),
		Date:   date,
	}

	slots, err := s.slotRepo.GetAllAvailable(s.ctx, filter)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), slots, 3)
}

func (s *IntegrationTestSuite) createTestSchedule(roomID uuid.UUID) *domain.Schedule {
	scheduleID := uuid.New()
	schedule := domain.NewSchedule(domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
		ID:         scheduleID,
		RoomID:     roomID,
		DaysOfWeek: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
		StartTime:  9 * time.Hour,
		EndTime:    18 * time.Hour,
	}))

	err := s.scheduleRepo.Create(s.ctx, schedule)
	require.NoError(s.T(), err)

	return schedule
}

func (s *IntegrationTestSuite) createTestUser() *domain.User {
	userID := uuid.New()
	user := domain.NewUser(domain.WithUserRestoreSpecs(&domain.UserRestoreSpecs{
		ID:           userID,
		Email:        "testuser_" + userID.String()[:8] + "@example.com",
		PasswordHash: "hashedpassword",
		Role:         domain.UserRole,
		CreatedAt:    time.Now().UTC(),
	}))

	err := s.userRepo.Create(s.ctx, user)
	require.NoError(s.T(), err)

	return user
}
