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

type RoomInitSpecs struct {
	Name        string
	Description string
	Capacity    int
}

type RoomRestoreSpecs struct {
	Id          uuid.UUID
	Name        string
	Description string
	Capacity    int
	CreatedAt   time.Time
}

func NewRoom(opts ...RoomOption) *Room {
	r := &Room{}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func WithRoomInitSpecs(sp *RoomInitSpecs) RoomOption {
	return func(r *Room) {
		r.id = uuid.New()
		r.name = sp.Name
		r.description = sp.Description
		r.capacity = sp.Capacity
		r.createdAt = time.Now()
	}
}

func WithRoomRestoreSpecs(sp *RoomRestoreSpecs) RoomOption {
	return func(r *Room) {
		r.id = sp.Id
		r.name = sp.Name
		r.description = sp.Description
		r.capacity = sp.Capacity
		r.createdAt = sp.CreatedAt
	}
}

type RoomOption func(*Room)

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
