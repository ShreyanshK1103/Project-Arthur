-- name: CreateProject :one
INSERT INTO projects (
    name, 
    user_id, 
    repo_url, 
    install_command, 
    build_command, 
    output_dir,
    github_repo_id,
    branch,
    auto_deploy,
    slug
)
VALUES(
    $1,$2,$3,$4,$5,$6,$7,$8,$9,$10
)
RETURNING *;

-- name: GetProjectByID :one
SELECT * FROM projects
WHERE id = $1;

-- name: GetProjectsByUsers :many
SELECT * FROM projects
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetProjectByGithubRepoID :one
SELECT *
FROM projects
WHERE github_repo_id = $1
AND auto_deploy = true;

-- name: GetProjectByIDAndUser :one
SELECT *
FROM projects
WHERE id = $1
AND user_id = $2;

-- name: GetProjectBySlug :one
SELECT * FROM projects
WHERE slug = $1;

-- name: SetActiveDeployment :exec
UPDATE projects
SET active_deployment_id = $2
WHERE id = $1;

-- name: PromoteDeployment :one
UPDATE projects p
SET active_deployment_id = d.id
FROM deployments d
WHERE p.id = $1
AND p.user_id = $2
AND d.id = $3
AND d.project_id = p.id
AND d.status = 'success'
RETURNING p.id;

-- name: DeleteProjectByIDAndUser :exec
DELETE FROM projects
WHERE id = $1
AND user_id = $2;