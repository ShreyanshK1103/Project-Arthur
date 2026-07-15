-- name: CreateUser :one
INSERT INTO users (
    name,
    email,
    password_hash,
    avatar_url,
    email_verified
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;
-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;


-- name: UpdateUser :one
UPDATE users
SET
    name = $2,
    avatar_url = $3,
    email_verified = $4
WHERE id = $1
RETURNING *;