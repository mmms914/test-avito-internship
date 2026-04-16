package converter_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/mmms914/test-avito-internship/internal/api/converter"
	"github.com/mmms914/test-avito-internship/internal/domain"
	"github.com/mmms914/test-avito-internship/pkg/ptr"
)

func TestRoomToObject(t *testing.T) {
	roomID := uuid.New()
	createdAt := time.Now().UTC()

	room := domain.NewRoom(domain.WithRoomRestoreSpecs(&domain.RoomRestoreSpecs{
		ID:          roomID,
		Name:        "Conference Room",
		Description: ptr.To("Spacious room with projector"),
		Capacity:    ptr.To(10),
		CreatedAt:   createdAt,
	}))

	result := converter.RoomToObject(room)

	assert.NotNil(t, result)
	assert.Equal(t, roomID, result.ID)
	assert.Equal(t, "Conference Room", result.Name)
	assert.Equal(t, ptr.To("Spacious room with projector"), result.Description)
	assert.Equal(t, ptr.To(10), result.Capacity)
	assert.Equal(t, createdAt, result.CreatedAt)
}

func TestRoomToObjectWithoutOptionalFields(t *testing.T) {
	roomID := uuid.New()
	createdAt := time.Now().UTC()

	room := domain.NewRoom(domain.WithRoomRestoreSpecs(&domain.RoomRestoreSpecs{
		ID:        roomID,
		Name:      "Small Room",
		CreatedAt: createdAt,
	}))

	result := converter.RoomToObject(room)

	assert.NotNil(t, result)
	assert.Equal(t, roomID, result.ID)
	assert.Equal(t, "Small Room", result.Name)
	assert.Nil(t, result.Description)
	assert.Nil(t, result.Capacity)
	assert.Equal(t, createdAt, result.CreatedAt)
}

func TestRoomToResponse(t *testing.T) {
	roomID := uuid.New()
	createdAt := time.Now().UTC()

	room := domain.NewRoom(domain.WithRoomRestoreSpecs(&domain.RoomRestoreSpecs{
		ID:          roomID,
		Name:        "Conference Room",
		Description: ptr.To("Spacious room"),
		Capacity:    ptr.To(10),
		CreatedAt:   createdAt,
	}))

	result := converter.RoomToResponse(room)

	assert.NotNil(t, result)
	assert.NotNil(t, result.Room)
	assert.Equal(t, roomID, result.Room.ID)
	assert.Equal(t, "Conference Room", result.Room.Name)
	assert.Equal(t, ptr.To("Spacious room"), result.Room.Description)
	assert.Equal(t, ptr.To(10), result.Room.Capacity)
	assert.Equal(t, createdAt, result.Room.CreatedAt)
}

func TestRoomArrayToResponse(t *testing.T) {
	tests := map[string]struct {
		rooms    []*domain.Room
		expected int
	}{
		"multiple rooms": {
			rooms: []*domain.Room{
				domain.NewRoom(domain.WithRoomRestoreSpecs(&domain.RoomRestoreSpecs{
					ID:        uuid.New(),
					Name:      "Room 1",
					CreatedAt: time.Now().UTC(),
				})),
				domain.NewRoom(domain.WithRoomRestoreSpecs(&domain.RoomRestoreSpecs{
					ID:        uuid.New(),
					Name:      "Room 2",
					CreatedAt: time.Now().UTC(),
				})),
			},
			expected: 2,
		},
		"single room": {
			rooms: []*domain.Room{
				domain.NewRoom(domain.WithRoomRestoreSpecs(&domain.RoomRestoreSpecs{
					ID:        uuid.New(),
					Name:      "Single Room",
					CreatedAt: time.Now().UTC(),
				})),
			},
			expected: 1,
		},
		"empty array": {
			rooms:    []*domain.Room{},
			expected: 0,
		},
		"nil array": {
			rooms:    nil,
			expected: 0,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := converter.RoomArrayToResponse(tt.rooms)

			assert.NotNil(t, result)
			assert.Len(t, result.Rooms, tt.expected)

			for i, room := range tt.rooms {
				if i < len(result.Rooms) {
					assert.Equal(t, room.ID(), result.Rooms[i].ID)
					assert.Equal(t, room.Name(), result.Rooms[i].Name)
					assert.Equal(t, room.Description(), result.Rooms[i].Description)
					assert.Equal(t, room.Capacity(), result.Rooms[i].Capacity)
					assert.Equal(t, room.CreatedAt(), result.Rooms[i].CreatedAt)
				}
			}
		})
	}
}
