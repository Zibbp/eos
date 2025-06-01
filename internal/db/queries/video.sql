-- name: GetVideoById :one
SELECT
    id,
    ext_id,
    title,
    description,
    upload_date,
    uploader,
    duration,
    view_count,
    like_count,
    dislike_count,
    format,
    height,
    width,
    resolution,
    fps,
    video_codec,
    vbr,
    audio_codec,
    abr,
    comment_count,
    video_path,
    thumbnail_path,
    info_path,
    subtitle_path,
    path,
    storyboard_path,
    channel_id,
    created_at,
    updated_at
FROM videos
WHERE id = $1 LIMIT 1;

-- name: GetVideoByExternalID :one
SELECT
    id,
    ext_id,
    title,
    description,
    upload_date,
    uploader,
    duration,
    view_count,
    like_count,
    dislike_count,
    format,
    height,
    width,
    resolution,
    fps,
    video_codec,
    vbr,
    audio_codec,
    abr,
    comment_count,
    video_path,
    thumbnail_path,
    info_path,
    subtitle_path,
    path,
    storyboard_path,
    channel_id,
    created_at,
    updated_at
FROM videos
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
SELECT
    id,
    ext_id,
    title,
    description,
    upload_date,
    uploader,
    duration,
    view_count,
    like_count,
    dislike_count,
    format,
    height,
    width,
    resolution,
    fps,
    video_codec,
    vbr,
    audio_codec,
    abr,
    comment_count,
    video_path,
    thumbnail_path,
    info_path,
    subtitle_path,
    path,
    storyboard_path,
    channel_id,
    created_at,
    updated_at
FROM videos
WHERE channel_id = $1
ORDER BY upload_date DESC
LIMIT $2 OFFSET $3;

-- name: GetVideosCount :one
SELECT COUNT(id) FROM videos
WHERE channel_id = $1;

-- name: FtsVideosFilter :many
WITH search_results AS (
  SELECT
    v.id,
    v.ext_id,
    v.title,
    v.description,
    v.upload_date,
    v.uploader,
    v.duration,
    v.view_count,
    v.like_count,
    v.dislike_count,
    v.format,
    v.height,
    v.width,
    v.resolution,
    v.fps,
    v.video_codec,
    v.vbr,
    v.audio_codec,
    v.abr,
    v.comment_count,
    v.video_path,
    v.thumbnail_path,
    v.info_path,
    v.subtitle_path,
    v.path,
    v.storyboard_path,
    v.channel_id,
    v.created_at,
    v.updated_at,
    c.name AS channel_name,
    c.image_path as channel_image_path,
    ts_rank(title_fts_en, websearch_to_tsquery($1)) AS rank
  FROM videos v
  JOIN channels c ON v.channel_id = c.id
  WHERE title_fts_en @@ websearch_to_tsquery($1)
)
SELECT 
  sr.*,
  COUNT(*) OVER () AS total_count
FROM search_results sr
ORDER BY sr.rank DESC
LIMIT $2 OFFSET $3;

-- name: GetTotalVideos :one
SELECT COUNT(*) AS total FROM videos;

-- name: GetTotalVideosByChannelId :one
SELECT COUNT(*) AS total FROM videos
WHERE channel_id = $1;
