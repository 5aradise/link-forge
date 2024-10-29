-- +goose Up
CREATE TABLE urls (
    id INTEGER PRIMARY KEY,
    alias TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL
);

-- +goose Down
DROP TABLE urls;
