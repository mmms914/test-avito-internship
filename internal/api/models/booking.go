package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateBookingRequest struct {
	SlotID               *string `json:"slotId" validate:"required,uuid"`
	CreateConferenceLink *bool   `json:"createConferenceLink"`
}

type ListBookingResponse struct {
	Bookings []*BookingObject `json:"bookings"`
}

type ListBookingResponseWithPag struct {
	Bookings   []*BookingObject `json:"bookings"`
	Pagination *Pagination      `json:"pagination"`
}
type BookingResponse struct {
	Booking *BookingObject `json:"booking"`
}

type BookingObject struct {
	ID             uuid.UUID `json:"id"`
	SlotID         uuid.UUID `json:"slotId"`
	UserID         uuid.UUID `json:"userId"`
	Status         string    `json:"status"`
	ConferenceLink *string   `json:"conferenceLink,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}
