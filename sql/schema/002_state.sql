-- +goose Up
CREATE TABLE state (
    id INTEGER PRIMARY KEY,
    alias_count INTEGER NOT NULL
);
INSERT INTO state (id, alias_count) VALUES (1, 0);

-- +goose Down
DROP TABLE state;