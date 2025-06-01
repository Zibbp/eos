-- name: InsertBlockedPath :one
INSERT INTO blocked_paths (id, path, error_text) VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteBlockedPathByPath :exec
DELETE FROM blocked_paths WHERE path = $1;

-- name: DeleteBlockedPathById :exec
DELETE FROM blocked_paths WHERE id = $1;

-- name: GetBlockedPath :one
SELECT * FROM blocked_paths
WHERE path = $1;

-- name: GetBlockedPaths :many
SELECT * FROM blocked_paths;

-- name: IncrementBlockedPathErrorCount :exec
UPDATE blocked_paths
SET
  error_count = error_count + 1,
  error_text = $2
WHERE path = $1;

-- name: SetBlockedPathAsBlocked :exec
UPDATE blocked_paths
SET is_blocked = true
WHERE path = $1;


-- name: GetTotalBlockedPaths :one
SELECT COUNT(*) AS total FROM blocked_paths;
