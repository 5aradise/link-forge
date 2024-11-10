-- name: LoadState :one
SELECT alias_count FROM state
WHERE id = 1;

-- name: StoreState :exec
UPDATE state
SET alias_count = ?
WHERE id = 1;