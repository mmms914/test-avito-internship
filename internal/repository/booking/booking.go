package booking

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/mmms914/test-avito-internship/internal/domain"
	"github.com/mmms914/test-avito-internship/internal/dto"
	"github.com/mmms914/test-avito-internship/internal/errs"
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

func (r *Repository) ListActive(ctx context.Context, filter *dto.BookingFilter) ([]*domain.Booking, error) {
	query := `
        SELECT b.id as id,
               b.slot_id as slot_id, 
               b.user_id as user_id,
               b.status as status, 
               b.conference_link as conference_link,
               b.created_at as created_at
        FROM bookings b 
        JOIN slots s ON b.slot_id = s.id
        WHERE b.status = 'active'
    `
	args := []any{}
	argIndex := 1

	if filter.UserID != nil {
		query += fmt.Sprintf(" AND b.user_id = $%d", argIndex)
		args = append(args, *filter.UserID)
		argIndex++
	}

	if filter.Time != nil {
		query += fmt.Sprintf(" AND s.start_time >= $%d", argIndex)
		args = append(args, *filter.Time)
		argIndex++
	}

	// Сортировка и пагинация
	query += " ORDER BY b.created_at DESC"

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

	res, err := r.db.ExecContext(ctx, query,
		model.Status.String(),
		model.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update bookings: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check if updated: %w", err)
	}

	if rowsAffected == 0 {
		return errs.ErrBookingNotFound
	}

	return err
}
