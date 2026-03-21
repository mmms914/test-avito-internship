package slot

import (
	"fmt"
	"time"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/dto"
	"github.com/google/uuid"
)

type Repository interface {
	GetAllAvailable(f *dto.SlotFilter) ([]*domain.Slot, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) GetAvailableSlots(roomID uuid.UUID, date time.Time) ([]*domain.Slot, error) {
	filter := &dto.SlotFilter{
		RoomID: roomID,
		Date:   date,
	}

	slots, err := s.repo.GetAllAvailable(filter)
	if err != nil {
		return nil, fmt.Errorf("getting available slots: %w", err)
	}

	return slots, nil
}
