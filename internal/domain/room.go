package domain

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	id          uuid.UUID
	name        string
	description string
	capacity    int
	createdAt   time.Time
}

func (room *Room) ID() uuid.UUID {
	return room.id
}
func (room *Room) Name() string {
	return room.name
}
func (room *Room) Description() string {
	return room.description
}
func (room *Room) Capacity() int {
	return room.capacity
}
func (room *Room) CreatedAt() time.Time {
	return room.createdAt
}
