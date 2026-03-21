package scheduler

import (
	"fmt"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/google/uuid"
)

type RoomRepository interface {
	Exists(roomID uuid.UUID) (bool, error)
}

type Repository interface {
	Create(schedule *domain.Schedule) (*domain.Schedule, error)
}

type Service struct {
	roomRepo RoomRepository
	repo     Repository
}

type Config struct {
	roomRepo     RoomRepository
	scheduleRepo Repository
}

func NewService(c *Config) *Service {
	return &Service{
		roomRepo: c.roomRepo,
		repo:     c.scheduleRepo,
	}
}
func (s *Service) Create(schedule *domain.Schedule) (*domain.Schedule, error) {
	roomExists, err := s.roomRepo.Exists(schedule.RoomID())
	if err != nil {
		return nil, fmt.Errorf("checking if room exists: %w", err)
	}

	if !roomExists {
		return nil, fmt.Errorf("room %s does not exist", schedule.RoomID())
	}

	createdSchedule, err := s.repo.Create(schedule)
	if err != nil {
		return nil, fmt.Errorf("creating schedule: %w", err)
	}

	return createdSchedule, nil
}
