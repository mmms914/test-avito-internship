package booking

import (
	"fmt"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/dto"
	"github.com/google/uuid"
)

type Repository interface {
	Create(booking *domain.Booking) (*domain.Booking, error)
	List(filter *dto.BookingFilter) (*domain.Booking, error)
	Cancel(bookingID uuid.UUID) (*domain.Booking, error)
}

//nolint:iface // interfaces has different ways for development
type SlotRepository interface {
	Exists(slotID uuid.UUID) (bool, error)
}

//nolint:iface // interfaces has different ways for development
type UserRepository interface {
	Exists(userID uuid.UUID) (bool, error)
}

type Service struct {
	repo     Repository
	slotRepo SlotRepository
	userRepo UserRepository
}

type Config struct {
	bookingRepo Repository
	slotRepo    SlotRepository
	userRepo    UserRepository
}

func NewService(c *Config) *Service {
	return &Service{
		repo:     c.bookingRepo,
		slotRepo: c.slotRepo,
		userRepo: c.userRepo,
	}
}

func (s *Service) Create(booking *domain.Booking) (*domain.Booking, error) {
	slotExists, err := s.slotRepo.Exists(booking.SlotID())
	if err != nil {
		return nil, fmt.Errorf("checking if slot exists: %w", err)
	}

	if !slotExists {
		return nil, fmt.Errorf("slot %s does not exist", booking.SlotID())
	}

	userExists, err := s.userRepo.Exists(booking.UserID())
	if err != nil {
		return nil, fmt.Errorf("checking if user exists: %w", err)
	}

	if !userExists {
		return nil, fmt.Errorf("user %s does not exist", booking.UserID())
	}

	createdBooking, err := s.repo.Create(booking)
	if err != nil {
		return nil, fmt.Errorf("creating booking: %w", err)
	}

	return createdBooking, nil
}
