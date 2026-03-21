package domain

import (
	"net/url"
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	id             uuid.UUID
	slotID         uuid.UUID
	userID         uuid.UUID
	status         BookingStatus
	conferenceLink url.URL
	createdAt      time.Time
}

func (b *Booking) ID() uuid.UUID {
	return b.id
}
func (b *Booking) SlotID() uuid.UUID {
	return b.slotID
}
func (b *Booking) UserID() uuid.UUID {
	return b.userID
}

func (b *Booking) Status() BookingStatus {
	return b.status
}

type BookingStatus string

const (
	ActiveBookingStatus    BookingStatus = "active"
	CancelledBookingStatus BookingStatus = "cancelled"
)

func (b *Booking) ConferenceLink() url.URL {
	return b.conferenceLink
}
func (b *Booking) CreatedAt() time.Time {
	return b.createdAt
}
