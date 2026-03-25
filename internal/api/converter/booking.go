package converter

import (
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
)

func BookingToResponse(b *domain.Booking) *models.BookingResponse {
	return &models.BookingResponse{
		Booking: BookingToResponseObject(b),
	}
}

func BookingToResponseObject(b *domain.Booking) *models.BookingObject {
	return &models.BookingObject{
		ID:             b.ID(),
		SlotID:         b.SlotID(),
		UserID:         b.UserID(),
		Status:         b.Status().String(),
		ConferenceLink: b.ConferenceLink(),
		CreatedAt:      b.CreatedAt(),
	}
}

func BookingArrayToResponse(b []*domain.Booking) []*models.BookingObject {
	res := make([]*models.BookingObject, len(b))
	for i := range b {
		res[i] = BookingToResponseObject(b[i])
	}

	return res
}
