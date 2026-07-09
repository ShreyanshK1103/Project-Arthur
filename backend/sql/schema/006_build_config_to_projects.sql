-- +goose Up

ALTER TABLE projects
ADD COLUMN install_command TEXT NOT NULL DEFAULT 'npm install',
ADD COLUMN build_command TEXT NOT NULL DEFAULT 'npm run build',
ADD COLUMN output_dir TEXT NOT NULL DEFAULT 'dist';


-- +goose Down

ALTER TABLE projects
DROP COLUMN install_command,
DROP COLUMN build_command,
DROP COLUMN output_dir;