-- name: CreateOAuthAccount :one
INSERT INTO oauth_accounts(
    user_id,
    provider,
    provider_id
)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetOAuthAccount :one
SELECT *
FROM oauth_accounts
WHERE provider = $1
AND provider_id = $2;

-- name: GetUserByOAuth :one
SELECT users.*
FROM users
JOIN oauth_accounts
ON users.id = oauth_accounts.user_id
WHERE oauth_accounts.provider = $1
AND oauth_accounts.provider_id = $2;

