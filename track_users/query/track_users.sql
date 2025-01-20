-- name: CreateTrackUser :exec
INSERT INTO track_user (name, counter)
VALUES ($1, $2);


-- name: UpdateUserCounter :exec
UPDATE track_user
SET counter = counter::BIGINT + 1::BIGINT
WHERE name = $1;

-- name: GetCounterUser :one
SELECT id
FROM track_user
WHERE name = $1;