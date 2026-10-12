-- +goose Up
CREATE TABLE project_env_vars (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    key TEXT NOT NULL,
    value_encrypted BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (project_id, key)
);

-- +goose Down
DROP TABLE project_env_vars;