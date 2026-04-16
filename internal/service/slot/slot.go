package slot

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/mmms914/test-avito-internship/internal/domain"
	"github.com/mmms914/test-avito-internship/internal/dto"
	"github.com/mmms914/test-avito-internship/internal/errs"
)

type Repository interface {
	GetAllAvailable(ctx context.Context, f *dto.SlotFilter) ([]*domain.Slot, error)
	Create(ctx context.Context, slot []*domain.Slot) error
	IsSlotsExist(ctx context.Context, f *dto.SlotFilter) (bool, error)
}

type ScheduleRepository interface {
	GetForRoom(ctx context.Context, roomID uuid.UUID) (*domain.Schedule, error)
}

type RoomRepository interface {
	Exists(ctx context.Context, slotID uuid.UUID) (bool, error)
}

type Service struct {
	repo         Repository
	scheduleRepo ScheduleRepository
	roomRepo     RoomRepository
}

type Config struct {
	SlotRepo     Repository
	ScheduleRepo ScheduleRepository
	RoomRepo     RoomRepository
}

func NewService(c *Config) *Service {
	return &Service{
		repo:         c.SlotRepo,
		scheduleRepo: c.ScheduleRepo,
		roomRepo:     c.RoomRepo,
	}
}

func (s *Service) GetAvailableSlots(ctx context.Context, filter *dto.SlotFilter) ([]*domain.Slot, error) {
	roomExists, err := s.roomRepo.Exists(ctx, filter.RoomID)
	if err != nil {
		return nil, fmt.Errorf("checking if room exists: %w", err)
	}

	if !roomExists {
		return nil, errs.ErrRoomNotExists
	}

	slotsExist, err := s.repo.IsSlotsExist(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("checking if slots exist: %w", err)
	}

	if !slotsExist {
		schedule, schErr := s.scheduleRepo.GetForRoom(ctx, filter.RoomID)
		if schErr != nil {
			return nil, fmt.Errorf("getting schedule: %w", schErr)
		}

		if schedule.IsAppliedForDate(filter.Date) {
			newSlots := domain.GenerateSlots(schedule, filter.Date)
			if createErr := s.repo.Create(ctx, newSlots); createErr != nil {
				return nil, fmt.Errorf("creating slots: %w", createErr)
			}
		}
	}

	slots, err := s.repo.GetAllAvailable(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("getting available slots: %w", err)
	}

	return slots, nil
}
