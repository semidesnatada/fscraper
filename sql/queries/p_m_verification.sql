-- name: GetNumberOfDistinctPlayersFieldedInLeagueByTeam :many
SELECT
    team_name,
    competition_name,
    COUNT(DISTINCT match_id) as matches_played,
    COUNT(DISTINCT player_id) as players_fielded
FROM(
SELECT
    PM.match_id as match_id,
    PM.player_id as player_id,
    HT.name as team_name,
    C.name as competition_name
FROM player_matches AS PM
INNER JOIN league_matches AS LM ON PM.match_id = LM.id
INNER JOIN teams AS HT ON HT.id = LM.home_team_id
INNER JOIN competitions AS C ON C.id = LM.competition_id
WHERE PM.at_home IS TRUE
GROUP BY team_name, competition_name, match_id, player_id
UNION ALL
SELECT
    PM.match_id as player_id,
    PM.player_id as player_id,
    AT.name as team_name,
    C.name as competition_name
FROM player_matches AS PM
INNER JOIN league_matches AS LM ON PM.match_id = LM.id
INNER JOIN teams AS AT ON AT.id = LM.away_team_id
INNER JOIN competitions AS C ON C.id = LM.competition_id
WHERE PM.at_home IS FALSE
GROUP BY team_name, competition_name, match_id, player_id
) s
INNER JOIN players ON player_id = players.id
WHERE NOT players.name = 'fakeRedCard'
GROUP BY team_name, competition_name
ORDER BY players_fielded DESC, matches_played DESC;

-- name: GetNumberOfGoalsScoredInEachLeagueSeasonByTeam :many
SELECT
    team_name,
    competition_name,
    competition_season,
    SUM(goals_scored) as total_goals_scored,
    COUNT(DISTINCT match_id) as matches_played,
    COUNT(DISTINCT player_id) as players_fielded
FROM(
SELECT
    PM.match_id as match_id,
    PM.player_id as player_id,
    SUM(PM.goals) as goals_scored,
    HT.name as team_name,
    C.name as competition_name,
    C.season as competition_season
FROM player_matches AS PM
INNER JOIN league_matches AS LM ON PM.match_id = LM.id
INNER JOIN teams AS HT ON HT.id = LM.home_team_id
INNER JOIN competitions AS C ON C.id = LM.competition_id
WHERE PM.at_home IS TRUE
GROUP BY team_name, competition_name, competition_season, match_id, player_id
UNION ALL
SELECT
    PM.match_id as match_id,
    PM.player_id as player_id,
    SUM(PM.goals) as goals_scored,
    AT.name as team_name,
    C.name as competition_name,
    C.season as competition_season
FROM player_matches AS PM
INNER JOIN league_matches AS LM ON PM.match_id = LM.id
INNER JOIN teams AS AT ON AT.id = LM.away_team_id
INNER JOIN competitions AS C ON C.id = LM.competition_id
WHERE PM.at_home IS FALSE
GROUP BY team_name, competition_name, competition_season, match_id, player_id
) s
INNER JOIN players ON player_id = players.id
WHERE NOT players.name = 'fakeRedCard'
GROUP BY team_name, competition_name, competition_season
ORDER BY total_goals_scored DESC;

-- name: GetPlayerLeagueStats :many
SELECT
    players.name as player_name,
    SUM(player_matches.goals) as total_goals,
    COUNT(*) as matches_played,
    competitions.name as competition_name,
    competitions.season as competition_season
FROM player_matches
INNER JOIN players ON player_matches.player_id = players.id
INNER JOIN league_matches ON league_matches.url = player_matches.match_url
INNER JOIN competitions ON league_matches.competition_id = competitions.id
GROUP BY player_name, players.id, competition_name, competition_season
ORDER BY total_goals DESC
LIMIT 50;

-- name: GetPlayersPlayedWithByUrl :many
SELECT
    P.name AS colleague_name,
    P.url AS player_url,
    P.id AS player_id,
    P.nationality AS colleague_nationality,
    SUM(CASE WHEN LEAST(OTHERS.last_minute, t_last_min) - GREATEST(OTHERS.first_minute, t_first_min)> 0 THEN LEAST(OTHERS.last_minute, t_last_min) - GREATEST(OTHERS.first_minute, t_first_min) ELSE 0 END) as total_mins_played
