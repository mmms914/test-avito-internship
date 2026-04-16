//go:build integration
// +build integration

package integration_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/mmms914/test-avito-internship/internal/domain"
	"github.com/mmms914/test-avito-internship/internal/dto"
	"github.com/mmms914/test-avito-internship/internal/errs"
	"github.com/mmms914/test-avito-internship/pkg/ptr"
)

func (s *IntegrationTestSuite) setupBookingTestData() (*domain.User, *domain.Room, *domain.Schedule, *domain.Slot) {
	userID := uuid.New()
	user := domain.NewUser(domain.WithUserRestoreSpecs(&domain.UserRestoreSpecs{
		ID:           userID,
		Email:        "bookinguser@example.com",
		PasswordHash: "hashed",
		Role:         domain.UserRole,
		CreatedAt:    time.Now().UTC(),
	}))
	err := s.userRepo.Create(s.ctx, user)
	require.NoError(s.T(), err)

	roomID := uuid.New()
	room := domain.NewRoom(domain.WithRoomRestoreSpecs(&domain.RoomRestoreSpecs{
		ID:        roomID,
		Name:      "Booking Test Room",
		CreatedAt: time.Now().UTC(),
	}))
	err = s.roomRepo.Create(s.ctx, room)
	require.NoError(s.T(), err)

	scheduleID := uuid.New()
	schedule := domain.NewSchedule(domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
		ID:         scheduleID,
		RoomID:     roomID,
		DaysOfWeek: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
		StartTime:  9 * time.Hour,
		EndTime:    18 * time.Hour,
	}))
	err = s.scheduleRepo.Create(s.ctx, schedule)
	require.NoError(s.T(), err)

	slotID := uuid.New()
	startTime := time.Now().UTC().Add(24 * time.Hour).Truncate(24 * time.Hour).Add(10 * time.Hour)
	slot := domain.NewSlot(domain.WithSlotRestoreSpecs(&domain.SlotRestoreSpecs{
		ID:        slotID,
		RoomID:    roomID,
		StartTime: startTime,
		EndTime:   startTime.Add(30 * time.Minute),
	}))
	err = s.slotRepo.Create(s.ctx, []*domain.Slot{slot})
	require.NoError(s.T(), err)

	return user, room, schedule, slot
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestBookingRepository_Create() {
	user, _, _, slot := s.setupBookingTestData()

	bookingID := uuid.New()
	booking := domain.NewBooking(domain.WithBookingRestoreSpecs(&domain.BookingRestoreSpecs{
		ID:             bookingID,
		SlotID:         slot.ID(),
		UserID:         user.ID(),
		Status:         domain.ActiveBookingStatus,
		ConferenceLink: ptr.To("https://meet.example.com/abc123"),
		CreatedAt:      time.Now().UTC(),
	}))

	err := s.bookingRepo.Create(s.ctx, booking)
	assert.NoError(s.T(), err)
}

func (s *IntegrationTestSuite) TestBookingRepository_GetByID() {
	user, _, _, slot := s.setupBookingTestData()

	bookingID := uuid.New()
	conferenceLink := "https://meet.example.com/test"
	booking := domain.NewBooking(domain.WithBookingRestoreSpecs(&domain.BookingRestoreSpecs{
		ID:             bookingID,
		SlotID:         slot.ID(),
		UserID:         user.ID(),
		Status:         domain.ActiveBookingStatus,
		ConferenceLink: ptr.To(conferenceLink),
		CreatedAt:      time.Now().UTC(),
	}))

	err := s.bookingRepo.Create(s.ctx, booking)
	require.NoError(s.T(), err)

	found, err := s.bookingRepo.GetByID(s.ctx, bookingID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), bookingID, found.ID())
	assert.Equal(s.T(), slot.ID(), found.SlotID())
	assert.Equal(s.T(), user.ID(), found.UserID())
	assert.Equal(s.T(), domain.ActiveBookingStatus, found.Status())
	assert.Equal(s.T(), ptr.To(conferenceLink), found.ConferenceLink())
}

func (s *IntegrationTestSuite) TestBookingRepository_GetByID_NotFound() {
	bookingID := uuid.New()
	_, err := s.bookingRepo.GetByID(s.ctx, bookingID)
	assert.ErrorIs(s.T(), err, errs.ErrBookingNotFound)
}

func (s *IntegrationTestSuite) TestBookingRepository_IsSlotAlreadyBooked() {
	user, _, _, slot := s.setupBookingTestData()

	booking := domain.NewBooking(domain.WithBookingRestoreSpecs(&domain.BookingRestoreSpecs{
		ID:             uuid.New(),
		SlotID:         slot.ID(),
		UserID:         user.ID(),
		Status:         domain.ActiveBookingStatus,
		ConferenceLink: nil,
		CreatedAt:      time.Now().UTC(),
	}))

	exists, err := s.bookingRepo.IsSlotAlreadyBooked(s.ctx, slot.ID())
	assert.NoError(s.T(), err)
	assert.False(s.T(), exists)

	err = s.bookingRepo.Create(s.ctx, booking)
	require.NoError(s.T(), err)

	exists, err = s.bookingRepo.IsSlotAlreadyBooked(s.ctx, slot.ID())
	assert.NoError(s.T(), err)
	assert.True(s.T(), exists)
}

