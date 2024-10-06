-- name: GetChannelById :one
SELECT * FROM channels
WHERE id = $1 LIMIT 1;

-- name: GetChannelByExternalId :one
SELECT * FROM channels
WHERE ext_id = $1 LIMIT 1;

-- name: GetChannelByName :one
SELECT * FROM channels
WHERE name = $1 LIMIT 1;

-- name: InsertChannel :one
INSERT INTO channels (id, ext_id, name, description, image_path, generate_thumbnails) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetChannels :many
SELECT * FROM channels
ORDER BY name ASC;

-- name: GetChannelNames :many
SELECT name FROM channels
ORDER BY name ASC;