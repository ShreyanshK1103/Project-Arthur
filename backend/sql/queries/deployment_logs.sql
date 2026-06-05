-- name: CreateDeploymentLog :exec
INSERT INTO deployment_logs (
    deployment_id,
    log
)
VALUES ($1, $2);

-- name: GetDeploymentLogs :many
SELECT *
FROM deployment_logs
WHERE deployment_id = $1
ORDER BY created_at ASC;