FROM (
SELECT
    match_id,
    at_home AS target_at_home,
    first_minute AS t_first_min,
    last_minute AS t_last_min
FROM player_matches
INNER JOIN players ON players.id = player_matches.player_id
WHERE players.url = $1
) AS GAMES_IN_SCOPE
INNER JOIN player_matches AS OTHERS ON (OTHERS.match_id = GAMES_IN_SCOPE.match_id AND OTHERS.at_home = GAMES_IN_SCOPE.target_at_home)
INNER JOIN players AS P ON P.id = OTHERS.player_id 
WHERE NOT P.name = 'fakeRedCard'
GROUP BY P.name, P.url, P.id
ORDER BY total_mins_played DESC;

-- name: GetStatsForPlayerUrl :one
SELECT
    players.name as player_name,
    SUM(last_minute - first_minute) AS total_mins_played,
    COUNT(*) AS matches_played,
    SUM(goals) AS total_goals,
    SUM(penalties) AS total_pens,
    SUM(own_goals) AS total_ogs,
    SUM(yellow_card) AS total_yellow_card,
    SUM(red_card) AS total_red_card
FROM player_matches
INNER JOIN players ON players.id = player_matches.player_id
WHERE players.url = $1
GROUP BY player_name, players.id;

-- name: GetLeagueAllTimeTopScorers :many
SELECT
    players.name as player_name,
    SUM(player_matches.goals) as total_goals,
    COUNT(*) as matches_played,
    competitions.name as competition_name
FROM player_matches
INNER JOIN players ON player_matches.player_id = players.id
INNER JOIN league_matches ON league_matches.url = player_matches.match_url
INNER JOIN competitions ON league_matches.competition_id = competitions.id
GROUP BY player_name, players.id, competition_name
ORDER BY total_goals DESC
LIMIT 50;

-- name: GetAllTimeTopScorers :many
SELECT
    players.name as player_name,
    SUM(player_matches.goals) as total_goals,
    COUNT(*) as matches_played
FROM player_matches
INNER JOIN players ON player_matches.player_id = players.id
INNER JOIN league_matches ON league_matches.url = player_matches.match_url
INNER JOIN competitions ON league_matches.competition_id = competitions.id
GROUP BY players.id, player_name
ORDER BY total_goals DESC
LIMIT 50;

-- name: GetMatchesWhereMinsDontAddUp :many
SELECT
    match_url,
    SUM(last_minute - first_minute) AS squad_mins
FROM player_matches
GROUP BY match_id, match_url, at_home
HAVING SUM(last_minute - first_minute) < 990 AND SUM(red_card) < 1
ORDER BY squad_mins DESC;

-- name: GetPlayerRecordsForGivenLeagueMatch :many
SELECT
    players.name as player_name,
    player_matches.first_minute as first_minute,
    player_matches.last_minute as last_minute,
    CASE WHEN player_matches.at_home IS TRUE THEN HOME_T.name ELSE AWAY_T.name END AS team_name,
    player_matches.red_card as red_card,
    player_matches.yellow_card as yellow_card,
    player_matches.goals as goals,
    player_matches.penalties as pens,
    player_matches.own_goals as ogs
FROM player_matches
INNER JOIN players ON player_matches.player_id = players.id
INNER JOIN league_matches ON player_matches.match_id = league_matches.id
INNER JOIN teams AS HOME_T ON league_matches.home_team_id = HOME_T.id
INNER JOIN teams AS AWAY_T ON league_matches.away_team_id = AWAY_T.id
WHERE player_matches.match_url = $1;

-- name: GetUrlsToRescrape :many
SELECT LM.url, LM.home_team_online_id, LM.away_team_online_id
FROM league_matches AS LM
INNER JOIN player_matches AS PM ON PM.match_id = LM.id
GROUP BY PM.match_id, LM.url, LM.home_team_online_id, LM.away_team_online_id
HAVING SUM(PM.last_minute - PM.first_minute) < 1980
ORDER BY LM.url;

