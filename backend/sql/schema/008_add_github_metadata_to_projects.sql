-- +goose Up

ALTER TABLE projects
ADD COLUMN github_repo_id BIGINT,
ADD COLUMN branch TEXT NOT NULL DEFAULT 'main',
ADD COLUMN auto_deploy BOOLEAN NOT NULL DEFAULT true;


-- +goose Down

ALTER TABLE projects
DROP COLUMN github_repo_id,
DROP COLUMN branch,
DROP COLUMN auto_deploy;