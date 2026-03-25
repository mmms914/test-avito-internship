package converter

import (
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
)

func BookingToResponse(b *domain.Booking) *models.BookingResponse {
	return &models.BookingResponse{
		Booking: BookingToObject(b),
	}
}

func BookingToObject(b *domain.Booking) *models.BookingObject {
	return &models.BookingObject{
		ID:             b.ID(),
		SlotID:         b.SlotID(),
		UserID:         b.UserID(),
		Status:         b.Status().String(),
		ConferenceLink: b.ConferenceLink(),
		CreatedAt:      b.CreatedAt(),
	}
}

func BookingArrayToListResponse(b []*domain.Booking) *models.ListBookingResponse {
	return &models.ListBookingResponse{
		Bookings: BookingArrayToObjects(b),
	}
}

func BookingArrayToListResponseWithPag(b []*domain.Booking, p *models.Pagination) *models.ListBookingResponseWithPag {
	return &models.ListBookingResponseWithPag{
		Bookings:   BookingArrayToObjects(b),
		Pagination: p,
	}
}

func BookingArrayToObjects(b []*domain.Booking) []*models.BookingObject {
	res := make([]*models.BookingObject, len(b))
	for i := range b {
		res[i] = BookingToObject(b[i])
	}

	return res
}
