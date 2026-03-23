-- +goose Up
CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY,
    slot_id UUID NOT NULL REFERENCES slots(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL,
    conference_link TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_bookings_slot_id ON bookings(slot_id);
CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_status ON bookings(status);

-- +goose Down
DROP TABLE IF EXISTS bookings;