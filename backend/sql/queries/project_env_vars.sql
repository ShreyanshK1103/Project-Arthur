-- name: UpsertEnvVar :exec
INSERT INTO project_env_vars (project_id, key, value_encrypted)
VALUES ($1, $2, $3)
ON CONFLICT (project_id, key)
DO UPDATE SET
    value_encrypted = EXCLUDED.value_encrypted,
    updated_at = NOW();

-- name: ListEnvVarsByProject :many
SELECT *
FROM project_env_vars
WHERE project_id = $1
ORDER BY key;

-- name: DeleteEnvVar :exec
DELETE FROM project_env_vars
WHERE project_id = $1
AND key = $2;

-- name: CountEnvVarsByProject :one
SELECT COUNT(*)
FROM project_env_vars
WHERE project_id = $1;