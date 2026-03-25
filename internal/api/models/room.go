package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateRoomRequest struct {
	Name        *string `json:"name" validate:"required"`
	Description *string `json:"description"`
	Capacity    *int    `json:"capacity"`
}

type ListRoomResponse struct {
	Rooms []*RoomObject `json:"rooms"`
}

type RoomResponse struct {
	Room *RoomObject `json:"room"`
}

type RoomObject struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Capacity    *int      `json:"capacity,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}
