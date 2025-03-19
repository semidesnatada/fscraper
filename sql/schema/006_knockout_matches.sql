-- +goose Up
CREATE TABLE knockout_matches (
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
    went_to_pens BOOLEAN NOT NULL,
    home_pens INTEGER,
    away_pens INTEGER,
    round TEXT NOT NULL,
    weekday TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    home_team_online_id TEXT NOT NULL,
    away_team_online_id TEXT NOT NULL,
    FOREIGN KEY (competition_id) REFERENCES competitions (id) ON DELETE CASCADE,
    FOREIGN KEY (home_team_id) REFERENCES teams (id) ON DELETE CASCADE,
    FOREIGN KEY (away_team_id) REFERENCES teams (id) ON DELETE CASCADE,
    FOREIGN KEY (referee_id) REFERENCES referees (id) ON DELETE CASCADE,
    FOREIGN KEY (venue_id) REFERENCES venues (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE knockout_matches;