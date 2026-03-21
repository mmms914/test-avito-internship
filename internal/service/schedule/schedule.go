package scheduler

import (
	"context"
	"fmt"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
	"github.com/google/uuid"
)

type RoomRepository interface {
	Exists(ctx context.Context, roomID uuid.UUID) (bool, error)
}

type Repository interface {
	Create(ctx context.Context, schedule *domain.Schedule) (*domain.Schedule, error)
}

type Service struct {
	roomRepo RoomRepository
	repo     Repository
}

type Config struct {
	roomRepo RoomRepository
	repo     Repository
}

func NewService(c *Config) *Service {
	return &Service{
		roomRepo: c.roomRepo,
		repo:     c.repo,
	}
}
func (s *Service) Create(ctx context.Context, schedule *domain.Schedule) (*domain.Schedule, error) {
	creds, err := domain.GetCredentialsFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	if creds.Role != domain.AdminRole {
		return nil, errs.ErrForbidden
	}

	roomExists, err := s.roomRepo.Exists(ctx, schedule.RoomID())
	if err != nil {
		return nil, fmt.Errorf("checking if room exists: %w", err)
	}

	if !roomExists {
		return nil, errs.ErrRoomNotExists
	}

	createdSchedule, err := s.repo.Create(ctx, schedule)
	if err != nil {
		return nil, fmt.Errorf("creating schedule: %w", err)
	}

	return createdSchedule, nil
}
