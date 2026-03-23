package schedule

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
)

type RoomRepository interface {
	Exists(ctx context.Context, roomID uuid.UUID) (bool, error)
}

type Repository interface {
	Create(ctx context.Context, schedule *domain.Schedule) (*domain.Schedule, error)
	ExistsForRoom(ctx context.Context, roomID uuid.UUID) (bool, error)
}

type Service struct {
	roomRepo RoomRepository
	repo     Repository
}

type Config struct {
	RoomRepo RoomRepository
	Repo     Repository
}

func NewService(c *Config) *Service {
	return &Service{
		roomRepo: c.RoomRepo,
		repo:     c.Repo,
	}
}
func (s *Service) Create(ctx context.Context, specs *domain.ScheduleInitSpecs) (*domain.Schedule, error) {
	creds, err := domain.GetCredentialsFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	if creds.Role != domain.AdminRole {
		return nil, errs.ErrForbidden
	}

	roomExists, err := s.roomRepo.Exists(ctx, specs.RoomID)
	if err != nil {
		return nil, fmt.Errorf("checking if room exists: %w", err)
	}

	if !roomExists {
		return nil, errs.ErrRoomNotExists
	}

	scheduleExists, err := s.repo.ExistsForRoom(ctx, specs.RoomID)
	if err != nil {
		return nil, fmt.Errorf("checking if schedule exists: %w", err)
	}

	if scheduleExists {
		return nil, errs.ErrScheduleExists
	}

	schedule := domain.NewSchedule(domain.WithScheduleInitSpecs(specs))
	createdSchedule, err := s.repo.Create(ctx, schedule)
	if err != nil {
		return nil, fmt.Errorf("creating schedule: %w", err)
	}

	return createdSchedule, nil
}
