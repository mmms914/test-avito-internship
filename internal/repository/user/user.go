package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
)

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Repository struct {
	db Executor
}

func NewRepository(db Executor) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, user *domain.User) error {
	query := `
        INSERT INTO users (id, email, password, role, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `

	_, err := r.db.ExecContext(ctx, query,
		user.ID(),
		user.Email(),
		user.PasswordHash(),
		user.Role().String(),
		user.CreatedAt(),
	)

	return err
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
        SELECT id, email, password, role, created_at
        FROM users
        WHERE email = $1
    `

	var specs domain.UserRestoreSpecs
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&specs.ID,
		&specs.Email,
		&specs.PasswordHash,
		&specs.Role,
		&specs.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return domain.NewUser(domain.WithUserRestoreSpecs(&specs)), nil
}

func (r *Repository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	return exists, err
}
