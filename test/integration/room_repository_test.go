//go:build integration
// +build integration

package integration_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/pkg/ptr"
)

func (s *IntegrationTestSuite) TestRoomRepository_Create() {
	roomID := uuid.New()
	room := domain.NewRoom(domain.WithRoomRestoreSpecs(&domain.RoomRestoreSpecs{
		ID:          roomID,
		Name:        "Conference Room",
		Description: ptr.To("Spacious room with projector"),
		Capacity:    ptr.To(10),
		CreatedAt:   time.Now().UTC(),
	}))

	err := s.roomRepo.Create(s.ctx, room)
	assert.NoError(s.T(), err)
}

func (s *IntegrationTestSuite) TestRoomRepository_CreateWithoutOptionalFields() {
	roomID := uuid.New()
	room := domain.NewRoom(domain.WithRoomRestoreSpecs(&domain.RoomRestoreSpecs{
		ID:        roomID,
		Name:      "Small Meeting Room",
		CreatedAt: time.Now().UTC(),
	}))

	err := s.roomRepo.Create(s.ctx, room)
	assert.NoError(s.T(), err)
}

func (s *IntegrationTestSuite) TestRoomRepository_GetAll() {
	room1 := domain.NewRoom(domain.WithRoomRestoreSpecs(&domain.RoomRestoreSpecs{
		ID:        uuid.New(),
		Name:      "Room 1",
		CreatedAt: time.Now().UTC(),
	}))
	room2 := domain.NewRoom(domain.WithRoomRestoreSpecs(&domain.RoomRestoreSpecs{
		ID:        uuid.New(),
		Name:      "Room 2",
		CreatedAt: time.Now().UTC().Add(-1 * time.Hour),
	}))
	room3 := domain.NewRoom(domain.WithRoomRestoreSpecs(&domain.RoomRestoreSpecs{
		ID:        uuid.New(),
		Name:      "Room 3",
		CreatedAt: time.Now().UTC().Add(-2 * time.Hour),
	}))

	err := s.roomRepo.Create(s.ctx, room1)
	require.NoError(s.T(), err)
	err = s.roomRepo.Create(s.ctx, room2)
	require.NoError(s.T(), err)
	err = s.roomRepo.Create(s.ctx, room3)
	require.NoError(s.T(), err)

	rooms, err := s.roomRepo.GetAll(s.ctx)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), rooms, 3)

	assert.Equal(s.T(), room1.ID(), rooms[0].ID())
	assert.Equal(s.T(), room2.ID(), rooms[1].ID())
	assert.Equal(s.T(), room3.ID(), rooms[2].ID())
}

func (s *IntegrationTestSuite) TestRoomRepository_GetAllEmpty() {
	rooms, err := s.roomRepo.GetAll(s.ctx)
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), rooms)
}

func (s *IntegrationTestSuite) TestRoomRepository_Exists() {
	roomID := uuid.New()
	room := domain.NewRoom(domain.WithRoomRestoreSpecs(&domain.RoomRestoreSpecs{
		ID:        roomID,
		Name:      "Test Room",
		CreatedAt: time.Now().UTC(),
	}))

	exists, err := s.roomRepo.Exists(s.ctx, roomID)
	assert.NoError(s.T(), err)
	assert.False(s.T(), exists)

	err = s.roomRepo.Create(s.ctx, room)
	require.NoError(s.T(), err)

	exists, err = s.roomRepo.Exists(s.ctx, roomID)
	assert.NoError(s.T(), err)
	assert.True(s.T(), exists)
}

func (s *IntegrationTestSuite) TestRoomRepository_Exists_NotFound() {
	roomID := uuid.New()
	exists, err := s.roomRepo.Exists(s.ctx, roomID)
	assert.NoError(s.T(), err)
	assert.False(s.T(), exists)
}
