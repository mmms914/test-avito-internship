package domain

import (
	"net/url"
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID             uuid.UUID
	SlotID         uuid.UUID
	UserID         uuid.UUID
	Status         BookingStatus
	ConferenceLink url.URL
	CreatedAt      time.Time
}

type BookingStatus string

const (
	ActiveBookingStatus    BookingStatus = "active"
	CancelledBookingStatus BookingStatus = "cancelled"
)
