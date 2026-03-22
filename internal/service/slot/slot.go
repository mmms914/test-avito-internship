package slot

import (
	"context"
	"fmt"
	"time"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/dto"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
	"github.com/google/uuid"
)

type Repository interface {
	GetAllAvailable(ctx context.Context, f *dto.SlotFilter) ([]*domain.Slot, error)
}

type RoomRepository interface {
	Exists(ctx context.Context, slotID uuid.UUID) (bool, error)
}

type Service struct {
	repo     Repository
	roomRepo RoomRepository
}

type Config struct {
	Repo     Repository
	RoomRepo RoomRepository
}

func NewService(c *Config) *Service {
	return &Service{
		repo:     c.Repo,
		roomRepo: c.RoomRepo,
	}
}

func (s *Service) GetAvailableSlots(ctx context.Context, roomID uuid.UUID, date time.Time) ([]*domain.Slot, error) {
	roomExists, err := s.roomRepo.Exists(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("checking if room exists: %w", err)
	}

	if !roomExists {
		return nil, errs.ErrRoomNotExists
	}

	filter := &dto.SlotFilter{
		RoomID: roomID,
		Date:   date,
	}

	slots, err := s.repo.GetAllAvailable(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("getting available slots: %w", err)
	}

	return slots, nil
}
