-- +goose Up
CREATE TABLE players (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    nationality TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE players;