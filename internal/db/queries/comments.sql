-- name: GetCommentById :one
SELECT * FROM comments
WHERE id = $1 LIMIT 1;

-- name: GetCommentsFilter :many
SELECT * FROM comments
WHERE video_id = $1
ORDER BY like_count DESC
LIMIT $2 OFFSET $3;

-- name: GetAllCommentsByVideoId :many
SELECT * FROM comments
WHERE video_id = $1
ORDER BY like_count DESC;

-- name: GetCommentCountByVideoId :one
SELECT count(id) FROM comments
WHERE video_id = $1;

-- name: GetRootCommentsForVideoId :many
WITH top_comments AS (
    SELECT *
    FROM comments
    WHERE comments.video_id = $1
    AND comments.parent IS NULL
    ORDER BY comments.like_count DESC NULLS LAST
    LIMIT $2 OFFSET $3
)
SELECT
    tc.*,
    COALESCE(r.reply_count, 0) AS reply_count
FROM top_comments tc
LEFT JOIN (
    SELECT parent, COUNT(*) AS reply_count
    FROM comments
    WHERE comments.parent IN (SELECT id FROM top_comments)
    GROUP BY comments.parent
) r ON tc.id = r.parent
ORDER BY tc.like_count DESC NULLS LAST;

-- name: GetCommentRepliesByCommentId :many
SELECT * FROM comments
WHERE parent = $1
ORDER BY like_count DESC;

-- -- name: InsertComment :one
-- INSERT INTO comments (
--   id, text, timestamp, like_count, is_favorited, author, author_id, author_thumbnail, author_is_uploader, parent, video_id
-- ) VALUES (
--   $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
-- ) RETURNING *;

-- https://docs.sqlc.dev/en/stable/howto/insert.html#using-copyfrom
-- name: BatchInsertComments :copyfrom
INSERT INTO comments (
  id, text, timestamp, like_count, is_favorited, author, author_id, author_thumbnail, author_is_uploader, parent, video_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
);
