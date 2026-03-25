package slot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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

func (r *Repository) Create(ctx context.Context, slots []*domain.Slot) error {
	if len(slots) == 0 {
		return nil
	}

	query := `
        INSERT INTO slots (id, room_id, start_time, end_time)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (room_id, start_time) DO NOTHING
    `

	for _, slot := range slots {
		_, err := r.db.ExecContext(ctx, query,
			slot.ID(),
			slot.RoomID(),
			slot.StartTime(),
			slot.EndTime(),
		)
		if err != nil {
			return fmt.Errorf("failed to create slot: %w", err)
		}
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, slotID uuid.UUID) (*domain.Slot, error) {
	query := `
        SELECT id, room_id, start_time, end_time
        FROM slots
        WHERE id = $1
    `

	var specs domain.SlotRestoreSpecs
	err := r.db.QueryRowContext(ctx, query, slotID).Scan(
		&specs.ID,
		&specs.RoomID,
		&specs.StartTime,
		&specs.EndTime,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrSlotNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get slot by id: %w", err)
	}

	return domain.NewSlot(domain.WithSlotRestoreSpecs(&specs)), nil
}
func (r *Repository) GetAllAvailable(ctx context.Context, filter *dto.SlotFilter) ([]*domain.Slot, error) {
	query := `
        SELECT s.id, s.room_id, s.start_time, s.end_time
        FROM slots s
        WHERE s.room_id = $1
          AND s.start_time >= $2
          AND s.start_time < $3
          AND NOT EXISTS (
              SELECT 1 FROM bookings b 
              WHERE b.slot_id = s.id 
                AND b.status = 'active'
          )
        ORDER BY s.start_time ASC
    `

	rows, err := r.db.QueryContext(ctx, query,
		filter.RoomID,
		filter.GetStartOfDay(),
		filter.GetEndOfDay(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get available slots: %w", err)
	}
	defer rows.Close()

	var slots []*domain.Slot
	for rows.Next() {
		var specs domain.SlotRestoreSpecs
		err = rows.Scan(
			&specs.ID,
			&specs.RoomID,
			&specs.StartTime,
			&specs.EndTime,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan slot: %w", err)
		}
		slots = append(slots, domain.NewSlot(domain.WithSlotRestoreSpecs(&specs)))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return slots, nil
}

func (r *Repository) IsSlotsExist(ctx context.Context, f *dto.SlotFilter) (bool, error) {
	startOfDay := time.Date(f.Date.Year(), f.Date.Month(), f.Date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.AddDate(0, 0, 1)

	query := `
        SELECT EXISTS(
            SELECT 1 FROM slots 
            WHERE start_time >= $1 
              AND start_time < $2
              AND room_id = $3
        )
    `

	var exists bool
	err := r.db.QueryRowContext(ctx, query, startOfDay, endOfDay, f.RoomID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check slots existence: %w", err)
	}

	return exists, nil
}
