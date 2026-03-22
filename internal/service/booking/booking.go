package booking

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/dto"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
	"github.com/avito-internships/test-backend-1-mmms914/pkg/ptr"
)

type Repository interface {
	GetByID(ctx context.Context, bookingID uuid.UUID) (*domain.Booking, error)
	IsSlotAlreadyBooked(ctx context.Context, slotID uuid.UUID) (bool, error)
	Create(ctx context.Context, booking *domain.Booking) (*domain.Booking, error)
	List(ctx context.Context, filter *dto.BookingFilter) ([]*domain.Booking, error)
	Update(ctx context.Context, bum *dto.BookingUpdateModel) (*domain.Booking, error)
}

//nolint:iface // interfaces has different ways for development
type SlotRepository interface {
	Exists(ctx context.Context, slotID uuid.UUID) (bool, error)
}

//nolint:iface // interfaces has different ways for development
type UserRepository interface {
	Exists(ctx context.Context, userID uuid.UUID) (bool, error)
}

type LinkManager interface {
	Create(ctx context.Context) (string, error)
	Cancel(ctx context.Context, link string) error
}

type Logger interface {
	Error(msg string)
}

type Service struct {
	logger Logger

	repo        Repository
	slotRepo    SlotRepository
	userRepo    UserRepository
	linkManager LinkManager
}

type Config struct {
	Logger Logger

	BookingRepo Repository
	SlotRepo    SlotRepository
	UserRepo    UserRepository
	LinkManager LinkManager
}

func NewService(c *Config) *Service {
	return &Service{
		logger: c.Logger,

		repo:        c.BookingRepo,
		slotRepo:    c.SlotRepo,
		userRepo:    c.UserRepo,
		linkManager: c.LinkManager,
	}
}

func (s *Service) Create(ctx context.Context, booking *dto.BookingCreateModel) (*domain.Booking, error) {
	creds, err := domain.GetCredentialsFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	if creds.Role != domain.UserRole {
		return nil, errs.ErrForbidden
	}

	slotExists, err := s.slotRepo.Exists(ctx, booking.SlotID)
	if err != nil {
		return nil, fmt.Errorf("checking if slot exists: %w", err)
	}

	if !slotExists {
		return nil, errs.ErrSlotNotFound
	}

	isAlreadyBooked, err := s.repo.IsSlotAlreadyBooked(ctx, booking.SlotID)
	if err != nil {
		return nil, fmt.Errorf("checking if slot is already booked: %w", err)
	}

	if isAlreadyBooked {
		return nil, errs.ErrSlotAlreadyBooked
	}

	userExists, err := s.userRepo.Exists(ctx, creds.ID)
	if err != nil {
		return nil, fmt.Errorf("checking if user exists: %w", err)
	}

	if !userExists {
		return nil, errs.ErrUserNotFound
	}

	var conferenceLink *string
	if booking.CreateConferenceLink != nil && *booking.CreateConferenceLink {
		cLink, linkErr := s.linkManager.Create(ctx)
		if linkErr != nil {
			return nil, fmt.Errorf("creating conference link: %w", linkErr)
		}

		conferenceLink = &cLink
	}

	createdBooking, err := s.repo.Create(ctx, domain.NewBooking(domain.WithBookingRestoreSpecs(
		&domain.BookingRestoreSpecs{
			ID:             uuid.New(),
			SlotID:         booking.SlotID,
			UserID:         creds.ID,
			Status:         domain.ActiveBookingStatus,
			ConferenceLink: conferenceLink,
			CreatedAt:      ptr.To(time.Now().UTC()),
		})))
	if err != nil {
		// компенсирующее действие по удалению созданной ссылки
		if conferenceLink != nil {
			cancelErr := s.linkManager.Cancel(ctx, *conferenceLink)
			if cancelErr != nil {
				s.logger.Error(fmt.Sprintf("Failed to cancel conference with link %s: %v", *conferenceLink, cancelErr))
			}
		}

		return nil, fmt.Errorf("creating booking: %w", err)
	}

	return createdBooking, nil
}

func (s *Service) List(ctx context.Context, filter *dto.BookingFilter) ([]*domain.Booking, error) {
	creds, err := domain.GetCredentialsFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	if creds.Role != domain.AdminRole {
		return nil, errs.ErrForbidden
	}

	bookings, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("listing bookings: %w", err)
	}

	return bookings, nil
}

func (s *Service) ListForUser(ctx context.Context) ([]*domain.Booking, error) {
	creds, err := domain.GetCredentialsFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	if creds.Role != domain.UserRole {
		return nil, errs.ErrForbidden
	}

	filter := &dto.BookingFilter{
		UserID: &creds.ID,
	}

	bookings, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("listing bookings: %w", err)
	}

	return bookings, nil
}

func (s *Service) Cancel(ctx context.Context, bookingID uuid.UUID) (*domain.Booking, error) {
	creds, err := domain.GetCredentialsFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	if creds.Role != domain.UserRole {
		return nil, errs.ErrForbidden
	}

	booking, err := s.repo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, fmt.Errorf("getting booking: %w", err)
	}

	if booking.UserID() != creds.ID {
		return nil, errs.ErrForbidden
	}

	if booking.ConferenceLink() != nil {
		cancelErr := s.linkManager.Cancel(ctx, *booking.ConferenceLink())
		if cancelErr != nil {
			s.logger.Error(
				fmt.Sprintf("Failed to cancel conference with link %s: %v", *booking.ConferenceLink(), cancelErr))
		}
	}

	model := &dto.BookingUpdateModel{
		ID:     booking.UserID(),
		Status: domain.CancelledBookingStatus,
	}

	updatedBooking, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, fmt.Errorf("updating booking: %w", err)
	}

	return updatedBooking, nil
}
