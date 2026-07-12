-- +goose Up

ALTER TABLE users
ADD COLUMN password_hash TEXT,
ADD COLUMN provider TEXT NOT NULL DEFAULT 'local',
ADD COLUMN provider_id TEXT,
ADD COLUMN avatar_url TEXT,
ADD COLUMN email_verified BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down

ALTER TABLE users
DROP COLUMN password_hash,
DROP COLUMN provider,
DROP COLUMN provider_id,
DROP COLUMN avatar_url,
DROP COLUMN email_verified;