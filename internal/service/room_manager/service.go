package room_manager

import (
	"fmt"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
)

type RoomRepository interface {
	GetAll() ([]domain.Room, error)
	Create(room *domain.Room) (*domain.Room, error)
}

type RoomManager struct {
	repo RoomRepository
}

func NewRoomManager(repo RoomRepository) *RoomManager {
	return &RoomManager{repo}
}

func (rm *RoomManager) GetAll() ([]domain.Room, error) {
	rooms, err := rm.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("getting all rooms: %w", err)
	}

	return rooms, nil
}

func (rm *RoomManager) Create(room *domain.Room) (*domain.Room, error) {
	createdRoom, err := rm.repo.Create(room)
	if err != nil {
		return nil, fmt.Errorf("creating room: %w", err)
	}

	return createdRoom, nil
}
