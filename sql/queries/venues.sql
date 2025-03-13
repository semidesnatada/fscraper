-- name: CreateVenue :one
INSERT INTO venues (id, name)
VALUES (
    $1,
    $2
)
RETURNING *;

-- name: CheckIfVenueExistsById :one
SELECT EXISTS (SELECT 1
FROM venues
WHERE id = $1
LIMIT 1);

-- name: CheckIfVenueExistsByName :one
SELECT EXISTS (SELECT 1
FROM venues
WHERE name = $1
LIMIT 1);

-- name: GetVenueIdFromName :one
SELECT id
FROM venues
WHERE name = $1;

-- name: GetVenueNameFromId :one
SELECT name
FROM venues
WHERE id = $1;

-- name: DeleteVenues :exec
DELETE FROM venues;