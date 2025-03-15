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
    team_name,
    SUM(goals_scored) AS goals_scored,
    SUM(goals_conceded) AS goals_conceded,
    SUM(goal_difference) AS goal_difference,
    SUM(games_played) AS games_played,
    SUM(wins) AS wins,
    SUM(draws) AS draws,
    SUM(losses) AS losses,
    SUM(points) AS points
FROM(
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
UNION ALL
SELECT 
    teams.name AS team_name,
    SUM(matches.away_goals) AS goals_scored,
    SUM(matches.home_goals) AS goals_conceded,
    COUNT(*) AS games_played,
    SUM(matches.away_goals) - SUM(matches.home_goals) AS goal_difference,
    SUM(CASE WHEN home_goals<away_goals THEN 1 ELSE 0 END) AS wins,
    SUM(CASE WHEN home_goals=away_goals THEN 1 ELSE 0 END) AS draws,
    SUM(CASE WHEN home_goals>away_goals THEN 1 ELSE 0 END) AS losses,
    SUM((CASE WHEN home_goals<away_goals THEN 1 ELSE 0 END)*3 + (CASE WHEN home_goals=away_goals THEN 1 ELSE 0 END)) AS points
FROM teams
INNER JOIN matches ON matches.away_team_id = teams.id
INNER JOIN competitions ON matches.competition_id = competitions.id
WHERE competitions.name = $1 AND competitions.season = $2
GROUP BY team_name
) s
GROUP BY team_name
ORDER BY points DESC, goal_difference DESC, goals_scored DESC, goals_conceded DESC, wins DESC;

-- name: GetUniqueCompetitionSeasons :many
SELECT name, season
FROM competitions;

-- name: GetAllTimeCompetitionTable :many
SELECT 
    team_name,
    SUM(goals_scored) AS goals_scored,
    SUM(goals_conceded) AS goals_conceded,
    SUM(goal_difference) AS goal_difference,
    SUM(games_played) AS games_played,
    SUM(wins) AS wins,
    SUM(draws) AS draws,
    SUM(losses) AS losses,
    SUM(points) AS points
FROM(
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
WHERE competitions.name = $1
GROUP BY team_name
UNION ALL
SELECT 
    teams.name AS team_name,
    SUM(matches.away_goals) AS goals_scored,
    SUM(matches.home_goals) AS goals_conceded,
    COUNT(*) AS games_played,
    SUM(matches.away_goals) - SUM(matches.home_goals) AS goal_difference,
    SUM(CASE WHEN home_goals<away_goals THEN 1 ELSE 0 END) AS wins,
    SUM(CASE WHEN home_goals=away_goals THEN 1 ELSE 0 END) AS draws,
    SUM(CASE WHEN home_goals>away_goals THEN 1 ELSE 0 END) AS losses,
    SUM((CASE WHEN home_goals<away_goals THEN 1 ELSE 0 END)*3 + (CASE WHEN home_goals=away_goals THEN 1 ELSE 0 END)) AS points
FROM teams
INNER JOIN matches ON matches.away_team_id = teams.id
INNER JOIN competitions ON matches.competition_id = competitions.id
WHERE competitions.name = $1
GROUP BY team_name
) s
GROUP BY team_name
ORDER BY points DESC, goal_difference DESC, goals_scored DESC, goals_conceded DESC, wins DESC;

-- name: DeleteCompetitions :exec
DELETE FROM competitions;