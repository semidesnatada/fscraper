-- name: CreatePlayerMatch :one
INSERT INTO player_matches (
    match_id,
    player_id,
    match_url,
    first_minute,
    last_minute,
    goals,
    penalties,
    yellow_card,
    red_card,
    own_goals,
    is_knockout,
    at_home
)
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
    $12
)
RETURNING *;

-- name: CheckIfPlayerMatchExistsByIds :one
SELECT EXISTS (SELECT 1
FROM player_matches
WHERE match_id = $1 and player_id = $2
LIMIT 1);

-- name: CheckIfMatchIDExistsInPlayerMatches :one
SELECT EXISTS (SELECT 1
FROM player_matches
WHERE match_id = $1
LIMIT 1);

-- name: CheckIfMatchUrlExistsInPlayerMatches :one
SELECT EXISTS (SELECT 1
FROM player_matches
WHERE match_url = $1
LIMIT 1);

-- name: GetPlayerMatchFromIds :one
SELECT *
FROM player_matches
WHERE match_id = $1 and player_id = $2;

-- name: DeletePlayerMatches :exec
DELETE FROM player_matches;