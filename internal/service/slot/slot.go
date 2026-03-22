package slot

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/dto"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
)

type Repository interface {
	GetAllAvailable(ctx context.Context, f *dto.SlotFilter) ([]*domain.Slot, error)
	Create(ctx context.Context, slot []*domain.Slot) error
	IsSlotsExistForDate(ctx context.Context, date time.Time) (bool, error)
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
	Repo         Repository
	ScheduleRepo ScheduleRepository
	RoomRepo     RoomRepository
}

func NewService(c *Config) *Service {
	return &Service{
		repo:         c.Repo,
		scheduleRepo: c.ScheduleRepo,
		roomRepo:     c.RoomRepo,
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

	slotsExist, err := s.repo.IsSlotsExistForDate(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("checking if slots exist: %w", err)
	}

	if !slotsExist {
		schedule, creationErr := s.scheduleRepo.GetForRoom(ctx, roomID)
		if creationErr != nil {
			return nil, fmt.Errorf("getting schedule: %w", creationErr)
		}

		if slices.Contains(schedule.DaysOfWeek(), date.Weekday()) {
			newSlots := domain.GenerateSlots(schedule, date)
			if createErr := s.repo.Create(ctx, newSlots); createErr != nil {
				return nil, fmt.Errorf("creating slots: %w", createErr)
			}
		}
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
