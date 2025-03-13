-- +goose Up
CREATE TABLE venues (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE venues;