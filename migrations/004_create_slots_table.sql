-- +goose Up
CREATE TABLE IF NOT EXISTS slots (
    id UUID PRIMARY KEY,
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_slots_room_id ON slots(room_id);
CREATE INDEX idx_slots_start_time ON slots(start_time);
CREATE INDEX idx_slots_start_time ON slots(end_time);

-- +goose Down
DROP TABLE IF EXISTS slots;