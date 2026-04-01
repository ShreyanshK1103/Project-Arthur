-- +goose Up
CREATE TABLE deployments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL,

    status TEXT NOT NULL,
    repo_url TEXT NOT NULL,
    url TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE INDEX idx_deployments_status_created_at
ON deployments (status, created_at);

-- +goose Down
DROP TABLE deployments;