package booking

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mmms914/test-avito-internship/internal/domain"
	"github.com/mmms914/test-avito-internship/internal/dto"
	"github.com/mmms914/test-avito-internship/internal/errs"
	"github.com/mmms914/test-avito-internship/pkg/ptr"
)

type Repository interface {
	GetByID(ctx context.Context, bookingID uuid.UUID) (*domain.Booking, error)
	IsSlotAlreadyBooked(ctx context.Context, slotID uuid.UUID) (bool, error)
	Create(ctx context.Context, booking *domain.Booking) error
	ListActive(ctx context.Context, filter *dto.BookingFilter) ([]*domain.Booking, error)
	Update(ctx context.Context, bum *dto.BookingUpdateModel) error
}

type SlotRepository interface {
	GetByID(ctx context.Context, slotID uuid.UUID) (*domain.Slot, error)
}

type LinkManager interface {
	Create(ctx context.Context) (string, error)
	Cancel(ctx context.Context, link string) error
}

type Logger interface {
	Error(msg string, args ...any)
}

type Service struct {
	logger Logger

	repo        Repository
	slotRepo    SlotRepository
	linkManager LinkManager
}

type Config struct {
	Logger Logger

	BookingRepo Repository
	SlotRepo    SlotRepository
	LinkManager LinkManager
}

func NewService(c *Config) *Service {
	return &Service{
		logger: c.Logger,

		repo:        c.BookingRepo,
		slotRepo:    c.SlotRepo,
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

	slot, err := s.slotRepo.GetByID(ctx, booking.SlotID)
	if err != nil {
		return nil, fmt.Errorf("getting slot: %w", err)
	}

	if slot.IsInPast() {
		return nil, errs.ErrSlotTimeInPast
	}

	isAlreadyBooked, err := s.repo.IsSlotAlreadyBooked(ctx, booking.SlotID)
	if err != nil {
		return nil, fmt.Errorf("checking if slot is already booked: %w", err)
	}

	if isAlreadyBooked {
		return nil, errs.ErrSlotAlreadyBooked
	}

	var conferenceLink *string
	if booking.CreateConferenceLink != nil && *booking.CreateConferenceLink {
		cLink, linkErr := s.linkManager.Create(ctx)
		if linkErr != nil {
			return nil, fmt.Errorf("creating conference link: %w", linkErr)
		}

		conferenceLink = &cLink
	}

	createdBooking := domain.NewBooking(domain.WithBookingRestoreSpecs(
		&domain.BookingRestoreSpecs{
			ID:             uuid.New(),
			SlotID:         booking.SlotID,
			UserID:         creds.ID,
			Status:         domain.ActiveBookingStatus,
			ConferenceLink: conferenceLink,
			CreatedAt:      time.Now().UTC(),
		}))

	if err = s.repo.Create(ctx, createdBooking); err != nil {
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

	bookings, err := s.repo.ListActive(ctx, filter)
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
		Time:   ptr.To(time.Now().UTC()),
	}

	bookings, err := s.repo.ListActive(ctx, filter)
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
		ID:     booking.ID(),
		Status: domain.CancelledBookingStatus,
	}

	booking.Cancel()

	if err = s.repo.Update(ctx, model); err != nil {
		return nil, fmt.Errorf("updating booking: %w", err)
	}

	return booking, nil
}
