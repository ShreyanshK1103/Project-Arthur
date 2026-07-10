-- +goose Up

ALTER TABLE projects
ADD COLUMN repo_url TEXT NOT NULL DEFAULT '',
ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT now();

ALTER TABLE deployments
DROP COLUMN repo_url;


-- +goose Down

ALTER TABLE deployments
ADD COLUMN repo_url TEXT NOT NULL DEFAULT '';

ALTER TABLE projects
DROP COLUMN repo_url,
DROP COLUMN updated_at;