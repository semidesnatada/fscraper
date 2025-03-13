-- +goose Up
CREATE TABLE teams (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE teams;

