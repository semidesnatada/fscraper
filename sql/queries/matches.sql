-- name: GetMatches :many
SELECT * FROM matches;

-- name: CreateMatch :one
INSERT INTO matches (id, home_team, away_team, home_goals, away_goals)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;