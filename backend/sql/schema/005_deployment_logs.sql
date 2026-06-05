-- +goose Up
CREATE TABLE deployment_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    deployment_id UUID NOT NULL,

    log TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (deployment_id)
    REFERENCES deployments(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE deployment_logs;

