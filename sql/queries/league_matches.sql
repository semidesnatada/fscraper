-- name: GetLeagueMatches :many
SELECT * FROM league_matches;

-- name: GetLeagueMatchesByClub :many
SELECT * FROM league_matches
WHERE home_team_id = $1 or away_team_id = $1;

-- name: CreateLeagueMatch :one
INSERT INTO league_matches (id, competition_id, home_team_id, away_team_id,
home_goals, away_goals, date, kick_off_time, referee_id, venue_id, attendance,
home_xg, away_xg, weekday, url, home_team_online_id, away_team_online_id)
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
    $15,
    $16,
    $17
)
RETURNING *;

-- name: GetLeagueGamesByTeamAndSeason :many
SELECT 
HT.name as home_team, 
AT.name as away_team, 
M.home_goals as home_goals, 
M.away_goals as away_goals,
M.date as date,
venues.name as stadium
FROM league_matches as M
INNER JOIN teams as HT on HT.id = M.home_team_id
INNER JOIN teams as AT on AT.id = M.away_team_id
INNER JOIN competitions on competitions.id = M.competition_id
INNER JOIN venues on M.venue_id = venues.id
WHERE (HT.name = $1 OR AT.name = $1) AND competitions.name = $2 AND competitions.season = $3;

-- name: GetLeagueMatchUrls :many
SELECT url
FROM league_matches;

-- name: GetLeagueMatchUrlsAndTeamOnlineIds :many
SELECT url, home_team_online_id, away_team_online_id
FROM league_matches;

-- name: GetLeagueMatchUrlsAndTeamOnlineIdsWOffset :many
SELECT url, home_team_online_id, away_team_online_id
FROM league_matches
ORDER BY url
LIMIT 1000
OFFSET $1;

-- name: GetLeagueMatchIDFromUrl :one
SELECT id
FROM league_matches
WHERE url = $1;

-- name: DeleteLeagueMatches :exec
DELETE FROM league_matches;