package booking

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/dto"
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

func (r *Repository) GetByID(ctx context.Context, bookingID uuid.UUID) (*domain.Booking, error) {
	query := `
        SELECT id, slot_id, user_id, status, conference_link, created_at
        FROM bookings
        WHERE id = $1
    `

	var specs domain.BookingRestoreSpecs
	err := r.db.QueryRowContext(ctx, query, bookingID).Scan(
		&specs.ID,
		&specs.SlotID,
		&specs.UserID,
		&specs.Status,
		&specs.ConferenceLink,
		&specs.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrBookingNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get booking by id: %w", err)
	}

	return domain.NewBooking(domain.WithBookingRestoreSpecs(&specs)), nil
}

func (r *Repository) IsSlotAlreadyBooked(ctx context.Context, slotID uuid.UUID) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1 FROM bookings 
            WHERE slot_id = $1 AND status = 'active'
        )
    `

	var exists bool
	err := r.db.QueryRowContext(ctx, query, slotID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if slot is booked: %w", err)
	}

	return exists, nil
}

func (r *Repository) Create(ctx context.Context, booking *domain.Booking) error {
	query := `
        INSERT INTO bookings (id, slot_id, user_id, status, conference_link, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	_, err := r.db.ExecContext(ctx, query,
		booking.ID(),
		booking.SlotID(),
		booking.UserID(),
		booking.Status().String(),
		booking.ConferenceLink(),
		booking.CreatedAt(),
	)

	return err
}

func (r *Repository) List(ctx context.Context, filter *dto.BookingFilter) ([]*domain.Booking, error) {
	query := `
        SELECT id, slot_id, user_id, status, conference_link, created_at
        FROM bookings
        WHERE 1=1
    `
	args := []any{}
	argIndex := 1

	if filter.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *filter.UserID)
		argIndex++
	}

	if filter.Time != nil {
		query += fmt.Sprintf(" AND start_time >= $%d", argIndex)
		args = append(args, *filter.Time)
		argIndex++
	}

	// Сортировка и пагинация
	query += " ORDER BY created_at DESC"

	if filter.Page != nil && filter.PageSize != nil {
		offset := (*filter.Page - 1) * *filter.PageSize
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
		args = append(args, *filter.PageSize, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list bookings: %w", err)
	}
	defer rows.Close()

	var bookings []*domain.Booking
	for rows.Next() {
		var specs domain.BookingRestoreSpecs
		scanErr := rows.Scan(
			&specs.ID,
			&specs.SlotID,
			&specs.UserID,
			&specs.Status,
			&specs.ConferenceLink,
			&specs.CreatedAt,
		)
		if scanErr != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", scanErr)
		}
		bookings = append(bookings, domain.NewBooking(domain.WithBookingRestoreSpecs(&specs)))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return bookings, nil
}

func (r *Repository) Update(ctx context.Context, model *dto.BookingUpdateModel) error {
	query := `
        UPDATE bookings
        SET status = $1
        WHERE id = $2
    `

	_, err := r.db.ExecContext(ctx, query,
		model.Status.String(),
		model.ID,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return errs.ErrBookingNotFound
	}

	return err
}
