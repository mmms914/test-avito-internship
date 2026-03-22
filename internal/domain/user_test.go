package domain_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/pkg/ptr"
)

func TestUser_GetCredentialsFromContext(t *testing.T) {
	tests := map[string]struct {
		ctx          context.Context
		expectedErr  bool
		expectedID   *uuid.UUID
		expectedRole *domain.Role
	}{
		"ID and role": {
			ctx: context.WithValue(
				context.WithValue(
					context.Background(), domain.UserIDKey, uuid.UUID{}), domain.UserRoleKey, domain.AdminRole),
			expectedErr:  false,
			expectedID:   ptr.To(uuid.UUID{}),
			expectedRole: ptr.To(domain.AdminRole),
		},
		"only ID": {
			ctx:         context.WithValue(context.Background(), domain.UserIDKey, uuid.UUID{}),
			expectedErr: true,
		},
		"only role": {
			ctx:         context.WithValue(context.Background(), domain.UserRoleKey, domain.UserRole),
			expectedErr: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			auth, err := domain.GetCredentialsFromContext(test.ctx)
			if test.expectedErr {
				require.Error(t, err)
			}

			if test.expectedID != nil {
				assert.Equal(t, *test.expectedID, auth.ID)
			}
			if test.expectedRole != nil {
				assert.Equal(t, *test.expectedRole, auth.Role)
			}
		})
	}
	t.Run("success", func(t *testing.T) {
		uid := uuid.New()

		ctx := context.WithValue(
			context.WithValue(context.Background(), domain.UserIDKey, uid), domain.UserRoleKey, domain.AdminRole)
		auth, err := domain.GetCredentialsFromContext(ctx)
		require.NoError(t, err)

		assert.Equal(t, domain.AdminRole, auth.Role)
		assert.Equal(t, uid, auth.ID)
	})
}

func TestNewUser_WithInitSpecs(t *testing.T) {
	email := "mail"
	passwordHash := "password123"
	role := domain.AdminRole

	user := domain.NewUser(
		domain.WithUserInitSpecs(&domain.UserInitSpecs{
			Email:        email,
			PasswordHash: passwordHash,
			Role:         role,
		}),
	)

	assert.Equal(t, email, user.Email())
	assert.Equal(t, passwordHash, user.PasswordHash())
	assert.Equal(t, role, user.Role())
}

func TestNewUser_WithRestoreSpecs(t *testing.T) {
	id := uuid.New()
	email := "mail"
	passwordHash := "password123"
	role := domain.AdminRole
	createdAt := time.Now()

	user := domain.NewUser(
		domain.WithUserRestoreSpecs(&domain.UserRestoreSpecs{
			ID:           id,
			Email:        email,
			PasswordHash: passwordHash,
			Role:         role,
			CreatedAt:    createdAt,
		}),
	)

	assert.Equal(t, id, user.ID())
	assert.Equal(t, email, user.Email())
	assert.Equal(t, passwordHash, user.PasswordHash())
	assert.Equal(t, role, user.Role())
	assert.Equal(t, createdAt, user.CreatedAt())
}
