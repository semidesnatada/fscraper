-- +goose Up
CREATE TABLE league_matches (
    id UUID PRIMARY KEY,
    competition_id UUID NOT NULL,
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
    url TEXT NOT NULL UNIQUE,
    FOREIGN KEY (competition_id) REFERENCES competitions (id) ON DELETE CASCADE,
    FOREIGN KEY (home_team_id) REFERENCES teams (id) ON DELETE CASCADE,
    FOREIGN KEY (away_team_id) REFERENCES teams (id) ON DELETE CASCADE,
    FOREIGN KEY (referee_id) REFERENCES referees (id) ON DELETE CASCADE,
    FOREIGN KEY (venue_id) REFERENCES venues (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE league_matches;