-- +goose Up
CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    name TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;