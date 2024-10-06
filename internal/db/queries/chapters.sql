-- https://docs.sqlc.dev/en/stable/howto/insert.html#using-copyfrom
-- name: BatchInsertChapters :copyfrom
INSERT INTO chapters (
  id, title, start_time, end_time, video_id
) VALUES (
  $1, $2, $3, $4, $5
);


-- name: GetChaptersByVideoId :many
SELECT * FROM chapters
WHERE video_id = $1
ORDER BY start_time ASC;
