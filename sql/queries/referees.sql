-- name: CreateReferee :one
INSERT INTO referees (id, name)
VALUES (
    $1,
    $2
)
RETURNING *;

-- name: CheckIfRefereeExistsById :one
SELECT EXISTS (SELECT 1
FROM referees
WHERE id = $1
LIMIT 1);

-- name: CheckIfRefereeExistsByName :one
SELECT EXISTS (SELECT 1
FROM referees
WHERE name = $1
LIMIT 1);

-- name: GetRefereeIdFromName :one
SELECT id
FROM referees
WHERE name = $1;

-- name: GetRefereeNameFromId :one
SELECT name
FROM referees
WHERE id = $1;

-- name: DeleteReferees :exec
DELETE FROM referees;