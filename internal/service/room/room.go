package room

import (
	"context"
	"fmt"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
)

type Repository interface {
	GetAll(ctx context.Context) ([]*domain.Room, error)
	Create(ctx context.Context, room *domain.Room) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}

func (s *Service) GetAll(ctx context.Context) ([]*domain.Room, error) {
	rooms, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting all rooms: %w", err)
	}

	return rooms, nil
}

func (s *Service) Create(ctx context.Context, specs *domain.RoomInitSpecs) (*domain.Room, error) {
	creds, err := domain.GetCredentialsFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	if creds.Role != domain.AdminRole {
		return nil, errs.ErrForbidden
	}

	room := domain.NewRoom(domain.WithRoomInitSpecs(specs))

	if err = s.repo.Create(ctx, room); err != nil {
		return nil, fmt.Errorf("creating room: %w", err)
	}

	return room, nil
}
