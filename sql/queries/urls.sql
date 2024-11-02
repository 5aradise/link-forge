-- name: CreateURL :one
INSERT INTO urls (alias, url)
VALUES (?, ?)
RETURNING *;

-- name: ListURLs :many
SELECT * FROM urls
ORDER BY id;

-- name: GetURLByAlias :one
SELECT * FROM urls
WHERE alias = ?;

-- name: DeleteURLByAlias :one
DELETE FROM urls
WHERE alias = ?
RETURNING *;
