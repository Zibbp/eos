-- name: GetVideoById :one
SELECT * FROM videos
WHERE id = $1 LIMIT 1;

-- name: GetVideoByExternalID :one
SELECT * FROM videos
WHERE ext_id = $1 LIMIT 1;

-- name: GetVideoInfoPaths :many
SELECT info_path FROM videos;

-- name: InsertVideo :one
INSERT INTO videos (
  id, ext_id, title, description, upload_date, uploader, duration, view_count,
  like_count, dislike_count, format, height, width, resolution, fps, video_codec,
  vbr, audio_codec, abr, comment_count, video_path, thumbnail_path, info_path,
  subtitle_path, path, storyboard_path, channel_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8,
  $9, $10, $11, $12, $13, $14, $15, $16,
  $17, $18, $19, $20, $21, $22, $23,
  $24, $25, $26, $27
) RETURNING *;

-- name: GetVideosFilter :many
SELECT * FROM videos
WHERE channel_id = $1
ORDER BY upload_date DESC
LIMIT $2 OFFSET $3;

-- name: GetVideosCount :one
SELECT COUNT(id) FROM videos
WHERE channel_id = $1;