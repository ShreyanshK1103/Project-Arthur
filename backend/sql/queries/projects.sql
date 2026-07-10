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
    auto_deploy
)
VALUES(
    $1,$2,$3,$4,$5,$6,$7,$8,$9
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