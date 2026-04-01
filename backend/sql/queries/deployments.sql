-- name: CreateDeployment :one
INSERT INTO deployments (project_id, status, repo_url)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetDeploymentByID :one
SELECT * FROM deployments
WHERE id = $1;

-- name: GetNextDeployment :one
UPDATE deployments
SET status = 'building',
    updated_at = NOW()
WHERE id = (
    SELECT id
    FROM deployments
    WHERE status = 'queued'
    ORDER BY created_at
    LIMIT 1
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