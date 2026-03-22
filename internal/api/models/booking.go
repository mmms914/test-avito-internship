package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateBookingRequest struct {
	SlotID               *uuid.UUID `json:"slotId" validate:"required"`
	CreateConferenceLink *bool      `json:"createConferenceLink"`
}

type CreateBookingResponse struct {
	Booking *BookingResponse `json:"booking"`
}
type ListBookingResponse struct {
	Bookings []*BookingResponse `json:"bookings"`
}

type ListBookingResponseWithPagination struct {
	Bookings   []*BookingResponse `json:"bookings"`
	Pagination *Pagination        `json:"pagination"`
}
type BookingResponse struct {
	ID             *uuid.UUID `json:"id"`
	SlotID         *uuid.UUID `json:"slotId"`
	UserID         *uuid.UUID `json:"userId"`
	Status         *string    `json:"status"`
	ConferenceLink *string    `json:"conferenceLink"`
	CreatedAt      *time.Time `json:"createdAt"`
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}
