package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/pkg/ptr"
)

func TestNewBooking_WithInitSpecs(t *testing.T) {
	slotID := uuid.New()
	userID := uuid.New()
	status := domain.ActiveBookingStatus
	conferenceLink := ptr.To("link")

	booking := domain.NewBooking(
		domain.WithBookingInitSpecs(&domain.BookingInitSpecs{
			SlotID:         slotID,
			UserID:         userID,
			Status:         status,
			ConferenceLink: conferenceLink,
		}),
	)

	assert.Equal(t, slotID, booking.SlotID())
	assert.Equal(t, userID, booking.UserID())
	assert.Equal(t, status, booking.Status())
	assert.Equal(t, conferenceLink, booking.ConferenceLink())
}

func TestNewBooking_WithRestoreSpecs(t *testing.T) {
	id := uuid.New()
	slotID := uuid.New()
	userID := uuid.New()
	status := domain.ActiveBookingStatus
	conferenceLink := ptr.To("link")
	createdAt := time.Now().UTC()

	booking := domain.NewBooking(
		domain.WithBookingRestoreSpecs(&domain.BookingRestoreSpecs{
			ID:             id,
			SlotID:         slotID,
			UserID:         userID,
			Status:         status,
			ConferenceLink: conferenceLink,
			CreatedAt:      createdAt,
		}),
	)

	assert.Equal(t, id, booking.ID())
	assert.Equal(t, slotID, booking.SlotID())
	assert.Equal(t, userID, booking.UserID())
	assert.Equal(t, status, booking.Status())
	assert.Equal(t, conferenceLink, booking.ConferenceLink())
	assert.Equal(t, createdAt, booking.CreatedAt())
}

func TestBooking_Cancel(t *testing.T) {
	id := uuid.New()
	slotID := uuid.New()
	userID := uuid.New()
	status := domain.ActiveBookingStatus
	conferenceLink := ptr.To("link")
	createdAt := time.Now().UTC()

	booking := domain.NewBooking(
		domain.WithBookingRestoreSpecs(&domain.BookingRestoreSpecs{
			ID:             id,
			SlotID:         slotID,
			UserID:         userID,
			Status:         status,
			ConferenceLink: conferenceLink,
			CreatedAt:      createdAt,
		}),
	)

	booking.Cancel()
	assert.Equal(t, domain.CancelledBookingStatus, booking.Status())
}
