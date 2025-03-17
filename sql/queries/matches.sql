-- name: GetMatches :many
SELECT * FROM matches;

-- name: GetMatchesByClub :many
SELECT * FROM matches
WHERE home_team_id = $1 or away_team_id = $1;

-- name: CreateMatch :one
INSERT INTO matches (id, competition_id, home_team_id, away_team_id,
home_goals, away_goals, date, kick_off_time, referee_id, venue_id, attendance, home_xg, away_xg, weekday, url)
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

-- name: GetGamesByTeamAndSeason :many
SELECT 
HT.name as home_team, 
AT.name as away_team, 
M.home_goals as home_goals, 
M.away_goals as away_goals,
M.date as date,
venues.name as stadium
FROM matches as M
INNER JOIN teams as HT on HT.id = M.home_team_id
INNER JOIN teams as AT on AT.id = M.away_team_id
INNER JOIN competitions on competitions.id = M.competition_id
INNER JOIN venues on M.venue_id = venues.id
WHERE (HT.name = $1 OR AT.name = $1) AND competitions.name = $2 AND competitions.season = $3;

-- name: DeleteMatches :exec
DELETE FROM matches;