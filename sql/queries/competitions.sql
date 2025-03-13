-- name: CreateCompetition :one
INSERT INTO competitions (id, name)
VALUES (
    $1,
    $2
)
RETURNING *;

-- name: CheckIfCompetitionExistsById :one
SELECT EXISTS (SELECT 1
FROM competitions
WHERE id = $1
LIMIT 1);

-- name: CheckIfCompetitionExistsByName :one
SELECT EXISTS (SELECT 1
FROM competitions
WHERE name = $1
LIMIT 1);

-- name: GetCompetitionIdFromName :one
SELECT id
FROM competitions
WHERE name = $1;

-- name: GetCompetitionNameFromId :one
SELECT name
FROM competitions
WHERE id = $1;

-- name: DeleteCompetitions :exec
DELETE FROM competitions;