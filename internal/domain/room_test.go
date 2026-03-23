package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
)

func TestNewRoom_WithInitSpecs(t *testing.T) {
	name := "name"
	description := "desc"
	capacity := 3

	room := domain.NewRoom(
		domain.WithRoomInitSpecs(&domain.RoomInitSpecs{
			Name:        name,
			Description: &description,
			Capacity:    &capacity,
		}),
	)

	assert.Equal(t, name, room.Name())
	assert.Equal(t, description, *room.Description())
	assert.Equal(t, capacity, *room.Capacity())
}

func TestNewRoom_WithRestoreSpecs(t *testing.T) {
	id := uuid.New()
	name := "name"
	description := "desc"
	capacity := 3
	createdAt := time.Now().UTC()

	room := domain.NewRoom(
		domain.WithRoomRestoreSpecs(&domain.RoomRestoreSpecs{
			ID:          id,
			Name:        name,
			Description: &description,
			Capacity:    &capacity,
			CreatedAt:   createdAt,
		}),
	)

	assert.Equal(t, id, room.ID())
	assert.Equal(t, name, room.Name())
	assert.Equal(t, description, *room.Description())
	assert.Equal(t, capacity, *room.Capacity())
	assert.Equal(t, createdAt, room.CreatedAt())
}