-- name: GetAllPlayersAndSharedMinsByID :many
SELECT
    P.id AS other_player_id,
    SUM(CASE WHEN LEAST(OTHERS.last_minute, t_last_min) - GREATEST(OTHERS.first_minute, t_first_min)> 0 THEN LEAST(OTHERS.last_minute, t_last_min) - GREATEST(OTHERS.first_minute, t_first_min) ELSE 0 END) as shared_minutes
FROM (
SELECT
    match_id,
    at_home AS target_at_home,
    first_minute AS t_first_min,
    last_minute AS t_last_min
FROM player_matches
INNER JOIN players ON players.id = player_matches.player_id
WHERE players.id = $1
) AS GAMES_IN_SCOPE
INNER JOIN player_matches AS OTHERS ON (OTHERS.match_id = GAMES_IN_SCOPE.match_id AND OTHERS.at_home = GAMES_IN_SCOPE.target_at_home)
RIGHT JOIN players AS P ON P.id = OTHERS.player_id
WHERE NOT P.name = 'fakeRedCard'
GROUP BY P.name, P.url, P.id
HAVING SUM(LEAST(OTHERS.last_minute, t_last_min) - GREATEST(OTHERS.first_minute, t_first_min))> 0
ORDER BY shared_minutes DESC;
-- LIMIT 3500;

-- name: GetSharedLeagueStatsForTwoPlayersByIDs :many
SELECT
    P.id AS other_player_id,
    SUM(CASE WHEN LEAST(OTHERS.last_minute, t_last_min) - GREATEST(OTHERS.first_minute, t_first_min)> 0 THEN LEAST(OTHERS.last_minute, t_last_min) - GREATEST(OTHERS.first_minute, t_first_min) ELSE 0 END) as shared_minutes,
    COUNT(*) AS shared_matches,
    competitions.name as comp_name,
    competitions.season as comp_season,
    TEXT(CASE WHEN GAMES_IN_SCOPE.target_at_home THEN HT.name ELSE AT.name END) AS team_name
FROM (
SELECT
    match_id,
    at_home AS target_at_home,
    first_minute AS t_first_min,
    last_minute AS t_last_min
FROM player_matches
INNER JOIN players ON players.id = player_matches.player_id
WHERE players.id = $1
) AS GAMES_IN_SCOPE
INNER JOIN player_matches AS OTHERS ON (OTHERS.match_id = GAMES_IN_SCOPE.match_id AND OTHERS.at_home = GAMES_IN_SCOPE.target_at_home)
RIGHT JOIN players AS P ON P.id = OTHERS.player_id 
INNER JOIN league_matches ON GAMES_IN_SCOPE.match_id = league_matches.id
INNER JOIN competitions ON league_matches.competition_id = competitions.id
INNER JOIN teams AS HT ON HT.id = league_matches.home_team_id
INNER JOIN teams AS AT ON AT.id = league_matches.away_team_id
WHERE P.id = $2
GROUP BY P.name, P.url, P.id, competitions.name, competitions.season, team_name
HAVING SUM(LEAST(OTHERS.last_minute, t_last_min) - GREATEST(OTHERS.first_minute, t_first_min))> 0
ORDER BY shared_minutes DESC;


-- name: GetCompTableFromPM :many
SELECT
    team_name,
    SUM(matches_played) as matches_played,
    SUM(goals_scored) as goals_scored,
    SUM(goals_conceded) as goals_conceded,
    SUM(goal_difference) as goal_difference,
    SUM(wins) as wins,
    SUM(draws) as draws,
    SUM(losses) as losses,
    SUM(points) as points
