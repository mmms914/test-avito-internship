package converter_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/avito-internships/test-backend-1-mmms914/internal/api/converter"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/pkg/ptr"
)

func TestBookingToObject(t *testing.T) {
	bookingID := uuid.New()
	slotID := uuid.New()
	userID := uuid.New()
	createdAt := time.Now().UTC()
	conferenceLink := "https://meet.example.com/abc123"

	booking := domain.NewBooking(domain.WithBookingRestoreSpecs(&domain.BookingRestoreSpecs{
		ID:             bookingID,
		SlotID:         slotID,
		UserID:         userID,
		Status:         domain.ActiveBookingStatus,
		ConferenceLink: ptr.To(conferenceLink),
		CreatedAt:      createdAt,
	}))

	result := converter.BookingToObject(booking)

	assert.NotNil(t, result)
	assert.Equal(t, bookingID, result.ID)
	assert.Equal(t, slotID, result.SlotID)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, "active", result.Status)
	assert.Equal(t, ptr.To(conferenceLink), result.ConferenceLink)
	assert.Equal(t, createdAt, result.CreatedAt)
}

func TestBookingToObjectCancelled(t *testing.T) {
	bookingID := uuid.New()
	slotID := uuid.New()
	userID := uuid.New()
	createdAt := time.Now().UTC()

	booking := domain.NewBooking(domain.WithBookingRestoreSpecs(&domain.BookingRestoreSpecs{
		ID:             bookingID,
		SlotID:         slotID,
		UserID:         userID,
		Status:         domain.CancelledBookingStatus,
		ConferenceLink: nil,
		CreatedAt:      createdAt,
	}))

	result := converter.BookingToObject(booking)

	assert.NotNil(t, result)
	assert.Equal(t, bookingID, result.ID)
	assert.Equal(t, slotID, result.SlotID)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, "cancelled", result.Status)
	assert.Nil(t, result.ConferenceLink)
	assert.Equal(t, createdAt, result.CreatedAt)
}

func TestBookingToResponse(t *testing.T) {
	bookingID := uuid.New()
	slotID := uuid.New()
	userID := uuid.New()
	createdAt := time.Now().UTC()

	booking := domain.NewBooking(domain.WithBookingRestoreSpecs(&domain.BookingRestoreSpecs{
		ID:             bookingID,
		SlotID:         slotID,
		UserID:         userID,
		Status:         domain.ActiveBookingStatus,
		ConferenceLink: nil,
		CreatedAt:      createdAt,
	}))

	result := converter.BookingToResponse(booking)

	assert.NotNil(t, result)
	assert.NotNil(t, result.Booking)
	assert.Equal(t, bookingID, result.Booking.ID)
	assert.Equal(t, slotID, result.Booking.SlotID)
	assert.Equal(t, userID, result.Booking.UserID)
	assert.Equal(t, "active", result.Booking.Status)
	assert.Equal(t, createdAt, result.Booking.CreatedAt)
}

func TestBookingArrayToObjects(t *testing.T) {
	bookings := []*domain.Booking{
		domain.NewBooking(domain.WithBookingRestoreSpecs(&domain.BookingRestoreSpecs{
			ID:             uuid.New(),
			SlotID:         uuid.New(),
			UserID:         uuid.New(),
			Status:         domain.ActiveBookingStatus,
			ConferenceLink: nil,
			CreatedAt:      time.Now().UTC(),
		})),
		domain.NewBooking(domain.WithBookingRestoreSpecs(&domain.BookingRestoreSpecs{
			ID:             uuid.New(),
			SlotID:         uuid.New(),
			UserID:         uuid.New(),
			Status:         domain.CancelledBookingStatus,
			ConferenceLink: nil,
			CreatedAt:      time.Now().UTC(),
		})),
	}

	result := converter.BookingArrayToObjects(bookings)

	assert.Len(t, result, 2)
	assert.Equal(t, bookings[0].ID(), result[0].ID)
	assert.Equal(t, bookings[0].SlotID(), result[0].SlotID)
	assert.Equal(t, bookings[0].UserID(), result[0].UserID)
	assert.Equal(t, "active", result[0].Status)
	assert.Equal(t, bookings[1].ID(), result[1].ID)
	assert.Equal(t, "cancelled", result[1].Status)
}

func TestBookingArrayToObjectsEmpty(t *testing.T) {
	result := converter.BookingArrayToObjects([]*domain.Booking{})

	assert.Empty(t, result)
}

func TestBookingArrayToObjectsNil(t *testing.T) {
	result := converter.BookingArrayToObjects(nil)

	assert.Empty(t, result)
}

func TestBookingArrayToListResponse(t *testing.T) {
	bookings := []*domain.Booking{
		domain.NewBooking(domain.WithBookingRestoreSpecs(&domain.BookingRestoreSpecs{
			ID:             uuid.New(),
			SlotID:         uuid.New(),
			UserID:         uuid.New(),
			Status:         domain.ActiveBookingStatus,
			ConferenceLink: nil,
			CreatedAt:      time.Now().UTC(),
		})),
	}

	result := converter.BookingArrayToListResponse(bookings)

	assert.NotNil(t, result)
	assert.Len(t, result.Bookings, 1)
	assert.Equal(t, bookings[0].ID(), result.Bookings[0].ID)
}

func TestBookingArrayToListResponseWithPag(t *testing.T) {
	bookings := []*domain.Booking{
		domain.NewBooking(domain.WithBookingRestoreSpecs(&domain.BookingRestoreSpecs{
			ID:             uuid.New(),
			SlotID:         uuid.New(),
			UserID:         uuid.New(),
			Status:         domain.ActiveBookingStatus,
			ConferenceLink: nil,
			CreatedAt:      time.Now().UTC(),
		})),
	}

	pagination := &models.Pagination{
		Page:     1,
		PageSize: 10,
		Total:    1,
	}

	result := converter.BookingArrayToListResponseWithPag(bookings, pagination)

	assert.NotNil(t, result)
	assert.Len(t, result.Bookings, 1)
	assert.Equal(t, pagination, result.Pagination)
	assert.Equal(t, bookings[0].ID(), result.Bookings[0].ID)
}
