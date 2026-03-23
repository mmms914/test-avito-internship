package converter

import (
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
)

func RoomToResponse(r *domain.Room) *models.RoomObject {
	return &models.RoomObject{
		ID:          r.ID(),
		Name:        r.Name(),
		Description: r.Description(),
		Capacity:    r.Capacity(),
		CreatedAt:   r.CreatedAt(),
	}
}

func RoomArrayToResponse(r []*domain.Room) []*models.RoomObject {
	res := make([]*models.RoomObject, len(r))
	for i := range r {
		res[i] = RoomToResponse(r[i])
	}

	return res
}