func (s *IntegrationTestSuite) TestBookingRepository_ListActive() {
	user, _, _, slot := s.setupBookingTestData()

	booking1 := domain.NewBooking(domain.WithBookingRestoreSpecs(&domain.BookingRestoreSpecs{
		ID:             uuid.New(),
		SlotID:         slot.ID(),
		UserID:         user.ID(),
		Status:         domain.ActiveBookingStatus,
		ConferenceLink: nil,
		CreatedAt:      time.Now().UTC(),
	}))

	booking2 := domain.NewBooking(domain.WithBookingRestoreSpecs(&domain.BookingRestoreSpecs{
		ID:             uuid.New(),
		SlotID:         slot.ID(),
		UserID:         user.ID(),
		Status:         domain.CancelledBookingStatus,
		ConferenceLink: nil,
		CreatedAt:      time.Now().UTC(),
	}))

	err := s.bookingRepo.Create(s.ctx, booking1)
	require.NoError(s.T(), err)
	err = s.bookingRepo.Create(s.ctx, booking2)
	require.NoError(s.T(), err)

	filter := &dto.BookingFilter{
		Page:     ptr.To(1),
		PageSize: ptr.To(10),
	}

	bookings, err := s.bookingRepo.ListActive(s.ctx, filter)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), bookings, 1)
	assert.Equal(s.T(), booking1.ID(), bookings[0].ID())
}

func (s *IntegrationTestSuite) TestBookingRepository_ListActive_WithUserFilter() {
	user1, _, _, slot1 := s.setupBookingTestData()

	booking1 := domain.NewBooking(domain.WithBookingRestoreSpecs(&domain.BookingRestoreSpecs{
		ID:             uuid.New(),
		SlotID:         slot1.ID(),
		UserID:         user1.ID(),
		Status:         domain.ActiveBookingStatus,
		ConferenceLink: nil,
		CreatedAt:      time.Now().UTC(),
	}))

	err := s.bookingRepo.Create(s.ctx, booking1)
	require.NoError(s.T(), err)

	filter := &dto.BookingFilter{
		UserID:   ptr.To(user1.ID()),
		Page:     ptr.To(1),
		PageSize: ptr.To(10),
	}

	bookings, err := s.bookingRepo.ListActive(s.ctx, filter)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), bookings, 1)
	assert.Equal(s.T(), booking1.ID(), bookings[0].ID())
}

func (s *IntegrationTestSuite) TestBookingRepository_ListActive_WithPagination() {
	user, _, _, slot := s.setupBookingTestData()

	for i := 0; i < 5; i++ {
		booking := domain.NewBooking(domain.WithBookingRestoreSpecs(&domain.BookingRestoreSpecs{
			ID:             uuid.New(),
			SlotID:         slot.ID(),
			UserID:         user.ID(),
			Status:         domain.ActiveBookingStatus,
			ConferenceLink: nil,
			CreatedAt:      time.Now().UTC(),
		}))
		err := s.bookingRepo.Create(s.ctx, booking)
		require.NoError(s.T(), err)
	}

	filter := &dto.BookingFilter{
		Page:     ptr.To(1),
		PageSize: ptr.To(2),
	}

	bookings, err := s.bookingRepo.ListActive(s.ctx, filter)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), bookings, 2)

	filter.Page = ptr.To(2)
	bookings, err = s.bookingRepo.ListActive(s.ctx, filter)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), bookings, 2)

	filter.Page = ptr.To(3)
	bookings, err = s.bookingRepo.ListActive(s.ctx, filter)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), bookings, 1)
}

func (s *IntegrationTestSuite) TestBookingRepository_Update() {
	user, _, _, slot := s.setupBookingTestData()

	bookingID := uuid.New()
	booking := domain.NewBooking(domain.WithBookingRestoreSpecs(&domain.BookingRestoreSpecs{
		ID:             bookingID,
		SlotID:         slot.ID(),
		UserID:         user.ID(),
		Status:         domain.ActiveBookingStatus,
		ConferenceLink: nil,
		CreatedAt:      time.Now().UTC(),
	}))

	err := s.bookingRepo.Create(s.ctx, booking)
	require.NoError(s.T(), err)

	updateModel := &dto.BookingUpdateModel{
		ID:     bookingID,
		Status: domain.CancelledBookingStatus,
	}

	err = s.bookingRepo.Update(s.ctx, updateModel)
	assert.NoError(s.T(), err)

	updated, err := s.bookingRepo.GetByID(s.ctx, bookingID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), domain.CancelledBookingStatus, updated.Status())
}

func (s *IntegrationTestSuite) TestBookingRepository_Update_NotFound() {
	updateModel := &dto.BookingUpdateModel{
		ID:     uuid.New(),
		Status: domain.CancelledBookingStatus,
	}

	err := s.bookingRepo.Update(s.ctx, updateModel)
	assert.ErrorIs(s.T(), err, errs.ErrBookingNotFound)
}
