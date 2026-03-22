package schedule_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/schedule"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/schedule/mocks"
)

var errInternal = errors.New("internal error")

type testMocks struct {
	repo     *mocks.Repository
	roomRepo *mocks.RoomRepository
}

func newMocks(t *testing.T) *testMocks {
	return &testMocks{
		repo:     mocks.NewRepository(t),
		roomRepo: mocks.NewRoomRepository(t),
	}
}

func TestService_Create(t *testing.T) {
	tests := map[string]struct {
		ctx           context.Context
		schedule      *domain.Schedule
		setupMocks    func(m *testMocks)
		expectedError error
	}{
		"no credentials in ctx": {
			ctx:           context.Background(),
			schedule:      &domain.Schedule{},
			expectedError: errs.ErrUnauthorized,
		},
		"user is NOT admin, forbidden": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			schedule:      &domain.Schedule{},
			expectedError: errs.ErrForbidden,
		},
		"error getting room": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.AdminRole),
				domain.UserIDKey, uuid.UUID{}),
			schedule: domain.NewSchedule(domain.WithScheduleRestoreSpecs(
				&domain.ScheduleRestoreSpecs{
					RoomID:     uuid.UUID{},
					DaysOfWeek: []domain.Day{domain.Monday, domain.Wednesday, domain.Friday},
					StartTime:  time.Time{},
					EndTime:    time.Time{},
				})),
			setupMocks: func(m *testMocks) {
				m.roomRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(false, errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"room does not exist": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.AdminRole),
				domain.UserIDKey, uuid.UUID{}),
			schedule: domain.NewSchedule(domain.WithScheduleRestoreSpecs(
				&domain.ScheduleRestoreSpecs{
					RoomID:     uuid.UUID{},
					DaysOfWeek: []domain.Day{domain.Monday, domain.Wednesday, domain.Friday},
					StartTime:  time.Time{},
					EndTime:    time.Time{},
				})),
			setupMocks: func(m *testMocks) {
				m.roomRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(false, nil).
					Once()
			},
			expectedError: errs.ErrRoomNotExists,
		},
		"error checking if schedule exists for room": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.AdminRole),
				domain.UserIDKey, uuid.UUID{}),
			schedule: domain.NewSchedule(domain.WithScheduleRestoreSpecs(
				&domain.ScheduleRestoreSpecs{
					RoomID:     uuid.UUID{},
					DaysOfWeek: []domain.Day{domain.Monday, domain.Wednesday, domain.Friday},
					StartTime:  time.Time{},
					EndTime:    time.Time{},
				})),
			setupMocks: func(m *testMocks) {
				m.roomRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(true, nil).
					Once()
				m.repo.
					On("ExistsForRoom", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(false, errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"schedule already exists": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.AdminRole),
				domain.UserIDKey, uuid.UUID{}),
			schedule: domain.NewSchedule(domain.WithScheduleRestoreSpecs(
				&domain.ScheduleRestoreSpecs{
					RoomID:     uuid.UUID{},
					DaysOfWeek: []domain.Day{domain.Monday, domain.Wednesday, domain.Friday},
					StartTime:  time.Time{},
					EndTime:    time.Time{},
				})),
			setupMocks: func(m *testMocks) {
				m.roomRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(true, nil).
					Once()
				m.repo.
					On("ExistsForRoom", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(true, nil).
					Once()
			},
			expectedError: errs.ErrScheduleExists,
		},
		"error creating schedule": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.AdminRole),
				domain.UserIDKey, uuid.UUID{}),
			schedule: domain.NewSchedule(domain.WithScheduleRestoreSpecs(
				&domain.ScheduleRestoreSpecs{
					RoomID:     uuid.UUID{},
					DaysOfWeek: []domain.Day{domain.Monday, domain.Wednesday, domain.Friday},
					StartTime:  time.Time{},
					EndTime:    time.Time{},
				})),
			setupMocks: func(m *testMocks) {
				m.roomRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(true, nil).
					Once()
				m.repo.
					On("ExistsForRoom", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(false, nil).
					Once()
				m.repo.
					On("Create", mock.Anything, mock.AnythingOfType("*domain.Schedule")).
					Return(nil, errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"success": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.AdminRole),
				domain.UserIDKey, uuid.UUID{}),
			schedule: domain.NewSchedule(domain.WithScheduleRestoreSpecs(
				&domain.ScheduleRestoreSpecs{
					RoomID:     uuid.UUID{},
					DaysOfWeek: []domain.Day{domain.Monday, domain.Wednesday, domain.Friday},
					StartTime:  time.Time{},
					EndTime:    time.Time{},
				})),
			setupMocks: func(m *testMocks) {
				m.roomRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(true, nil).
					Once()
				m.repo.
					On("ExistsForRoom", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(false, nil).
					Once()
				m.repo.
					On("Create", mock.Anything, mock.AnythingOfType("*domain.Schedule")).
					Return(&domain.Schedule{}, nil).
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

			s := schedule.NewService(&schedule.Config{
				RoomRepo: m.roomRepo,
				Repo:     m.repo,
			})

			_, err := s.Create(test.ctx, test.schedule)

			if test.expectedError != nil {
				require.ErrorIs(t, err, test.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
