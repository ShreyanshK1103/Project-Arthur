-- +goose Up
ALTER TABLE deployment_logs ADD COLUMN seq BIGSERIAL NOT NULL;
CREATE INDEX deployment_logs_deployment_seq_idx ON deployment_logs (deployment_id, seq);

-- +goose Down
DROP INDEX deployment_logs_deployment_seq_idx;
ALTER TABLE deployment_logs DROP COLUMN seq;