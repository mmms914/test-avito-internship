package dto

import (
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/google/uuid"
)

type BookingFilter struct {
	UserID   *uuid.UUID
	Page     *int
	PageSize *int
}

type BookingCreateModel struct {
	SlotID               uuid.UUID
	CreateConferenceLink *bool
}

type BookingUpdateModel struct {
	ID     uuid.UUID
	Status domain.BookingStatus
}