FROM (
SELECT
    home_team_name as team_name,
    COUNT(*) as matches_played,
    SUM(home_goals_scored) + SUM(away_ogs) as goals_scored,
    SUM(away_goals_scored) + SUM(home_ogs) as goals_conceded,
    SUM(home_goals_scored) + SUM(away_ogs) - SUM(away_goals_scored) - SUM(home_ogs) as goal_difference,
    SUM(CASE WHEN home_goals_scored + away_ogs > away_goals_scored + home_ogs THEN 1 ELSE 0 END) as wins,
    SUM(CASE WHEN home_goals_scored + away_ogs = away_goals_scored + home_ogs THEN 1 ELSE 0 END) as draws,
    SUM(CASE WHEN home_goals_scored + away_ogs < away_goals_scored + home_ogs THEN 1 ELSE 0 END) as losses,
    SUM((CASE WHEN home_goals_scored + away_ogs > away_goals_scored + home_ogs THEN 1 ELSE 0 END)*3 + (CASE WHEN home_goals_scored + away_ogs = away_goals_scored + home_ogs THEN 1 ELSE 0 END)) as points
FROM (SELECT
    PM.match_id as match_id,
    HT.name as home_team_name,
    AT.name as away_team_name,
    SUM(CASE WHEN PM.at_home THEN PM.goals ELSE 0 END) as home_goals_scored,
    SUM(CASE WHEN NOT PM.at_home THEN PM.goals ELSE 0 END) as away_goals_scored,
    SUM(CASE WHEN PM.at_home THEN PM.red_card ELSE 0 END) as home_red_cards,
    SUM(CASE WHEN NOT PM.at_home THEN PM.red_card ELSE 0 END) as away_red_cards,
    SUM(CASE WHEN PM.at_home THEN PM.yellow_card ELSE 0 END) as home_yellow_cards,
    SUM(CASE WHEN NOT PM.at_home THEN PM.yellow_card ELSE 0 END) as away_yellow_cards,
    SUM(CASE WHEN PM.at_home THEN PM.penalties ELSE 0 END) as home_penalties,
    SUM(CASE WHEN NOT PM.at_home THEN PM.penalties ELSE 0 END) as away_penalties,
    SUM(CASE WHEN PM.at_home THEN PM.own_goals ELSE 0 END) as home_ogs,
    SUM(CASE WHEN NOT PM.at_home THEN PM.own_goals ELSE 0 END) as away_ogs
FROM player_matches AS PM
INNER JOIN league_matches AS LM ON PM.match_id = LM.id
INNER JOIN competitions AS C ON C.id = LM.competition_id
INNER JOIN players ON PM.player_id = players.id
INNER JOIN teams AS HT ON HT.id = LM.home_team_id
INNER JOIN teams AS AT ON AT.id = LM.away_team_id
WHERE C.name = $1 AND C.season = $2 AND NOT players.name = 'fakeRedCard'
GROUP BY HT.name, AT.name, PM.match_id
) s
GROUP BY home_team_name
UNION ALL
SELECT
    away_team_name as team_name,
    COUNT(*) as matches_played,
    SUM(away_goals_scored) + SUM(home_ogs) as goals_scored,
    SUM(home_goals_scored) + SUM(away_ogs) as goals_conceded,
    SUM(away_goals_scored) + SUM(home_ogs) - SUM(home_goals_scored) - SUM(away_ogs) as goal_difference,
    SUM(CASE WHEN home_goals_scored + away_ogs < away_goals_scored + home_ogs THEN 1 ELSE 0 END) as wins,
    SUM(CASE WHEN home_goals_scored + away_ogs = away_goals_scored + home_ogs THEN 1 ELSE 0 END) as draws,
    SUM(CASE WHEN home_goals_scored + away_ogs > away_goals_scored + home_ogs THEN 1 ELSE 0 END) as losses,
    SUM((CASE WHEN home_goals_scored + away_ogs < away_goals_scored + home_ogs THEN 1 ELSE 0 END)*3 + (CASE WHEN home_goals_scored + away_ogs = away_goals_scored + home_ogs THEN 1 ELSE 0 END)) as points
FROM (SELECT
    PM.match_id as match_id,
    HT.name as home_team_name,
    AT.name as away_team_name,
    SUM(CASE WHEN PM.at_home THEN PM.goals ELSE 0 END) as home_goals_scored,
    SUM(CASE WHEN NOT PM.at_home THEN PM.goals ELSE 0 END) as away_goals_scored,
    SUM(CASE WHEN PM.at_home THEN PM.red_card ELSE 0 END) as home_red_cards,
    SUM(CASE WHEN NOT PM.at_home THEN PM.red_card ELSE 0 END) as away_red_cards,
    SUM(CASE WHEN PM.at_home THEN PM.yellow_card ELSE 0 END) as home_yellow_cards,
    SUM(CASE WHEN NOT PM.at_home THEN PM.yellow_card ELSE 0 END) as away_yellow_cards,
    SUM(CASE WHEN PM.at_home THEN PM.penalties ELSE 0 END) as home_penalties,
    SUM(CASE WHEN NOT PM.at_home THEN PM.penalties ELSE 0 END) as away_penalties,
    SUM(CASE WHEN PM.at_home THEN PM.own_goals ELSE 0 END) as home_ogs,
    SUM(CASE WHEN NOT PM.at_home THEN PM.own_goals ELSE 0 END) as away_ogs
FROM player_matches AS PM
INNER JOIN league_matches AS LM ON PM.match_id = LM.id
INNER JOIN competitions AS C ON C.id = LM.competition_id
INNER JOIN players ON PM.player_id = players.id
INNER JOIN teams AS HT ON HT.id = LM.home_team_id
INNER JOIN teams AS AT ON AT.id = LM.away_team_id
WHERE C.name = $1 AND C.season = $2 AND NOT players.name = 'fakeRedCard'
GROUP BY HT.name, AT.name, PM.match_id
) s
GROUP BY away_team_name
) s
GROUP BY team_name
ORDER BY points DESC, goal_difference DESC, goals_scored DESC, goals_conceded DESC, wins DESC;

