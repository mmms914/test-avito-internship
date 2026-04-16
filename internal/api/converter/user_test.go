package converter_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mmms914/test-avito-internship/internal/api/converter"
	"github.com/mmms914/test-avito-internship/internal/domain"
)

func TestUserRoleFromString(t *testing.T) {
	tests := map[string]struct {
		strRole       string
		expectedRole  domain.Role
		expectedError bool
	}{
		"success - admin role": {
			strRole:       "admin",
			expectedRole:  domain.AdminRole,
			expectedError: false,
		},
		"success - user role": {
			strRole:       "user",
			expectedRole:  domain.UserRole,
			expectedError: false,
		},
		"failure - invalid role": {
			strRole:       "superuser",
			expectedRole:  "",
			expectedError: true,
		},
		"failure - empty role": {
			strRole:       "",
			expectedRole:  "",
			expectedError: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := converter.UserRoleFromString(tt.strRole)

			if tt.expectedError {
				require.Error(t, err)
				assert.Equal(t, tt.expectedRole, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedRole, result)
			}
		})
	}
}

func TestUserToResponse(t *testing.T) {
	userID := uuid.New()
	createdAt := time.Now().UTC()

	user := domain.NewUser(domain.WithUserRestoreSpecs(&domain.UserRestoreSpecs{
		ID:           userID,
		Email:        "test@example.com",
		PasswordHash: "hashed",
		Role:         domain.UserRole,
		CreatedAt:    createdAt,
	}))

	result := converter.UserToResponse(user)

	assert.NotNil(t, result)
	assert.NotNil(t, result.User)
	assert.Equal(t, userID, result.User.ID)
	assert.Equal(t, "test@example.com", result.User.Email)
	assert.Equal(t, "user", result.User.Role)
	assert.Equal(t, createdAt, result.User.CreatedAt)
}

func TestUserToObject(t *testing.T) {
	tests := map[string]struct {
		user *domain.User
	}{
		"user with all fields": {
			user: domain.NewUser(domain.WithUserRestoreSpecs(&domain.UserRestoreSpecs{
				ID:           uuid.New(),
				Email:        "user@example.com",
				PasswordHash: "hashed",
				Role:         domain.UserRole,
				CreatedAt:    time.Now().UTC(),
			})),
		},
		"admin user": {
			user: domain.NewUser(domain.WithUserRestoreSpecs(&domain.UserRestoreSpecs{
				ID:           uuid.New(),
				Email:        "admin@example.com",
				PasswordHash: "hashed",
				Role:         domain.AdminRole,
				CreatedAt:    time.Now().UTC(),
			})),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := converter.UserToObject(tt.user)

			assert.NotNil(t, result)
			assert.Equal(t, tt.user.ID(), result.ID)
			assert.Equal(t, tt.user.Email(), result.Email)
			assert.Equal(t, tt.user.Role().String(), result.Role)
			assert.Equal(t, tt.user.CreatedAt(), result.CreatedAt)
		})
	}
}
