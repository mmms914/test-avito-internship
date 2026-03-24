package user_test

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
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/user"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/user/mocks"
)

var errInternal = errors.New("internal error")

type testMocks struct {
	repo   *mocks.Repository
	hasher *mocks.Hasher
}

func newMocks(t *testing.T) *testMocks {
	return &testMocks{
		repo:   mocks.NewRepository(t),
		hasher: mocks.NewHasher(t),
	}
}

func TestService_Register(t *testing.T) {
	tests := map[string]struct {
		sp            *dto.UserCreateModel
		setupMocks    func(m *testMocks)
		expectedError error
	}{
		"error getting user": {
			sp: &dto.UserCreateModel{
				Email:    "mail@test.ru",
				Password: "123",
				Role:     domain.AdminRole,
			},
			setupMocks: func(m *testMocks) {
				m.repo.
					On("ExistsByEmail", mock.Anything, mock.AnythingOfType("string")).
					Return(false, errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"user already exists": {
			sp: &dto.UserCreateModel{
				Email:    "mail@test.ru",
				Password: "123",
				Role:     domain.AdminRole,
			},
			setupMocks: func(m *testMocks) {
				m.repo.
					On("ExistsByEmail", mock.Anything, mock.AnythingOfType("string")).
					Return(true, nil).
					Once()
			},
			expectedError: errs.ErrUserAlreadyExists,
		},
		"error hashing password": {
			sp: &dto.UserCreateModel{
				Email:    "mail@test.ru",
				Password: "123",
				Role:     domain.AdminRole,
			},
			setupMocks: func(m *testMocks) {
				m.repo.
					On("ExistsByEmail", mock.Anything, mock.AnythingOfType("string")).
					Return(false, nil).
					Once()
				m.hasher.
					On("Hash", mock.Anything, mock.AnythingOfType("string")).
					Return("", errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"error creating user": {
			sp: &dto.UserCreateModel{
				Email:    "mail@test.ru",
				Password: "123",
				Role:     domain.AdminRole,
			},
			setupMocks: func(m *testMocks) {
				m.repo.
					On("ExistsByEmail", mock.Anything, mock.AnythingOfType("string")).
					Return(false, nil).
					Once()
				m.hasher.
					On("Hash", mock.Anything, mock.AnythingOfType("string")).
					Return("12312", nil).
					Once()
				m.repo.
					On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).
					Return(errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"success": {
			sp: &dto.UserCreateModel{
				Email:    "mail@test.ru",
				Password: "123",
				Role:     domain.AdminRole,
			},
			setupMocks: func(m *testMocks) {
				m.repo.
					On("ExistsByEmail", mock.Anything, mock.AnythingOfType("string")).
					Return(false, nil).
					Once()
				m.hasher.
					On("Hash", mock.Anything, mock.AnythingOfType("string")).
					Return("12312", nil).
					Once()
				m.repo.
					On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).
					Return(nil).
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

			s := user.NewService(&user.Config{
				UserRepo: m.repo,
				Hasher:   m.hasher,
			})

			_, err := s.Register(context.Background(), test.sp)

			if test.expectedError != nil {
				require.ErrorIs(t, err, test.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_Login(t *testing.T) {
	tests := map[string]struct {
		cred          *dto.UserCredentials
		setupMocks    func(m *testMocks)
		expectedError error
	}{
		"error getting user": {
			cred: &dto.UserCredentials{
				Email:    "mail@test.ru",
				Password: "123",
			},
			setupMocks: func(m *testMocks) {
				m.repo.
					On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"user does not exists": {
			cred: &dto.UserCredentials{
				Email:    "mail@test.ru",
				Password: "123",
			},
			setupMocks: func(m *testMocks) {
				m.repo.
					On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, nil).
					Once()
			},
			expectedError: errs.ErrUserNotFound,
		},
		"error hashing password": {
			cred: &dto.UserCredentials{
				Email:    "mail@test.ru",
				Password: "123",
			},
			setupMocks: func(m *testMocks) {
				m.repo.
					On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).
					Return(&domain.User{}, nil).
					Once()
				m.hasher.
					On("Hash", mock.Anything, mock.AnythingOfType("string")).
					Return("", errInternal).
					Once()
			},
			expectedError: errInternal,
		},
		"hashes not equal": {
			cred: &dto.UserCredentials{
				Email:    "mail@test.ru",
				Password: "123",
			},
			setupMocks: func(m *testMocks) {
				m.repo.
					On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).
					Return(domain.NewUser(domain.WithUserRestoreSpecs(&domain.UserRestoreSpecs{
						ID:           uuid.UUID{},
						Email:        "mail@test.ru",
						PasswordHash: "123",
						Role:         domain.AdminRole,
						CreatedAt:    time.Time{},
					})), nil).
					Once()
				m.hasher.
					On("Hash", mock.Anything, mock.AnythingOfType("string")).
					Return("321", nil).
					Once()
			},
			expectedError: errs.ErrUnauthorized,
		},
		"success": {
			cred: &dto.UserCredentials{
				Email:    "mail@test.ru",
				Password: "123",
			},
			setupMocks: func(m *testMocks) {
				m.repo.
					On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).
					Return(domain.NewUser(domain.WithUserRestoreSpecs(&domain.UserRestoreSpecs{
						ID:           uuid.UUID{},
						Email:        "mail@test.ru",
						PasswordHash: "123",
						Role:         domain.AdminRole,
						CreatedAt:    time.Time{},
					})), nil).
					Once()
				m.hasher.
					On("Hash", mock.Anything, mock.AnythingOfType("string")).
					Return("123", nil).
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

			s := user.NewService(&user.Config{
				UserRepo: m.repo,
				Hasher:   m.hasher,
			})

			_, err := s.Login(context.Background(), test.cred)

			if test.expectedError != nil {
				require.ErrorIs(t, err, test.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
