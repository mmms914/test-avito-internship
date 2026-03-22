package booking_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/dto"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/booking"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/booking/mocks"
	"github.com/avito-internships/test-backend-1-mmms914/pkg/ptr"
)

var (
	errInternal      = errors.New("internal error")
	errOtherInternal = errors.New("other internal error")
)

type testMocks struct {
	logger *mocks.Logger

	repo        *mocks.Repository
	slotRepo    *mocks.SlotRepository
	userRepo    *mocks.UserRepository
	linkManager *mocks.LinkManager
}

func newMocks(t *testing.T) *testMocks {
	return &testMocks{
		logger: mocks.NewLogger(t),

		repo:        mocks.NewRepository(t),
		slotRepo:    mocks.NewSlotRepository(t),
		userRepo:    mocks.NewUserRepository(t),
		linkManager: mocks.NewLinkManager(t),
	}
}

func TestService_Create(t *testing.T) {
	tests := map[string]struct {
		ctx           context.Context
		model         *dto.BookingCreateModel
		setupMocks    func(m *testMocks)
		expectedError error
	}{
		"no credentials in ctx": {
			ctx: context.Background(),
			model: &dto.BookingCreateModel{
				SlotID:               uuid.UUID{},
				CreateConferenceLink: ptr.To(true),
			},
			expectedError: errs.ErrUnauthorized,
		},
		"user is admin, forbidden": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.AdminRole),
				domain.UserIDKey, uuid.UUID{}),
			model: &dto.BookingCreateModel{
				SlotID:               uuid.UUID{},
				CreateConferenceLink: ptr.To(true),
			},
			expectedError: errs.ErrForbidden,
		},
		"error getting slot": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			model: &dto.BookingCreateModel{
				SlotID:               uuid.UUID{},
				CreateConferenceLink: ptr.To(true),
			},
			setupMocks: func(m *testMocks) {
				m.slotRepo.
					On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(nil, errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"slot does not exist": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			model: &dto.BookingCreateModel{
				SlotID:               uuid.UUID{},
				CreateConferenceLink: ptr.To(true),
			},
			setupMocks: func(m *testMocks) {
				m.slotRepo.
					On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(nil, errs.ErrSlotNotFound).
					Once()
			},
			expectedError: errs.ErrSlotNotFound,
		},
		"slot is in past": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			model: &dto.BookingCreateModel{
				SlotID:               uuid.UUID{},
				CreateConferenceLink: ptr.To(true),
			},
			setupMocks: func(m *testMocks) {
				m.slotRepo.
					On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(domain.NewSlot(domain.WithSlotRestoreSpecs(
						&domain.SlotRestoreSpecs{
							ID:        uuid.UUID{},
							RoomID:    uuid.UUID{},
							StartTime: time.Now().Add(-time.Hour).UTC(),
							EndTime:   time.Now().Add(-time.Minute).UTC(),
						})), nil).
					Once()
			},
			expectedError: errs.ErrSlotTimeInPast,
		},
		"error checking if slot already booked": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			model: &dto.BookingCreateModel{
				SlotID:               uuid.UUID{},
				CreateConferenceLink: ptr.To(true),
			},
			setupMocks: func(m *testMocks) {
				m.slotRepo.
					On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(domain.NewSlot(domain.WithSlotRestoreSpecs(
						&domain.SlotRestoreSpecs{
							ID:        uuid.UUID{},
							RoomID:    uuid.UUID{},
							StartTime: time.Now().Add(time.Hour).UTC(),
							EndTime:   time.Now().Add(2 * time.Hour).UTC(),
						})), nil).
					Once()
				m.repo.
					On("IsSlotAlreadyBooked", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(false, errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"slot already booked": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			model: &dto.BookingCreateModel{
				SlotID:               uuid.UUID{},
				CreateConferenceLink: ptr.To(true),
			},
			setupMocks: func(m *testMocks) {
				m.slotRepo.
					On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(domain.NewSlot(domain.WithSlotRestoreSpecs(
						&domain.SlotRestoreSpecs{
							ID:        uuid.UUID{},
							RoomID:    uuid.UUID{},
							StartTime: time.Now().Add(time.Hour).UTC(),
							EndTime:   time.Now().Add(2 * time.Hour).UTC(),
						})), nil).
					Once()
				m.repo.
					On("IsSlotAlreadyBooked", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(true, nil).
					Once()
			},
			expectedError: errs.ErrSlotAlreadyBooked,
		},
		"error getting user": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			model: &dto.BookingCreateModel{
				SlotID:               uuid.UUID{},
				CreateConferenceLink: ptr.To(true),
			},
			setupMocks: func(m *testMocks) {
				m.slotRepo.
					On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(domain.NewSlot(domain.WithSlotRestoreSpecs(
						&domain.SlotRestoreSpecs{
							ID:        uuid.UUID{},
							RoomID:    uuid.UUID{},
							StartTime: time.Now().Add(time.Hour).UTC(),
							EndTime:   time.Now().Add(2 * time.Hour).UTC(),
						})), nil).
					Once()
				m.repo.
					On("IsSlotAlreadyBooked", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(false, nil).
					Once()
				m.userRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(false, errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"user does not exist": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			model: &dto.BookingCreateModel{
				SlotID:               uuid.UUID{},
				CreateConferenceLink: ptr.To(true),
			},
			setupMocks: func(m *testMocks) {
				m.slotRepo.
					On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(domain.NewSlot(domain.WithSlotRestoreSpecs(
						&domain.SlotRestoreSpecs{
							ID:        uuid.UUID{},
							RoomID:    uuid.UUID{},
							StartTime: time.Now().Add(time.Hour).UTC(),
							EndTime:   time.Now().Add(2 * time.Hour).UTC(),
						})), nil).
					Once()
				m.repo.
					On("IsSlotAlreadyBooked", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(false, nil).
					Once()
				m.userRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(false, nil).
					Once()
			},
			expectedError: errs.ErrUserNotFound,
		},
		"error creating conference link": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			model: &dto.BookingCreateModel{
				SlotID:               uuid.UUID{},
				CreateConferenceLink: ptr.To(true),
			},
			setupMocks: func(m *testMocks) {
				m.slotRepo.
					On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(domain.NewSlot(domain.WithSlotRestoreSpecs(
						&domain.SlotRestoreSpecs{
							ID:        uuid.UUID{},
							RoomID:    uuid.UUID{},
							StartTime: time.Now().Add(time.Hour).UTC(),
							EndTime:   time.Now().Add(2 * time.Hour).UTC(),
						})), nil).
					Once()
				m.repo.
					On("IsSlotAlreadyBooked", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(false, nil).
					Once()
				m.userRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(true, nil).
					Once()
				m.linkManager.
					On("Create", mock.Anything).
					Return("", errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"error creating booking in repo, cancelling of conference is success": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			model: &dto.BookingCreateModel{
				SlotID:               uuid.UUID{},
				CreateConferenceLink: ptr.To(true),
			},
			setupMocks: func(m *testMocks) {
				m.slotRepo.
					On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(domain.NewSlot(domain.WithSlotRestoreSpecs(
						&domain.SlotRestoreSpecs{
							ID:        uuid.UUID{},
							RoomID:    uuid.UUID{},
							StartTime: time.Now().Add(time.Hour).UTC(),
							EndTime:   time.Now().Add(2 * time.Hour).UTC(),
						})), nil).
					Once()
				m.repo.
					On("IsSlotAlreadyBooked", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(false, nil).
					Once()
				m.userRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(true, nil).
					Once()
				m.linkManager.
					On("Create", mock.Anything).
					Return("link", nil).
					Once()
				m.repo.
					On("Create", mock.Anything, mock.AnythingOfType("*domain.Booking")).
					Return(nil, errInternal).
					Once()
				m.linkManager.
					On("Cancel", mock.Anything, mock.AnythingOfType("string")).
					Return(nil).
					Once()
			},
			expectedError: errInternal,
		},
		"error creating booking in repo, cancelling of conference is NOT success": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			model: &dto.BookingCreateModel{
				SlotID:               uuid.UUID{},
				CreateConferenceLink: ptr.To(true),
			},
			setupMocks: func(m *testMocks) {
				m.slotRepo.
					On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(domain.NewSlot(domain.WithSlotRestoreSpecs(
						&domain.SlotRestoreSpecs{
							ID:        uuid.UUID{},
							RoomID:    uuid.UUID{},
							StartTime: time.Now().Add(time.Hour).UTC(),
							EndTime:   time.Now().Add(2 * time.Hour).UTC(),
						})), nil).
					Once()
				m.repo.
					On("IsSlotAlreadyBooked", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(false, nil).
					Once()
				m.userRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(true, nil).
					Once()
				m.linkManager.
					On("Create", mock.Anything).
					Return("link", nil).
					Once()
				m.repo.
					On("Create", mock.Anything, mock.AnythingOfType("*domain.Booking")).
					Return(nil, errInternal).
					Once()
				m.linkManager.
					On("Cancel", mock.Anything, mock.AnythingOfType("string")).
					Return(errOtherInternal).
					Once()
				m.logger.
					On("Error", mock.AnythingOfType("string")).
					Once()
			},
			expectedError: errInternal,
		},
		"success": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			model: &dto.BookingCreateModel{
				SlotID:               uuid.UUID{},
				CreateConferenceLink: ptr.To(true),
			},
			setupMocks: func(m *testMocks) {
				m.slotRepo.
					On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(domain.NewSlot(domain.WithSlotRestoreSpecs(
						&domain.SlotRestoreSpecs{
							ID:        uuid.UUID{},
							RoomID:    uuid.UUID{},
							StartTime: time.Now().Add(time.Hour).UTC(),
							EndTime:   time.Now().Add(2 * time.Hour).UTC(),
						})), nil).
					Once()
				m.repo.
					On("IsSlotAlreadyBooked", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(false, nil).
					Once()
				m.userRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(true, nil).
					Once()
				m.linkManager.
					On("Create", mock.Anything).
					Return("link", nil).
					Once()
				m.repo.
					On("Create", mock.Anything, mock.AnythingOfType("*domain.Booking")).
					Return(&domain.Booking{}, nil).
					Once()
			},
			expectedError: nil,
		},
		"success, conference link is not needed": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			model: &dto.BookingCreateModel{
				SlotID:               uuid.UUID{},
				CreateConferenceLink: ptr.To(false),
			},
			setupMocks: func(m *testMocks) {
				m.slotRepo.
					On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(domain.NewSlot(domain.WithSlotRestoreSpecs(
						&domain.SlotRestoreSpecs{
							ID:        uuid.UUID{},
							RoomID:    uuid.UUID{},
							StartTime: time.Now().Add(time.Hour).UTC(),
							EndTime:   time.Now().Add(2 * time.Hour).UTC(),
						})), nil).
					Once()
				m.repo.
					On("IsSlotAlreadyBooked", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(false, nil).
					Once()
				m.userRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(true, nil).
					Once()
				m.linkManager.
					On("Create", mock.Anything).
					Return("", errInternal).
					Maybe()
				m.repo.
					On("Create", mock.Anything, mock.AnythingOfType("*domain.Booking")).
					Return(&domain.Booking{}, nil).
					Once()
			},
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := newMocks(t)

			if test.setupMocks != nil {
				test.setupMocks(m)
			}

			s := booking.NewService(&booking.Config{
				Logger:      m.logger,
				BookingRepo: m.repo,
				SlotRepo:    m.slotRepo,
				UserRepo:    m.userRepo,
				LinkManager: m.linkManager,
			})

			_, err := s.Create(test.ctx, test.model)

			if test.expectedError != nil {
				require.ErrorIs(t, err, test.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_List(t *testing.T) {
	tests := map[string]struct {
		ctx           context.Context
		filter        *dto.BookingFilter
		setupMocks    func(m *testMocks)
		expectedError error
	}{
		"no credentials in ctx": {
			ctx: context.Background(),
			filter: &dto.BookingFilter{
				Page:     ptr.To(1),
				PageSize: ptr.To(20),
			},
			expectedError: errs.ErrUnauthorized,
		},
		"user is NOT admin, forbidden": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			filter: &dto.BookingFilter{
				Page:     ptr.To(1),
				PageSize: ptr.To(20),
			},
			expectedError: errs.ErrForbidden,
		},
		"error getting booking list": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.AdminRole),
				domain.UserIDKey, uuid.UUID{}),
			filter: &dto.BookingFilter{
				Page:     ptr.To(1),
				PageSize: ptr.To(20),
			},
			setupMocks: func(m *testMocks) {
				m.repo.
					On("List", mock.Anything, mock.AnythingOfType("*dto.BookingFilter")).
					Return([]*domain.Booking{}, errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"success": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.AdminRole),
				domain.UserIDKey, uuid.UUID{}),
			filter: &dto.BookingFilter{
				Page:     ptr.To(1),
				PageSize: ptr.To(20),
			},
			setupMocks: func(m *testMocks) {
				m.repo.
					On("List", mock.Anything, mock.AnythingOfType("*dto.BookingFilter")).
					Return([]*domain.Booking{}, nil).
					Once()
			},
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := newMocks(t)

			if test.setupMocks != nil {
				test.setupMocks(m)
			}

			s := booking.NewService(&booking.Config{
				Logger:      m.logger,
				BookingRepo: m.repo,
				SlotRepo:    m.slotRepo,
				UserRepo:    m.userRepo,
				LinkManager: m.linkManager,
			})

			_, err := s.List(test.ctx, test.filter)

			if test.expectedError != nil {
				require.ErrorIs(t, err, test.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_ListForUser(t *testing.T) {
	tests := map[string]struct {
		ctx           context.Context
		setupMocks    func(m *testMocks)
		expectedError error
	}{
		"no credentials in ctx": {
			ctx:           context.Background(),
			expectedError: errs.ErrUnauthorized,
		},
		"user is admin, forbidden": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.AdminRole),
				domain.UserIDKey, uuid.UUID{}),
			expectedError: errs.ErrForbidden,
		},
		"error getting booking list": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			setupMocks: func(m *testMocks) {
				m.repo.
					On("List", mock.Anything, mock.AnythingOfType("*dto.BookingFilter")).
					Return(nil, errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"success": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			setupMocks: func(m *testMocks) {
				m.repo.
					On("List", mock.Anything, mock.AnythingOfType("*dto.BookingFilter")).
					Return([]*domain.Booking{}, nil).
					Once()
			},
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := newMocks(t)

			if test.setupMocks != nil {
				test.setupMocks(m)
			}

			s := booking.NewService(&booking.Config{
				Logger:      m.logger,
				BookingRepo: m.repo,
				SlotRepo:    m.slotRepo,
				UserRepo:    m.userRepo,
				LinkManager: m.linkManager,
			})

			_, err := s.ListForUser(test.ctx)

			if test.expectedError != nil {
				require.ErrorIs(t, err, test.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_Cancel(t *testing.T) {
	tests := map[string]struct {
		ctx           context.Context
		bookingID     uuid.UUID
		setupMocks    func(m *testMocks)
		expectedError error
	}{
		"no credentials in ctx": {
			ctx:           context.Background(),
			bookingID:     uuid.UUID{},
			expectedError: errs.ErrUnauthorized,
		},
		"user is admin, forbidden": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.AdminRole),
				domain.UserIDKey, uuid.UUID{}),
			bookingID:     uuid.UUID{},
			expectedError: errs.ErrForbidden,
		},
		"error getting booking": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			bookingID: uuid.UUID{},
			setupMocks: func(m *testMocks) {
				m.repo.
					On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(nil, errs.ErrBookingNotFound).
					Once()
			},
			expectedError: errs.ErrBookingNotFound,
		},
		"user is not owner of booking": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			bookingID: uuid.UUID{},
			setupMocks: func(m *testMocks) {
				m.repo.
					On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(domain.NewBooking(domain.WithBookingRestoreSpecs(
						&domain.BookingRestoreSpecs{
							ID:        uuid.UUID{},
							SlotID:    uuid.UUID{},
							UserID:    uuid.New(),
							Status:    domain.ActiveBookingStatus,
							CreatedAt: ptr.To(time.Now().UTC()),
						})), nil).
					Once()
			},
			expectedError: errs.ErrForbidden,
		},
		"error cancelling conference, other is good": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			bookingID: uuid.UUID{},
			setupMocks: func(m *testMocks) {
				m.repo.
					On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(domain.NewBooking(domain.WithBookingRestoreSpecs(
						&domain.BookingRestoreSpecs{
							ID:             uuid.UUID{},
							SlotID:         uuid.UUID{},
							UserID:         uuid.UUID{},
							Status:         domain.ActiveBookingStatus,
							ConferenceLink: ptr.To("link"),
							CreatedAt:      ptr.To(time.Now().UTC()),
						})), nil).
					Once()
				m.linkManager.
					On("Cancel", mock.Anything, mock.AnythingOfType("string")).
					Return(errOtherInternal).
					Once()
				m.logger.
					On("Error", mock.AnythingOfType("string")).
					Once()
				m.repo.
					On("Update", mock.Anything, mock.AnythingOfType("*dto.BookingUpdateModel")).
					Return(&domain.Booking{}, nil).
					Once()
			},
			expectedError: nil,
		},
		"error updating booking": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			bookingID: uuid.UUID{},
			setupMocks: func(m *testMocks) {
				m.repo.
					On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(domain.NewBooking(domain.WithBookingRestoreSpecs(
						&domain.BookingRestoreSpecs{
							ID:        uuid.UUID{},
							SlotID:    uuid.UUID{},
							UserID:    uuid.UUID{},
							Status:    domain.ActiveBookingStatus,
							CreatedAt: ptr.To(time.Now().UTC()),
						})), nil).
					Once()

				m.repo.
					On("Update", mock.Anything, mock.AnythingOfType("*dto.BookingUpdateModel")).
					Return(nil, errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"success, booking is active": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			bookingID: uuid.UUID{},
			setupMocks: func(m *testMocks) {
				m.repo.
					On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(domain.NewBooking(domain.WithBookingRestoreSpecs(
						&domain.BookingRestoreSpecs{
							ID:        uuid.UUID{},
							SlotID:    uuid.UUID{},
							UserID:    uuid.UUID{},
							Status:    domain.ActiveBookingStatus,
							CreatedAt: ptr.To(time.Now().UTC()),
						})), nil).
					Once()

				m.repo.
					On("Update", mock.Anything, mock.AnythingOfType("*dto.BookingUpdateModel")).
					Return(&domain.Booking{}, nil).
					Once()
			},
			expectedError: nil,
		},
		"success, booking is cancelled": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			bookingID: uuid.UUID{},
			setupMocks: func(m *testMocks) {
				m.repo.
					On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(domain.NewBooking(domain.WithBookingRestoreSpecs(
						&domain.BookingRestoreSpecs{
							ID:             uuid.UUID{},
							SlotID:         uuid.UUID{},
							UserID:         uuid.UUID{},
							Status:         domain.CancelledBookingStatus,
							ConferenceLink: ptr.To("link"),
							CreatedAt:      ptr.To(time.Now().UTC()),
						})), nil).
					Once()
				m.linkManager.
					On("Cancel", mock.Anything, mock.AnythingOfType("string")).
					Return(nil).
					Once()
				m.repo.
					On("Update", mock.Anything, mock.AnythingOfType("*dto.BookingUpdateModel")).
					Return(&domain.Booking{}, nil).
					Once()
			},
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := newMocks(t)

			if test.setupMocks != nil {
				test.setupMocks(m)
			}

			s := booking.NewService(&booking.Config{
				Logger:      m.logger,
				BookingRepo: m.repo,
				SlotRepo:    m.slotRepo,
				UserRepo:    m.userRepo,
				LinkManager: m.linkManager,
			})

			_, err := s.Cancel(test.ctx, test.bookingID)

			if test.expectedError != nil {
				require.ErrorIs(t, err, test.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
