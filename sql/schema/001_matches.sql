-- +goose Up
CREATE TABLE matches (
    id UUID PRIMARY KEY,
    home_team TEXT NOT NULL,
    away_team TEXT NOT NULL,
    home_goals INTEGER NOT NULL,
    away_goals INTEGER NOT NULL
);

-- +goose Down
DROP TABLE matches;