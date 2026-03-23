package room_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/room"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/room/mocks"
)

var errInternal = errors.New("internal error")

type testMocks struct {
	repo *mocks.Repository
}

func newMocks(t *testing.T) *testMocks {
	return &testMocks{
		repo: mocks.NewRepository(t),
	}
}

func TestService_GetAll(t *testing.T) {
	tests := map[string]struct {
		setupMocks    func(m *testMocks)
		expectedError error
	}{
		"error getting rooms": {
			setupMocks: func(m *testMocks) {
				m.repo.
					On("GetAll", mock.Anything).
					Return(nil, errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"success": {
			setupMocks: func(m *testMocks) {
				m.repo.
					On("GetAll", mock.Anything).
					Return([]*domain.Room{}, nil).
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

			s := room.NewService(m.repo)

			_, err := s.GetAll(context.Background())

			if test.expectedError != nil {
				require.ErrorIs(t, err, test.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_Create(t *testing.T) {
	tests := map[string]struct {
		ctx           context.Context
		room          *domain.RoomInitSpecs
		setupMocks    func(m *testMocks)
		expectedError error
	}{
		"no credentials in ctx": {
			ctx:           context.Background(),
			room:          &domain.RoomInitSpecs{},
			expectedError: errs.ErrUnauthorized,
		},
		"user is NOT admin, forbidden": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
				domain.UserIDKey, uuid.UUID{}),
			room:          &domain.RoomInitSpecs{},
			expectedError: errs.ErrForbidden,
		},
		"error creating room": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.AdminRole),
				domain.UserIDKey, uuid.UUID{}),
			room: &domain.RoomInitSpecs{},
			setupMocks: func(m *testMocks) {
				m.repo.
					On("Create", mock.Anything, mock.AnythingOfType("*domain.Room")).
					Return(nil, errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"success": {
			ctx: context.WithValue(context.WithValue(context.Background(), domain.UserRoleKey, domain.AdminRole),
				domain.UserIDKey, uuid.UUID{}),
			room: &domain.RoomInitSpecs{},
			setupMocks: func(m *testMocks) {
				m.repo.
					On("Create", mock.Anything, mock.AnythingOfType("*domain.Room")).
					Return(&domain.Room{}, nil).
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

			s := room.NewService(m.repo)

			_, err := s.Create(test.ctx, test.room)

			if test.expectedError != nil {
				require.ErrorIs(t, err, test.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