-- name: GetMatchRecordsFromPM :many
SELECT
    PM.match_id as match_id,
    HT.name as home_team_name,
    AT.name as away_team_name,
    SUM(CASE WHEN PM.at_home THEN PM.goals ELSE 0 END) as home_goals_scored,
    SUM(CASE WHEN NOT PM.at_home THEN PM.goals ELSE 0 END) as away_goals_scored,
    SUM(CASE WHEN PM.at_home THEN PM.red_card ELSE 0 END) as home_red_cards,
    SUM(CASE WHEN NOT PM.at_home THEN PM.red_card ELSE 0 END) as away_red_cards,
    SUM(CASE WHEN PM.at_home THEN PM.yellow_card ELSE 0 END) as home_yellow_cards,
    SUM(CASE WHEN NOT PM.at_home THEN PM.yellow_card ELSE 0 END) as away_yellow_cards,
    SUM(CASE WHEN PM.at_home THEN PM.penalties ELSE 0 END) as home_penalties,
    SUM(CASE WHEN NOT PM.at_home THEN PM.penalties ELSE 0 END) as away_penalties,
    SUM(CASE WHEN PM.at_home THEN PM.own_goals ELSE 0 END) as home_ogs,
    SUM(CASE WHEN NOT PM.at_home THEN PM.own_goals ELSE 0 END) as away_ogs
FROM player_matches AS PM
INNER JOIN league_matches AS LM ON PM.match_id = LM.id
INNER JOIN competitions AS C ON C.id = LM.competition_id
INNER JOIN players ON PM.player_id = players.id
INNER JOIN teams AS HT ON HT.id = LM.home_team_id
INNER JOIN teams AS AT ON AT.id = LM.away_team_id
WHERE C.name = $1 AND C.season = $2 AND NOT players.name = 'fakeRedCard'
GROUP BY HT.name, AT.name, PM.match_id;



