-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, email, hashed_password)
VALUES ( NOW(), NOW(), $1, $2)
RETURNING *;

-- name: UpdateUser :one
UPDATE users 
SET updated_at = NOW(), email = $2, hashed_password = $3
WHERE id = $1
RETURNING *;