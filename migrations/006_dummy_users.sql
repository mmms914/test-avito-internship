-- +goose Up
INSERT INTO users (id, email, password, role, created_at)
VALUES
    ('849f574c-80c9-4b1e-9862-93bb48848b32', 'admin@example.com', '$2a$10$', 'admin', NOW()),
    ('25a8330f-eb80-41b7-ae55-e3e7350517e7', 'user@example.com', '$2a$10$', 'user', NOW())
    ON CONFLICT (id) DO NOTHING;

-- +goose Down
DELETE FROM users WHERE id IN ('849f574c-80c9-4b1e-9862-93bb48848b32', '25a8330f-eb80-41b7-ae55-e3e7350517e7');