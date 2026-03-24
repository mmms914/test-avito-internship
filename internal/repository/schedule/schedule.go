package schedule

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
	"github.com/avito-internships/test-backend-1-mmms914/internal/repository/converter"
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

func (r *Repository) Create(ctx context.Context, schedule *domain.Schedule) error {
	query := `
        INSERT INTO schedules (id, room_id, days_of_week, start_time, end_time)
        VALUES ($1, $2, $3, $4, $5)
    `

	daysOfWeek := converter.WeekdaysToIntArray(schedule.DaysOfWeek())

	_, err := r.db.ExecContext(ctx, query,
		schedule.ID(),
		schedule.RoomID(),
		daysOfWeek,
		schedule.StartTime(),
		schedule.EndTime(),
	)

	if err != nil {
		return fmt.Errorf("failed to create schedule: %w", err)
	}

	return nil
}

func (r *Repository) GetForRoom(ctx context.Context, roomID uuid.UUID) (*domain.Schedule, error) {
	query := `
        SELECT id, room_id, days_of_week, start_time, end_time
        FROM schedules
        WHERE room_id = $1
    `

	var specs domain.ScheduleRestoreSpecs
	var daysOfWeek []int

	err := r.db.QueryRowContext(ctx, query, roomID).Scan(
		&specs.ID,
		&specs.RoomID,
		&daysOfWeek,
		&specs.StartTime,
		&specs.EndTime,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrScheduleNotExists
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find schedule by room id: %w", err)
	}

	days, err := converter.IntArrayToWeekdays(&daysOfWeek)
	if err != nil {
		return nil, fmt.Errorf("failed to convert days of week: %w", err)
	}

	specs.DaysOfWeek = days

	return domain.NewSchedule(domain.WithScheduleRestoreSpecs(&specs)), nil
}

func (r *Repository) ExistsForRoom(ctx context.Context, roomID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM schedules WHERE room_id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, roomID).Scan(&exists)
	return exists, err
}
