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
    SUM(league_matches.home_goals) AS goals_scored,
    SUM(league_matches.away_goals) AS goals_conceded,
    COUNT(*) AS games_played,
    SUM(league_matches.home_goals) - SUM(league_matches.away_goals) AS goal_difference,
    SUM(CASE WHEN home_goals>away_goals THEN 1 ELSE 0 END) AS wins,
    SUM(CASE WHEN home_goals=away_goals THEN 1 ELSE 0 END) AS draws,
    SUM(CASE WHEN home_goals<away_goals THEN 1 ELSE 0 END) AS losses,
    SUM((CASE WHEN home_goals>away_goals THEN 1 ELSE 0 END)*3 + (CASE WHEN home_goals=away_goals THEN 1 ELSE 0 END)) AS points
FROM teams
INNER JOIN league_matches ON league_matches.home_team_id = teams.id
INNER JOIN competitions ON league_matches.competition_id = competitions.id
WHERE competitions.name = $1 AND competitions.season = $2
GROUP BY team_name
UNION ALL
SELECT 
    teams.name AS team_name,
    SUM(league_matches.away_goals) AS goals_scored,
    SUM(league_matches.home_goals) AS goals_conceded,
    COUNT(*) AS games_played,
    SUM(league_matches.away_goals) - SUM(league_matches.home_goals) AS goal_difference,
    SUM(CASE WHEN home_goals<away_goals THEN 1 ELSE 0 END) AS wins,
    SUM(CASE WHEN home_goals=away_goals THEN 1 ELSE 0 END) AS draws,
    SUM(CASE WHEN home_goals>away_goals THEN 1 ELSE 0 END) AS losses,
    SUM((CASE WHEN home_goals<away_goals THEN 1 ELSE 0 END)*3 + (CASE WHEN home_goals=away_goals THEN 1 ELSE 0 END)) AS points
FROM teams
INNER JOIN league_matches ON league_matches.away_team_id = teams.id
INNER JOIN competitions ON league_matches.competition_id = competitions.id
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
    SUM(league_matches.home_goals) AS goals_scored,
    SUM(league_matches.away_goals) AS goals_conceded,
    COUNT(*) AS games_played,
    SUM(league_matches.home_goals) - SUM(league_matches.away_goals) AS goal_difference,
    SUM(CASE WHEN home_goals>away_goals THEN 1 ELSE 0 END) AS wins,
    SUM(CASE WHEN home_goals=away_goals THEN 1 ELSE 0 END) AS draws,
    SUM(CASE WHEN home_goals<away_goals THEN 1 ELSE 0 END) AS losses,
    SUM((CASE WHEN home_goals>away_goals THEN 1 ELSE 0 END)*3 + (CASE WHEN home_goals=away_goals THEN 1 ELSE 0 END)) AS points
FROM teams
INNER JOIN league_matches ON league_matches.home_team_id = teams.id
INNER JOIN competitions ON league_matches.competition_id = competitions.id
WHERE competitions.name = $1
GROUP BY team_name
UNION ALL
SELECT 
    teams.name AS team_name,
    SUM(league_matches.away_goals) AS goals_scored,
    SUM(league_matches.home_goals) AS goals_conceded,
    COUNT(*) AS games_played,
    SUM(league_matches.away_goals) - SUM(league_matches.home_goals) AS goal_difference,
    SUM(CASE WHEN home_goals<away_goals THEN 1 ELSE 0 END) AS wins,
    SUM(CASE WHEN home_goals=away_goals THEN 1 ELSE 0 END) AS draws,
    SUM(CASE WHEN home_goals>away_goals THEN 1 ELSE 0 END) AS losses,
    SUM((CASE WHEN home_goals<away_goals THEN 1 ELSE 0 END)*3 + (CASE WHEN home_goals=away_goals THEN 1 ELSE 0 END)) AS points
FROM teams
INNER JOIN league_matches ON league_matches.away_team_id = teams.id
INNER JOIN competitions ON league_matches.competition_id = competitions.id
WHERE competitions.name = $1
GROUP BY team_name
) s
GROUP BY team_name
ORDER BY points DESC, goal_difference DESC, goals_scored DESC, goals_conceded DESC, wins DESC;

-- name: GetAllClubCompetitionResults :many
SELECT 
    team_name,
    competition_name,
    competition_season,
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
    competitions.id AS competition_id,
    competitions.name AS competition_name,
    competitions.season AS competition_season,
    SUM(league_matches.home_goals) AS goals_scored,
    SUM(league_matches.away_goals) AS goals_conceded,
    COUNT(*) AS games_played,
    SUM(league_matches.home_goals) - SUM(league_matches.away_goals) AS goal_difference,
    SUM(CASE WHEN home_goals>away_goals THEN 1 ELSE 0 END) AS wins,
    SUM(CASE WHEN home_goals=away_goals THEN 1 ELSE 0 END) AS draws,
    SUM(CASE WHEN home_goals<away_goals THEN 1 ELSE 0 END) AS losses,
    SUM((CASE WHEN home_goals>away_goals THEN 1 ELSE 0 END)*3 + (CASE WHEN home_goals=away_goals THEN 1 ELSE 0 END)) AS points
FROM teams
INNER JOIN league_matches ON league_matches.home_team_id = teams.id
INNER JOIN competitions ON league_matches.competition_id = competitions.id
WHERE teams.name = $1
GROUP BY competitions.id, teams.name, competitions.name, competitions.season
UNION ALL
SELECT 
    teams.name AS team_name,
    competitions.id AS competition_id,
    competitions.name AS competition_name,
    competitions.season AS competition_season,
    SUM(league_matches.away_goals) AS goals_scored,
    SUM(league_matches.home_goals) AS goals_conceded,
    COUNT(*) AS games_played,
    SUM(league_matches.away_goals) - SUM(league_matches.home_goals) AS goal_difference,
    SUM(CASE WHEN home_goals<away_goals THEN 1 ELSE 0 END) AS wins,
    SUM(CASE WHEN home_goals=away_goals THEN 1 ELSE 0 END) AS draws,
    SUM(CASE WHEN home_goals>away_goals THEN 1 ELSE 0 END) AS losses,
    SUM((CASE WHEN home_goals<away_goals THEN 1 ELSE 0 END)*3 + (CASE WHEN home_goals=away_goals THEN 1 ELSE 0 END)) AS points
FROM teams
INNER JOIN league_matches ON league_matches.away_team_id = teams.id
INNER JOIN competitions ON league_matches.competition_id = competitions.id
WHERE teams.name = $1
GROUP BY competitions.id, teams.name, competitions.name, competitions.season
) s
GROUP BY competition_id, team_name, competition_name, competition_season
ORDER BY points DESC, goal_difference DESC, goals_scored DESC, goals_conceded DESC, wins DESC;


-- name: DeleteCompetitions :exec
DELETE FROM competitions;