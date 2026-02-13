-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    UUID(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;
