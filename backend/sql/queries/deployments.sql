-- name: CreateDeployment :one
INSERT INTO deployments (project_id, status)
VALUES ($1, $2)
RETURNING *;

-- name: GetDeploymentByID :one
SELECT * FROM deployments
WHERE id = $1;

-- name: GetNextDeployment :one
WITH next_job AS (
    SELECT id
    FROM deployments
    WHERE status = 'queued'
    ORDER BY created_at
    LIMIT 1
    FOR UPDATE SKIP LOCKED
)
UPDATE deployments
SET status = 'building',
    updated_at = NOW()
WHERE id = (
    SELECT id FROM next_job
)
RETURNING *;

-- name: MarkDeploymentSuccess :exec
UPDATE deployments
SET STATUS = 'success',
    url = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: MarkDeploymentFailed :exec
UPDATE deployments
SET status = 'failed',
    updated_at = NOW()
WHERE id = $1;

-- name: GetDeploymentByPrefix :one
SELECT *
FROM deployments
WHERE id::text LIKE $1 || '%'
LIMIT 1;

-- name: ResetBuildingDeployments :exec
UPDATE deployments
SET status = 'queued',
    updated_at = NOW()
WHERE status = 'building'
AND updated_at < NOW() - INTERVAL '5 minutes';

-- name: GetDeploymentsByProject :many
SELECT *
FROM deployments
WHERE project_id = $1
ORDER BY created_at DESC;