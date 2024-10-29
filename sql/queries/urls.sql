-- name: CreateURL :one
INSERT INTO urls (alias, url)
VALUES (?, ?)
RETURNING *;

-- name: ListURLs :many
SELECT * FROM urls
ORDER BY id;
