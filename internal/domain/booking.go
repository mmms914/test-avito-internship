package domain

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	id             uuid.UUID
	slotID         uuid.UUID
	userID         uuid.UUID
	status         BookingStatus
	conferenceLink *string
	createdAt      *time.Time
}

type BookingRestoreSpecs struct {
	ID             uuid.UUID
	SlotID         uuid.UUID
	UserID         uuid.UUID
	Status         BookingStatus
	ConferenceLink *string
	CreatedAt      *time.Time
}

type Option func(*Booking)

func NewBooking(opts ...Option) *Booking {
	b := &Booking{}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func WithBookingRestoreSpecs(sp *BookingRestoreSpecs) Option {
	return func(b *Booking) {
		b.id = sp.ID
		b.slotID = sp.SlotID
		b.userID = sp.UserID
		b.status = sp.Status

		if sp.ConferenceLink != nil {
			b.conferenceLink = sp.ConferenceLink
		}

		if sp.CreatedAt != nil {
			b.createdAt = sp.CreatedAt
		}
	}
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

func (b *Booking) ConferenceLink() *string {
	return b.conferenceLink
}
func (b *Booking) CreatedAt() *time.Time {
	return b.createdAt
}
