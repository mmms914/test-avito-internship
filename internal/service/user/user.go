package user

import (
	"context"
	"fmt"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/dto"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
)

type Repository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type Hasher interface {
	Hash(password string) (string, error)
	CompareHashAndPassword(hash string, password string) error
}

type Service struct {
	repo   Repository
	hasher Hasher
}

type Config struct {
	UserRepo Repository
	Hasher   Hasher
}

func NewService(c *Config) *Service {
	return &Service{
		repo:   c.UserRepo,
		hasher: c.Hasher,
	}
}

func (s *Service) Register(ctx context.Context, specs *dto.UserCreateModel) (*domain.User, error) {
	userExists, err := s.repo.ExistsByEmail(ctx, specs.Email)
	if err != nil {
		return nil, fmt.Errorf("checking if user exists: %w", err)
	}

	if userExists {
		return nil, errs.ErrUserAlreadyExists
	}

	passHash, err := s.hasher.Hash(specs.Password)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user := domain.NewUser(domain.WithUserInitSpecs(&domain.UserInitSpecs{
		Email:        specs.Email,
		PasswordHash: passHash,
		Role:         specs.Role,
	}))

	if dbErr := s.repo.Create(ctx, user); dbErr != nil {
		return nil, fmt.Errorf("creating user: %w", dbErr)
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, cred *dto.UserCredentials) (*domain.User, error) {
	user, err := s.repo.GetByEmail(ctx, cred.Email)
	if err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}

	if err = s.hasher.CompareHashAndPassword(user.PasswordHash(), cred.Password); err != nil {
		return nil, errs.ErrUnauthorized
	}

	return user, nil
}
