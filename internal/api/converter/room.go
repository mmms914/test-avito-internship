package converter

import (
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
)

func RoomToResponse(r *domain.Room) *models.RoomResponse {
	return &models.RoomResponse{
		Room: RoomToObject(r),
	}
}

func RoomArrayToResponse(r []*domain.Room) *models.ListRoomResponse {
	res := make([]*models.RoomObject, len(r))
	for i := range r {
		res[i] = RoomToObject(r[i])
	}

	return &models.ListRoomResponse{
		Rooms: res,
	}
}

func RoomToObject(r *domain.Room) *models.RoomObject {
	return &models.RoomObject{
		ID:          r.ID(),
		Name:        r.Name(),
		Description: r.Description(),
		Capacity:    r.Capacity(),
		CreatedAt:   r.CreatedAt(),
	}
}
