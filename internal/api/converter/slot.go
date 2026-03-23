package converter

import (
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
)

func SlotToResponse(s *domain.Slot) *models.SlotObject {
	return &models.SlotObject{
		ID:        s.ID(),
		RoomID:    s.RoomID(),
		StartTime: s.StartTime(),
		EndTime:   s.EndTime(),
	}
}

func SlotArrayToResponse(s []*domain.Slot) []*models.SlotObject {
	res := make([]*models.SlotObject, len(s))
	for i := range s {
		res[i] = SlotToResponse(s[i])
	}

	return res
}
