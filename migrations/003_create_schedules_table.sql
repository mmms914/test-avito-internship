-- +goose Up
CREATE TABLE IF NOT EXISTS schedules (
    id UUID PRIMARY KEY,
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    days_of_week INT[] NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    UNIQUE(room_id)
);

CREATE INDEX idx_schedules_room_id ON schedules(room_id);

-- +goose Down
DROP TABLE IF EXISTS schedules;