-- SELECT
--     team_name,
--     SUM(home_goals_scored) + SUM(away_goals_scored) as total_goals_scored,
--     SUM(home_goals_conceded) + SUM(away_goals_conceded) as total_goals_conceded,
--     SUM(home_goals_scored) + SUM(away_goals_scored) - SUM(home_goals_conceded) - SUM(away_goals_conceded) as goal_difference,
--     -- SUM(wins) as total_wins,
--     -- SUM(draws) as total_draws,
--     -- SUM(losses) as total_losses,
--     -- SUM(points) as total_points,
--     SUM((CASE WHEN home_goals_scored>home_goals_conceded THEN 1 ELSE 0 END)+(CASE WHEN away_goals_scored>away_goals_conceded THEN 1 ELSE 0 END)) AS wins,
--     SUM((CASE WHEN home_goals_scored=home_goals_conceded THEN 1 ELSE 0 END)+(CASE WHEN away_goals_scored=away_goals_conceded THEN 1 ELSE 0 END)) AS draws,
--     SUM((CASE WHEN home_goals_scored<home_goals_conceded THEN 1 ELSE 0 END)+(CASE WHEN away_goals_scored<away_goals_conceded THEN 1 ELSE 0 END)) AS losses,
--     -- SUM((CASE WHEN goals_scored>goals_conceded THEN 1 ELSE 0 END)*3 + (CASE WHEN goals_scored=goals_conceded THEN 1 ELSE 0 END)) AS points,
--     COUNT(DISTINCT match_id) as matches_played,
--     -- COUNT(DISTINCT player_id) as players_fielded,
--     SUM(red_cards) as total_red_cards,
--     SUM(yellow_cards) as total_yellow_cards,
--     SUM(pens) as total_pens,
--     SUM(ogs) as total_ogs
-- FROM(
-- SELECT
--     PM.match_id as match_id,
--     SUM(CASE WHEN PM.at_home IS TRUE THEN PM.goals ELSE 0 END) as home_goals_scored,
--     SUM(CASE WHEN PM.at_home IS FALSE THEN PM.goals ELSE 0 END) as home_goals_conceded,
--     0 as away_goals_scored,
--     0 as away_goals_conceded,
--     SUM(CASE WHEN PM.at_home IS TRUE THEN PM.red_card ELSE 0 END) as red_cards,
--     SUM(CASE WHEN PM.at_home IS TRUE THEN PM.yellow_card ELSE 0 END) as yellow_cards,
--     SUM(CASE WHEN PM.at_home IS TRUE THEN PM.penalties ELSE 0 END) as pens,
--     SUM(CASE WHEN PM.at_home IS TRUE THEN PM.own_goals ELSE 0 END) as ogs,
--     HT.name as team_name
-- FROM player_matches AS PM
-- INNER JOIN league_matches AS LM ON PM.match_id = LM.id
-- INNER JOIN teams AS HT ON HT.id = LM.home_team_id
-- INNER JOIN competitions AS C ON C.id = LM.competition_id
-- INNER JOIN players ON PM.player_id = players.id 
-- WHERE C.name = $1 AND C.season = $2 AND NOT players.name = 'fakeRedCard'
-- GROUP BY team_name, match_id
-- UNION ALL
-- SELECT
--     PM.match_id as match_id,
--     0 as home_goals_scored,
--     0 as home_goals_conceded,
--     SUM(CASE WHEN PM.at_home IS FALSE THEN PM.goals ELSE 0 END) as away_goals_scored,
--     SUM(CASE WHEN PM.at_home IS TRUE THEN PM.goals ELSE 0 END) as away_goals_conceded,
--     SUM(CASE WHEN PM.at_home IS FALSE THEN PM.red_card ELSE 0 END) as red_cards,
--     SUM(CASE WHEN PM.at_home IS FALSE THEN PM.yellow_card ELSE 0 END) as yellow_cards,
--     SUM(CASE WHEN PM.at_home IS FALSE THEN PM.penalties ELSE 0 END) as pens,
--     SUM(CASE WHEN PM.at_home IS FALSE THEN PM.own_goals ELSE 0 END) as ogs,
--     AT.name as team_name
-- FROM player_matches AS PM
-- INNER JOIN league_matches AS LM ON PM.match_id = LM.id
-- INNER JOIN teams AS AT ON AT.id = LM.away_team_id
-- INNER JOIN competitions AS C ON C.id = LM.competition_id
-- INNER JOIN players ON PM.player_id = players.id 
-- WHERE C.name = $1 AND C.season = $2 AND NOT players.name = 'fakeRedCard'
-- GROUP BY team_name, match_id
-- ) s
-- GROUP BY team_name
-- ORDER BY total_goals_scored DESC;
