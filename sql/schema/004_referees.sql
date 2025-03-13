-- +goose Up
CREATE TABLE referees (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE IF EXISTS referees;