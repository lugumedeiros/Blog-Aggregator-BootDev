-- +goose Up
CREATE TABLE feeds(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL
);

-- +goose Down
DROP TABLE feeds;