package converter

import (
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
)

func BookingToResponse(b *domain.Booking) *models.BookingObject {
	return &models.BookingObject{
		ID:             b.ID(),
		SlotID:         b.SlotID(),
		UserID:         b.UserID(),
		Status:         b.Status().String(),
		ConferenceLink: b.ConferenceLink(),
		CreatedAt:      b.CreatedAt(),
	}
}

func RoomArrayToResponse(r []*domain.Room) []*models.RoomObject {
	res := make([]*models.RoomObject, len(r))
	for i := range r {
		res[i] = RoomToResponse(r[i])
	}

	return res
}
