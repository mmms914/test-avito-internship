package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateRoomRequest struct {
	Name        *string `json:"name" validate:"required"`
	Description *string `json:"description"`
	Capacity    *string `json:"capacity"`
}

type ListRoomResponse struct {
	Rooms []*RoomResponse `json:"rooms"`
}

type RoomResponse struct {
	ID          *uuid.UUID `json:"id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Capacity    *int       `json:"capacity"`
	CreatedAt   *time.Time `json:"createdAt"`
}
