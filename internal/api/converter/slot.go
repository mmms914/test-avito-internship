package converter

import (
	"github.com/mmms914/test-avito-internship/internal/api/models"
	"github.com/mmms914/test-avito-internship/internal/domain"
)

func SlotToObject(s *domain.Slot) *models.SlotObject {
	return &models.SlotObject{
		ID:        s.ID(),
		RoomID:    s.RoomID(),
		StartTime: s.StartTime(),
		EndTime:   s.EndTime(),
	}
}

func SlotArrayToResponse(s []*domain.Slot) *models.SlotsResponse {
	res := make([]*models.SlotObject, len(s))
	for i := range s {
		res[i] = SlotToObject(s[i])
	}

	return &models.SlotsResponse{
		Slots: res,
	}
}
