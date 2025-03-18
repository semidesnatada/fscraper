-- name: GetKnockoutMatches :many
SELECT * FROM knockout_matches;

-- name: GetKnockoutMatchesByClub :many
SELECT * FROM knockout_matches
WHERE home_team_id = $1 or away_team_id = $1;

-- name: CreateKnockoutMatch :one
INSERT INTO knockout_matches (id, competition_id, home_team_id, away_team_id,
home_goals, away_goals, date, kick_off_time, referee_id, venue_id, attendance, 
home_xg, away_xg, went_to_pens, home_pens, away_pens, round, weekday, url)
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
    $17,
    $18,
    $19
)
RETURNING *;

-- name: GetKnockoutGamesByTeamAndSeason :many
SELECT 
HT.name as home_team, 
AT.name as away_team, 
M.home_goals as home_goals, 
M.away_goals as away_goals,
M.date as date,
venues.name as stadium
FROM knockout_matches as M
INNER JOIN teams as HT on HT.id = M.home_team_id
INNER JOIN teams as AT on AT.id = M.away_team_id
INNER JOIN competitions on competitions.id = M.competition_id
INNER JOIN venues on M.venue_id = venues.id
WHERE (HT.name = $1 OR AT.name = $1) AND competitions.name = $2 AND competitions.season = $3;

-- name: GetKnockoutGamesByRoundAndSeason :many
SELECT
    M.home_goals AS home_goals,
    M.away_goals AS away_goals,
    M.went_to_pens AS went_to_pens,
    M.home_pens AS home_pens,
    M.away_pens AS away_pens,
    M.date AS date,
    M.kick_off_time AS kick_off_time,
    M.attendance AS attendance,
    venues.name AS stadium,
    referees.name AS referee,
    HT.name AS home_team,
    AT.name AS away_team
FROM knockout_matches as M
INNER JOIN teams as HT on HT.id = M.home_team_id
INNER JOIN teams as AT on AT.id = M.away_team_id
INNER JOIN venues on M.venue_id = venues.id
INNER JOIN referees on M.referee_id = referees.id
INNER JOIN competitions ON M.competition_id = competitions.id
WHERE round = $1 AND competitions.name = $2 AND competitions.season = $3;

-- name: GetMatchesInEachRoundForGivenComp :many
SELECT 
    round,
    COUNT(*) AS no_of_games
FROM knockout_matches
INNER JOIN competitions ON knockout_matches.competition_id = competitions.id
WHERE competitions.name = $1 AND competitions.season = $2
GROUP BY round
ORDER BY no_of_games;

-- name: GetAllMatchesFromComp :many
SELECT * FROM knockout_matches
INNER JOIN competitions on competitions.id = knockout_matches.competition_id
WHERE competitions.name = $1 AND competitions.season = $2 AND round = $3;

-- name: DeleteKnockoutMatches :exec
DELETE FROM knockout_matches;