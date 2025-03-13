-- name: CreateTeam :one
INSERT INTO teams (id, name)
VALUES (
    $1,
    $2
)
RETURNING *;

-- name: CheckIfTeamExistsById :one
SELECT EXISTS (SELECT 1
FROM teams
WHERE id = $1
LIMIT 1);

-- name: CheckIfTeamExistsByName :one
SELECT EXISTS (SELECT 1
FROM teams
WHERE name = $1
LIMIT 1);

-- name: GetTeamIdFromName :one
SELECT id
FROM teams
WHERE name = $1;

-- name: GetTeamNameFromId :one
SELECT name
FROM teams
WHERE id = $1;

-- name: DeleteTeams :exec
DELETE FROM teams;