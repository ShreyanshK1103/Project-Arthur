-- +goose Up

ALTER TABLE users
DROP COLUMN provider,
DROP COLUMN provider_id;

-- +goose Down

ALTER TABLE users
ADD COLUMN provider TEXT NOT NULL DEFAULT 'local',
ADD COLUMN provider_id TEXT;