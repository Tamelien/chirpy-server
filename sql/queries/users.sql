-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, email, hashed_password)
VALUES ( NOW(), NOW(), $1, $2)
RETURNING *;

-- name: UpdateUser :one
UPDATE users 
SET updated_at = NOW(), email = $2, hashed_password = $3
WHERE id = $1
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserFromRefreshToken :one
SELECT * FROM users
WHERE id IN (
    SELECT user_id FROM refresh_tokens
    WHERE (token = $1 AND revoked_at IS NULL AND expires_at > NOW())
);