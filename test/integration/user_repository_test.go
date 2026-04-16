//go:build integration
// +build integration

package integration_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mmms914/test-avito-internship/internal/domain"
	"github.com/mmms914/test-avito-internship/internal/errs"
)

func (s *IntegrationTestSuite) TestUserRepository_Create() {
	userID := uuid.New()
	user := domain.NewUser(domain.WithUserRestoreSpecs(&domain.UserRestoreSpecs{
		ID:           userID,
		Email:        "create@example.com",
		PasswordHash: "hashedPasswordHash",
		Role:         domain.UserRole,
		CreatedAt:    time.Now().UTC(),
	}))

	err := s.userRepo.Create(s.ctx, user)
	assert.NoError(s.T(), err)
}

func (s *IntegrationTestSuite) TestUserRepository_GetByEmail() {
	userID := uuid.New()
	email := "getbyemail@example.com"
	user := domain.NewUser(domain.WithUserRestoreSpecs(&domain.UserRestoreSpecs{
		ID:           userID,
		Email:        email,
		PasswordHash: "hashedPasswordHash",
		Role:         domain.UserRole,
		CreatedAt:    time.Now().UTC(),
	}))

	err := s.userRepo.Create(s.ctx, user)
	require.NoError(s.T(), err)

	found, err := s.userRepo.GetByEmail(s.ctx, email)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), userID, found.ID())
	assert.Equal(s.T(), email, found.Email())
	assert.Equal(s.T(), domain.UserRole, found.Role())
}

func (s *IntegrationTestSuite) TestUserRepository_GetByEmail_NotFound() {
	_, err := s.userRepo.GetByEmail(s.ctx, "nonexistent@example.com")
	assert.ErrorIs(s.T(), err, errs.ErrUserNotFound)
}

func (s *IntegrationTestSuite) TestUserRepository_ExistsByEmail() {
	email := "exists@example.com"
	user := domain.NewUser(domain.WithUserRestoreSpecs(&domain.UserRestoreSpecs{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: "hashedPasswordHash",
		Role:         domain.UserRole,
		CreatedAt:    time.Now().UTC(),
	}))

	exists, err := s.userRepo.ExistsByEmail(s.ctx, email)
	assert.NoError(s.T(), err)
	assert.False(s.T(), exists)

	err = s.userRepo.Create(s.ctx, user)
	require.NoError(s.T(), err)

	exists, err = s.userRepo.ExistsByEmail(s.ctx, email)
	assert.NoError(s.T(), err)
	assert.True(s.T(), exists)
}

func (s *IntegrationTestSuite) TestUserRepository_Create_Admin() {
	userID := uuid.New()
	user := domain.NewUser(domain.WithUserRestoreSpecs(&domain.UserRestoreSpecs{
		ID:           userID,
		Email:        "admin@example.com",
		PasswordHash: "hashedPasswordHash",
		Role:         domain.AdminRole,
		CreatedAt:    time.Now().UTC(),
	}))

	err := s.userRepo.Create(s.ctx, user)
	assert.NoError(s.T(), err)

	found, err := s.userRepo.GetByEmail(s.ctx, "admin@example.com")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), domain.AdminRole, found.Role())
}

func (s *IntegrationTestSuite) TestUserRepository_GetByEmail_ReturnsCorrectUser() {
	user1 := domain.NewUser(domain.WithUserRestoreSpecs(&domain.UserRestoreSpecs{
		ID:           uuid.New(),
		Email:        "user1@example.com",
		PasswordHash: "hashed1",
		Role:         domain.UserRole,
		CreatedAt:    time.Now().UTC(),
	}))
	user2 := domain.NewUser(domain.WithUserRestoreSpecs(&domain.UserRestoreSpecs{
		ID:           uuid.New(),
		Email:        "user2@example.com",
		PasswordHash: "hashed2",
		Role:         domain.UserRole,
		CreatedAt:    time.Now().UTC(),
	}))

	err := s.userRepo.Create(s.ctx, user1)
	require.NoError(s.T(), err)
	err = s.userRepo.Create(s.ctx, user2)
	require.NoError(s.T(), err)

	found, err := s.userRepo.GetByEmail(s.ctx, "user1@example.com")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user1.ID(), found.ID())

	found, err = s.userRepo.GetByEmail(s.ctx, "user2@example.com")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user2.ID(), found.ID())
}
