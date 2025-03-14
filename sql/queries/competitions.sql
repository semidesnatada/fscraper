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

-- name: DeleteCompetitions :exec
DELETE FROM competitions;