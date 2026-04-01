-- name: CreateProject :one
INSERT INTO projects (name, user_id)
VALUES($1, $2)
RETURNING *;

-- name: GetProjectByID :one
SELECT * FROM projects
WHERE id = $1;

-- name: GetProjectsByUsers :many
SELECT * FROM projects
WHERE user_id = $1
ORDER BY created_at DESC;