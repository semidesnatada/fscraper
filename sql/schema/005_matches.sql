-- +goose Up
CREATE TABLE matches (
    id UUID PRIMARY KEY,
    competition_id UUID NOT NULL,
    competition_season_id TEXT NOT NULL,
    home_team_id UUID NOT NULL,
    away_team_id UUID NOT NULL,
    home_goals INTEGER NOT NULL,
    away_goals INTEGER NOT NULL,
    date DATE NOT NULL,
    kick_off_time TIME,
    referee_id UUID,
    venue_id UUID,
    attendance INTEGER,
    home_xg REAL,
    away_xg REAL,
    weekday TEXT NOT NULL,
    -- FOREIGN KEY (competition_id) REFERENCES competitions (id),
    FOREIGN KEY (home_team_id) REFERENCES teams (id),
    FOREIGN KEY (away_team_id) REFERENCES teams (id),
    FOREIGN KEY (referee_id) REFERENCES referees (id),
    FOREIGN KEY (venue_id) REFERENCES venues (id)
);

-- +goose Down
DROP TABLE matches;