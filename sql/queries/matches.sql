-- name: GetMatches :many
SELECT * FROM matches;

-- name: GetMatchesByClub :many
SELECT * FROM matches
WHERE home_team_id = $1 or away_team_id = $1;

-- name: CreateMatch :one
INSERT INTO matches (id, competition_id, competition_season_id, home_team_id, away_team_id,
home_goals, away_goals, date, kick_off_time, referee_id, venue_id, attendance, home_xg, away_xg, weekday)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13,
    $14,
    $15
)
RETURNING *;

-- name: DeleteMatches :exec
DELETE FROM matches;