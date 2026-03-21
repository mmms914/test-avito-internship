package booking

import (
	"fmt"
	"time"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/dto"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
	"github.com/avito-internships/test-backend-1-mmms914/pkg/ptr"
	"github.com/google/uuid"
)

type Repository interface {
	GetByID(bookingID uuid.UUID) (*domain.Booking, error)
	Create(booking *domain.Booking) (*domain.Booking, error)
	List(filter *dto.BookingFilter) (*domain.Booking, error)
	Update(bum *dto.BookingUpdateModel) (*domain.Booking, error)
}

//nolint:iface // interfaces has different ways for development
type SlotRepository interface {
	Exists(slotID uuid.UUID) (bool, error)
}

//nolint:iface // interfaces has different ways for development
type UserRepository interface {
	Exists(userID uuid.UUID) (bool, error)
}

type LinkManager interface {
	Create() (string, error)
}

type Service struct {
	repo        Repository
	slotRepo    SlotRepository
	userRepo    UserRepository
	linkManager LinkManager
}

type Config struct {
	bookingRepo Repository
	slotRepo    SlotRepository
	userRepo    UserRepository
	linkManager LinkManager
}

func NewService(c *Config) *Service {
	return &Service{
		repo:        c.bookingRepo,
		slotRepo:    c.slotRepo,
		userRepo:    c.userRepo,
		linkManager: c.linkManager,
	}
}

func (s *Service) Create(booking *dto.BookingCreateModel, userID uuid.UUID) (*domain.Booking, error) {
	slotExists, err := s.slotRepo.Exists(booking.SlotID)
	if err != nil {
		return nil, fmt.Errorf("checking if slot exists: %w", err)
	}

	if !slotExists {
		return nil, errs.ErrSlotNotFound
	}

	userExists, err := s.userRepo.Exists(userID)
	if err != nil {
		return nil, fmt.Errorf("checking if user exists: %w", err)
	}

	if !userExists {
		return nil, errs.ErrUserNotFound
	}

	var conferenceLink *string
	if booking.CreateConferenceLink != nil && *booking.CreateConferenceLink {
		cLink, err := s.linkManager.Create()
		if err != nil {
			return nil, fmt.Errorf("creating conference link: %w", err)
		}

		conferenceLink = &cLink
	}

	createdBooking, err := s.repo.Create(domain.NewBooking(domain.WithBookingRestoreSpecs(
		&domain.BookingRestoreSpecs{
			ID:             uuid.New(),
			SlotID:         booking.SlotID,
			UserID:         userID,
			Status:         domain.ActiveBookingStatus,
			ConferenceLink: conferenceLink,
			CreatedAt:      ptr.To(time.Now().UTC()),
		})))
	if err != nil {
		// TODO: cancel conference for link
		return nil, fmt.Errorf("creating booking: %w", err)
	}

	return createdBooking, nil
}

func (s *Service) List(filter *dto.BookingFilter) (*domain.Booking, error) {
	if filter.UserID != nil {
		userExists, err := s.userRepo.Exists(*filter.UserID)
		if err != nil {
			return nil, fmt.Errorf("checking if user exists: %w", err)
		}

		if !userExists {
			return nil, errs.ErrUserNotFound
		}
	}

	bookings, err := s.repo.List(filter)
	if err != nil {
		return nil, fmt.Errorf("listing bookings: %w", err)
	}

	return bookings, nil
}

func (s *Service) ListForUser(userID uuid.UUID) (*domain.Booking, error) {
	userExists, err := s.userRepo.Exists(userID)
	if err != nil {
		return nil, fmt.Errorf("checking if user exists: %w", err)
	}

	if !userExists {
		return nil, errs.ErrUserNotFound
	}

	filter := &dto.BookingFilter{
		UserID: &userID,
	}

	bookings, err := s.repo.List(filter)
	if err != nil {
		return nil, fmt.Errorf("listing bookings: %w", err)
	}

	return bookings, nil
}

func (s *Service) Cancel(bookingID uuid.UUID, userID uuid.UUID) (*domain.Booking, error) {
	booking, err := s.repo.GetByID(bookingID)
	if err != nil {
		return nil, fmt.Errorf("getting booking: %w", err)
	}

	if booking.UserID() != userID {
		return nil, errs.ErrForbidden
	}

	model := &dto.BookingUpdateModel{
		ID:     booking.UserID(),
		Status: booking.Status(),
	}

	updatedBooking, err := s.repo.Update(model)
	if err != nil {
		return nil, fmt.Errorf("updating booking: %w", err)
	}

	return updatedBooking, nil
}
