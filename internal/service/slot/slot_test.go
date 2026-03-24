package slot_test

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
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/slot"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/slot/mocks"
)

var errInternal = errors.New("internal error")

type testMocks struct {
	repo         *mocks.Repository
	scheduleRepo *mocks.ScheduleRepository
	roomRepo     *mocks.RoomRepository
}

func newMocks(t *testing.T) *testMocks {
	return &testMocks{
		repo:         mocks.NewRepository(t),
		scheduleRepo: mocks.NewScheduleRepository(t),
		roomRepo:     mocks.NewRoomRepository(t),
	}
}

func TestService_GetAvailableSlots(t *testing.T) {
	tests := map[string]struct {
		setupMocks    func(m *testMocks)
		roomID        uuid.UUID
		date          time.Time
		expectedError error
	}{
		"error getting room": {
			roomID: uuid.UUID{},
			date:   time.Time{},
			setupMocks: func(m *testMocks) {
				m.roomRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(false, errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"room does not exist": {
			roomID: uuid.UUID{},
			date:   time.Time{},
			setupMocks: func(m *testMocks) {
				m.roomRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(false, nil).
					Once()
			},
			expectedError: errs.ErrRoomNotExists,
		},
		"error getting if slots exists": {
			roomID: uuid.UUID{},
			date:   time.Time{},
			setupMocks: func(m *testMocks) {
				m.roomRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(true, nil).
					Once()
				m.repo.
					On("IsSlotsExistForDate", mock.Anything, mock.AnythingOfType("time.Time")).
					Return(false, errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"slots don't exist, error getting schedule": {
			roomID: uuid.UUID{},
			date:   time.Time{},
			setupMocks: func(m *testMocks) {
				m.roomRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(true, nil).
					Once()
				m.repo.
					On("IsSlotsExistForDate", mock.Anything, mock.AnythingOfType("time.Time")).
					Return(false, nil).
					Once()
				m.scheduleRepo.
					On("GetForRoom", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(nil, errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"slots don't exist, schedule is not for that day, success": {
			roomID: uuid.UUID{},
			date:   time.Date(2026, 3, 17, 0, 0, 0, 0, time.UTC),
			setupMocks: func(m *testMocks) {
				m.roomRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(true, nil).
					Once()
				m.repo.
					On("IsSlotsExistForDate", mock.Anything, mock.AnythingOfType("time.Time")).
					Return(false, nil).
					Once()
				m.scheduleRepo.
					On("GetForRoom", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(domain.NewSchedule(domain.WithScheduleRestoreSpecs(
						&domain.ScheduleRestoreSpecs{
							ID:         uuid.UUID{},
							RoomID:     uuid.UUID{},
							DaysOfWeek: []time.Weekday{time.Monday},
							StartTime:  7*time.Hour + 30*time.Minute,
							EndTime:    15 * time.Hour,
						})), nil).
					Once()
				m.repo.
					On("GetAllAvailable", mock.Anything, mock.AnythingOfType("*dto.SlotFilter")).
					Return([]*domain.Slot{}, nil).
					Once()
			},
			expectedError: nil,
		},
		"slots don't exist, generate slots, error creating slots in repo": {
			roomID: uuid.UUID{},
			date:   time.Date(2026, 3, 17, 0, 0, 0, 0, time.UTC),
			setupMocks: func(m *testMocks) {
				m.roomRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(true, nil).
					Once()
				m.repo.
					On("IsSlotsExistForDate", mock.Anything, mock.AnythingOfType("time.Time")).
					Return(false, nil).
					Once()
				m.scheduleRepo.
					On("GetForRoom", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(domain.NewSchedule(domain.WithScheduleRestoreSpecs(
						&domain.ScheduleRestoreSpecs{
							ID:         uuid.UUID{},
							RoomID:     uuid.UUID{},
							DaysOfWeek: []time.Weekday{time.Monday, time.Tuesday},
							StartTime:  7*time.Hour + 30*time.Minute,
							EndTime:    15 * time.Hour,
						})), nil).
					Once()
				m.repo.
					On("Create", mock.Anything, mock.AnythingOfType("[]*domain.Slot")).
					Return(errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"slots don't exist, generate slots, success": {
			roomID: uuid.UUID{},
			date:   time.Date(2026, 3, 17, 0, 0, 0, 0, time.UTC),
			setupMocks: func(m *testMocks) {
				m.roomRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(true, nil).
					Once()
				m.repo.
					On("IsSlotsExistForDate", mock.Anything, mock.AnythingOfType("time.Time")).
					Return(false, nil).
					Once()
				m.scheduleRepo.
					On("GetForRoom", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(domain.NewSchedule(domain.WithScheduleRestoreSpecs(
						&domain.ScheduleRestoreSpecs{
							ID:         uuid.UUID{},
							RoomID:     uuid.UUID{},
							DaysOfWeek: []time.Weekday{time.Monday, time.Tuesday},
							StartTime:  7*time.Hour + 30*time.Minute,
							EndTime:    15 * time.Hour,
						})), nil).
					Once()
				m.repo.
					On("Create", mock.Anything, mock.AnythingOfType("[]*domain.Slot")).
					Return(nil).
					Once()
				m.repo.
					On("GetAllAvailable", mock.Anything, mock.AnythingOfType("*dto.SlotFilter")).
					Return([]*domain.Slot{}, nil).
					Once()
			},
			expectedError: nil,
		},
		"error getting available slots": {
			roomID: uuid.UUID{},
			date:   time.Time{},
			setupMocks: func(m *testMocks) {
				m.roomRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(true, nil).
					Once()
				m.repo.
					On("IsSlotsExistForDate", mock.Anything, mock.AnythingOfType("time.Time")).
					Return(true, nil).
					Once()
				m.repo.
					On("GetAllAvailable", mock.Anything, mock.AnythingOfType("*dto.SlotFilter")).
					Return(nil, errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"success": {
			roomID: uuid.UUID{},
			date:   time.Time{},
			setupMocks: func(m *testMocks) {
				m.roomRepo.
					On("Exists", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(true, nil).
					Once()
				m.repo.
					On("IsSlotsExistForDate", mock.Anything, mock.AnythingOfType("time.Time")).
					Return(true, nil).
					Once()
				m.repo.
					On("GetAllAvailable", mock.Anything, mock.AnythingOfType("*dto.SlotFilter")).
					Return([]*domain.Slot{}, nil).
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

			s := slot.NewService(&slot.Config{
				RoomRepo:     m.roomRepo,
				ScheduleRepo: m.scheduleRepo,
				SlotRepo:     m.repo,
			})

			_, err := s.GetAvailableSlots(context.Background(), test.roomID, test.date)

			if test.expectedError != nil {
				require.ErrorIs(t, err, test.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
