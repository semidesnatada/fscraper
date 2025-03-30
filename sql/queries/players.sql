-- name: CreatePlayer :one
INSERT INTO players (id, name, nationality, url)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: CheckIfPlayerExistsById :one
SELECT EXISTS (SELECT 1
FROM players
WHERE id = $1
LIMIT 1);

-- name: CheckIfPlayerExistsByUrl :one
SELECT EXISTS (SELECT 1
FROM players
WHERE url = $1
LIMIT 1);

-- name: GetPlayerNameFromId :one
SELECT name
FROM players
WHERE id = $1;

-- name: GetPlayerFromId :one
SELECT *
FROM players
WHERE id = $1;

-- name: GetPlayersByName :many
SELECT * FROM players
WHERE name = $1;

-- name: GetPlayerIdsFromName :many
SELECT id
FROM players
WHERE name = $1;

-- name: GetPlayerIdFromUrl :one
SELECT id
FROM players
WHERE url = $1;

-- name: GetPlayerUrlFromId :one
SELECT url
FROM players
WHERE id = $1;

-- name: GetPlayerUUIDsOrderedByUrl :many
SELECT id
FROM players
ORDER BY url
LIMIT 5000;

-- name: DeletePlayers :exec
DELETE FROM players;