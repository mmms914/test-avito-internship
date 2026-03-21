package room

import (
	"fmt"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
)

type Repository interface {
	GetAll() ([]*domain.Room, error)
	Create(room *domain.Room) (*domain.Room, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) GetAll() ([]*domain.Room, error) {
	rooms, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("getting all rooms: %w", err)
	}

	return rooms, nil
}

func (s *Service) Create(room *domain.Room) (*domain.Room, error) {
	createdRoom, err := s.repo.Create(room)
	if err != nil {
		return nil, fmt.Errorf("creating room: %w", err)
	}

	return createdRoom, nil
}
