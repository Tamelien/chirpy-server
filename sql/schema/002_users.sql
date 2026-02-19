-- +goose Up
DROP TABLE users;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT UNIQUE NOT NULL
);
-- +goose Down
DROP TABLE users;

CREATE TABLE users(
    id INTEGER PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT UNIQUE NOT NULL
);