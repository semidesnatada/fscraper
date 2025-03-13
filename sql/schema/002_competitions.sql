-- +goose Up
CREATE TABLE competitions (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE competitions;