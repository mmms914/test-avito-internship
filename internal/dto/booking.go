package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
)

type BookingFilter struct {
	UserID   *uuid.UUID
	Time     *time.Time
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
