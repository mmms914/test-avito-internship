package dto

import "github.com/google/uuid"

type BookingFilter struct {
	UserID   *uuid.UUID
	Page     *int
	PageSize *int
}
