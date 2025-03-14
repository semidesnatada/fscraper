-- name: CreateCompetition :one
INSERT INTO competitions (id, name, season, url)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: CheckIfCompetitionExistsById :one
SELECT EXISTS (SELECT 1
FROM competitions
WHERE id = $1
LIMIT 1);

-- name: CheckIfCompetitionExistsByNameAndSeason :one
SELECT EXISTS (SELECT 1
FROM competitions
WHERE name = $1 AND season = $2
LIMIT 1);

-- name: GetCompetitionIdFromNameAndSeason :one
SELECT id
FROM competitions
WHERE name = $1 AND season = $2;

-- name: GetCompetitionNameAndSeasonFromId :one
SELECT name, season
FROM competitions
WHERE id = $1;

-- name: GetCompetitionTable :many
SELECT 
    teams.name AS team_name,
    SUM(matches.home_goals) AS goals_scored,
    SUM(matches.away_goals) AS goals_conceded,
    COUNT(*) AS games_played,
    SUM(matches.home_goals) - SUM(matches.away_goals) AS goal_difference,
    SUM(CASE WHEN home_goals>away_goals THEN 1 ELSE 0 END) AS wins,
    SUM(CASE WHEN home_goals=away_goals THEN 1 ELSE 0 END) AS draws,
    SUM(CASE WHEN home_goals<away_goals THEN 1 ELSE 0 END) AS losses,
    SUM((CASE WHEN home_goals>away_goals THEN 1 ELSE 0 END)*3 + (CASE WHEN home_goals=away_goals THEN 1 ELSE 0 END)) AS points
FROM teams
INNER JOIN matches ON matches.home_team_id = teams.id
INNER JOIN competitions ON matches.competition_id = competitions.id
WHERE competitions.name = $1 AND competitions.season = $2
GROUP BY team_name
ORDER BY points DESC, goal_difference DESC, goals_scored DESC, goals_conceded DESC, wins DESC;

-- name: DeleteCompetitions :exec
DELETE FROM competitions;