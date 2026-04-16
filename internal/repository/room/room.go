package room

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/mmms914/test-avito-internship/internal/domain"
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

func (r *Repository) Create(ctx context.Context, room *domain.Room) error {
	query := `
        INSERT INTO rooms (id, name, description, capacity, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `

	_, err := r.db.ExecContext(ctx, query,
		room.ID(),
		room.Name(),
		room.Description(),
		room.Capacity(),
		room.CreatedAt(),
	)

	return err
}

func (r *Repository) GetAll(ctx context.Context) ([]*domain.Room, error) {
	query := `
        SELECT id, name, description, capacity, created_at
        FROM rooms
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list rooms: %w", err)
	}
	defer rows.Close()

	var rooms []*domain.Room
	for rows.Next() {
		var specs domain.RoomRestoreSpecs
		if err = rows.Scan(
			&specs.ID,
			&specs.Name,
			&specs.Description,
			&specs.Capacity,
			&specs.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan room: %w", err)
		}
		rooms = append(rooms, domain.NewRoom(domain.WithRoomRestoreSpecs(&specs)))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return rooms, nil
}

func (r *Repository) Exists(ctx context.Context, roomID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM rooms WHERE id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, roomID).Scan(&exists)
	return exists, err
}
