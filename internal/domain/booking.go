package domain

import (
	"time"

	"github.com/google/uuid"

	"github.com/avito-internships/test-backend-1-mmms914/pkg/ptr"
)

type Booking struct {
	id             uuid.UUID
	slotID         uuid.UUID
	userID         uuid.UUID
	status         BookingStatus
	conferenceLink *string
	createdAt      *time.Time
}

type BookingInitSpecs struct {
	SlotID         uuid.UUID
	UserID         uuid.UUID
	Status         BookingStatus
	ConferenceLink *string
}

type BookingRestoreSpecs struct {
	ID             uuid.UUID
	SlotID         uuid.UUID
	UserID         uuid.UUID
	Status         BookingStatus
	ConferenceLink *string
	CreatedAt      *time.Time
}

type BookingStatus string

const (
	ActiveBookingStatus    BookingStatus = "active"
	CancelledBookingStatus BookingStatus = "cancelled"
)

type BookingOption func(*Booking)

func NewBooking(opts ...BookingOption) *Booking {
	b := &Booking{}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func WithBookingInitSpecs(sp *BookingInitSpecs) BookingOption {
	return func(b *Booking) {
		b.id = uuid.New()
		b.slotID = sp.SlotID
		b.userID = sp.UserID
		b.status = sp.Status
		b.createdAt = ptr.To(time.Now().UTC())

		if sp.ConferenceLink != nil {
			b.conferenceLink = sp.ConferenceLink
		}
	}
}

func WithBookingRestoreSpecs(sp *BookingRestoreSpecs) BookingOption {
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
func (b *Booking) ConferenceLink() *string {
	return b.conferenceLink
}
func (b *Booking) CreatedAt() *time.Time {
	return b.createdAt
}

func (bs BookingStatus) String() string {
	return string(bs)
}
