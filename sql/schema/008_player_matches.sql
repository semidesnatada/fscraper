-- +goose Up
CREATE TABLE player_matches (
    match_id UUID NOT NULL,
    player_id UUID NOT NULL,
    match_url TEXT NOT NULL,
    first_minute INTEGER NOT NULL,
    last_minute INTEGER NOT NULL,
    goals INTEGER NOT NULL,
    penalties INTEGER NOT NULL,
    yellow_card INTEGER NOT NULL,
    red_card INTEGER NOT NULL,
    own_goals INTEGER NOT NULL,
    is_knockout BOOLEAN NOT NULL,
    at_home BOOLEAN NOT NULL,
    UNIQUE(match_id, player_id)
);

-- +goose Down
DROP TABLE player_matches;