-- +goose Up
CREATE TABLE competitions (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    season TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    UNIQUE (name, season)
);

-- +goose Down
DROP TABLE competitions;