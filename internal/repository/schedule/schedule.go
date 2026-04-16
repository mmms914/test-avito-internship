package schedule

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/mmms914/test-avito-internship/internal/domain"
	"github.com/mmms914/test-avito-internship/internal/errs"
	"github.com/mmms914/test-avito-internship/internal/repository/converter"
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
		converter.DurationToStringTime(schedule.StartTime()),
		converter.DurationToStringTime(schedule.EndTime()),
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
	var daysOfWeekStr, startTimeStr, endTimeStr string

	err := r.db.QueryRowContext(ctx, query, roomID).Scan(
		&specs.ID,
		&specs.RoomID,
		&daysOfWeekStr,
		&startTimeStr,
		&endTimeStr,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrScheduleNotExists
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find schedule by room id: %w", err)
	}

	daysOfWeek, err := converter.StringToIntArray(daysOfWeekStr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert postgres format of array: %w", err)
	}

	days, err := converter.IntArrayToWeekdays(&daysOfWeek)
	if err != nil {
		return nil, fmt.Errorf("failed to convert days of week: %w", err)
	}

	specs.DaysOfWeek = days

	startTime, err := converter.StringToDuration(startTimeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert postgres format of time to duration: %w", err)
	}

	endTime, err := converter.StringToDuration(endTimeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert postgres format of time to duration: %w", err)
	}

	specs.StartTime = startTime
	specs.EndTime = endTime

	return domain.NewSchedule(domain.WithScheduleRestoreSpecs(&specs)), nil
}

func (r *Repository) ExistsForRoom(ctx context.Context, roomID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM schedules WHERE room_id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, roomID).Scan(&exists)
	return exists, err
}
