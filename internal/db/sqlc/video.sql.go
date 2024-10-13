// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: video.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const ftsVideosFilter = `-- name: FtsVideosFilter :many
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
  sr.id, sr.ext_id, sr.title, sr.description, sr.upload_date, sr.uploader, sr.duration, sr.view_count, sr.like_count, sr.dislike_count, sr.format, sr.height, sr.width, sr.resolution, sr.fps, sr.video_codec, sr.vbr, sr.audio_codec, sr.abr, sr.comment_count, sr.video_path, sr.thumbnail_path, sr.info_path, sr.subtitle_path, sr.path, sr.storyboard_path, sr.channel_id, sr.created_at, sr.updated_at, sr.channel_name, sr.channel_image_path, sr.rank,
  COUNT(*) OVER () AS total_count
FROM search_results sr
ORDER BY sr.rank DESC
LIMIT $2 OFFSET $3
`

type FtsVideosFilterParams struct {
	WebsearchToTsquery string
	Limit              int32
	Offset             int32
}

type FtsVideosFilterRow struct {
	ID               pgtype.UUID
	ExtID            *string
	Title            string
	Description      *string
	UploadDate       pgtype.Timestamptz
	Uploader         *string
	Duration         int32
	ViewCount        int64
	LikeCount        *int64
	DislikeCount     *int64
	Format           *string
	Height           *int32
	Width            *int32
	Resolution       *string
	Fps              *float32
	VideoCodec       *string
	Vbr              *float32
	AudioCodec       *string
	Abr              *float32
	CommentCount     *int32
	VideoPath        string
	ThumbnailPath    string
	InfoPath         string
	SubtitlePath     []string
	Path             string
	StoryboardPath   *string
	ChannelID        pgtype.UUID
	CreatedAt        pgtype.Timestamptz
	UpdatedAt        pgtype.Timestamptz
	ChannelName      string
	ChannelImagePath *string
	Rank             float32
	TotalCount       int64
}

func (q *Queries) FtsVideosFilter(ctx context.Context, arg FtsVideosFilterParams) ([]FtsVideosFilterRow, error) {
	rows, err := q.db.Query(ctx, ftsVideosFilter, arg.WebsearchToTsquery, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FtsVideosFilterRow
	for rows.Next() {
		var i FtsVideosFilterRow
		if err := rows.Scan(
			&i.ID,
			&i.ExtID,
			&i.Title,
			&i.Description,
			&i.UploadDate,
			&i.Uploader,
			&i.Duration,
			&i.ViewCount,
			&i.LikeCount,
			&i.DislikeCount,
			&i.Format,
			&i.Height,
			&i.Width,
			&i.Resolution,
			&i.Fps,
			&i.VideoCodec,
			&i.Vbr,
			&i.AudioCodec,
			&i.Abr,
			&i.CommentCount,
			&i.VideoPath,
			&i.ThumbnailPath,
			&i.InfoPath,
			&i.SubtitlePath,
			&i.Path,
			&i.StoryboardPath,
			&i.ChannelID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ChannelName,
			&i.ChannelImagePath,
			&i.Rank,
			&i.TotalCount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getVideoByExternalID = `-- name: GetVideoByExternalID :one
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
WHERE ext_id = $1 LIMIT 1
`

type GetVideoByExternalIDRow struct {
	ID             pgtype.UUID
	ExtID          *string
	Title          string
	Description    *string
	UploadDate     pgtype.Timestamptz
	Uploader       *string
	Duration       int32
	ViewCount      int64
	LikeCount      *int64
	DislikeCount   *int64
	Format         *string
	Height         *int32
	Width          *int32
	Resolution     *string
	Fps            *float32
	VideoCodec     *string
	Vbr            *float32
	AudioCodec     *string
	Abr            *float32
	CommentCount   *int32
	VideoPath      string
	ThumbnailPath  string
	InfoPath       string
	SubtitlePath   []string
	Path           string
	StoryboardPath *string
	ChannelID      pgtype.UUID
	CreatedAt      pgtype.Timestamptz
	UpdatedAt      pgtype.Timestamptz
}

func (q *Queries) GetVideoByExternalID(ctx context.Context, extID *string) (GetVideoByExternalIDRow, error) {
	row := q.db.QueryRow(ctx, getVideoByExternalID, extID)
	var i GetVideoByExternalIDRow
	err := row.Scan(
		&i.ID,
		&i.ExtID,
		&i.Title,
		&i.Description,
		&i.UploadDate,
		&i.Uploader,
		&i.Duration,
		&i.ViewCount,
		&i.LikeCount,
		&i.DislikeCount,
		&i.Format,
		&i.Height,
		&i.Width,
		&i.Resolution,
		&i.Fps,
		&i.VideoCodec,
		&i.Vbr,
		&i.AudioCodec,
		&i.Abr,
		&i.CommentCount,
		&i.VideoPath,
		&i.ThumbnailPath,
		&i.InfoPath,
		&i.SubtitlePath,
		&i.Path,
		&i.StoryboardPath,
		&i.ChannelID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getVideoById = `-- name: GetVideoById :one
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
WHERE id = $1 LIMIT 1
`

type GetVideoByIdRow struct {
	ID             pgtype.UUID
	ExtID          *string
	Title          string
	Description    *string
	UploadDate     pgtype.Timestamptz
	Uploader       *string
	Duration       int32
	ViewCount      int64
	LikeCount      *int64
	DislikeCount   *int64
	Format         *string
	Height         *int32
	Width          *int32
	Resolution     *string
	Fps            *float32
	VideoCodec     *string
	Vbr            *float32
	AudioCodec     *string
	Abr            *float32
	CommentCount   *int32
	VideoPath      string
	ThumbnailPath  string
	InfoPath       string
	SubtitlePath   []string
	Path           string
	StoryboardPath *string
	ChannelID      pgtype.UUID
	CreatedAt      pgtype.Timestamptz
	UpdatedAt      pgtype.Timestamptz
}

func (q *Queries) GetVideoById(ctx context.Context, id pgtype.UUID) (GetVideoByIdRow, error) {
	row := q.db.QueryRow(ctx, getVideoById, id)
	var i GetVideoByIdRow
	err := row.Scan(
		&i.ID,
		&i.ExtID,
		&i.Title,
		&i.Description,
		&i.UploadDate,
		&i.Uploader,
		&i.Duration,
		&i.ViewCount,
		&i.LikeCount,
		&i.DislikeCount,
		&i.Format,
		&i.Height,
		&i.Width,
		&i.Resolution,
		&i.Fps,
		&i.VideoCodec,
		&i.Vbr,
		&i.AudioCodec,
		&i.Abr,
		&i.CommentCount,
		&i.VideoPath,
		&i.ThumbnailPath,
		&i.InfoPath,
		&i.SubtitlePath,
		&i.Path,
		&i.StoryboardPath,
		&i.ChannelID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getVideoInfoPaths = `-- name: GetVideoInfoPaths :many
SELECT info_path FROM videos
`

func (q *Queries) GetVideoInfoPaths(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, getVideoInfoPaths)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var info_path string
		if err := rows.Scan(&info_path); err != nil {
			return nil, err
		}
		items = append(items, info_path)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getVideosCount = `-- name: GetVideosCount :one
SELECT COUNT(id) FROM videos
WHERE channel_id = $1
`

func (q *Queries) GetVideosCount(ctx context.Context, channelID pgtype.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, getVideosCount, channelID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getVideosFilter = `-- name: GetVideosFilter :many
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
LIMIT $2 OFFSET $3
`

type GetVideosFilterParams struct {
	ChannelID pgtype.UUID
	Limit     int32
	Offset    int32
}

type GetVideosFilterRow struct {
	ID             pgtype.UUID
	ExtID          *string
	Title          string
	Description    *string
	UploadDate     pgtype.Timestamptz
	Uploader       *string
	Duration       int32
	ViewCount      int64
	LikeCount      *int64
	DislikeCount   *int64
	Format         *string
	Height         *int32
	Width          *int32
	Resolution     *string
	Fps            *float32
	VideoCodec     *string
	Vbr            *float32
	AudioCodec     *string
	Abr            *float32
	CommentCount   *int32
	VideoPath      string
	ThumbnailPath  string
	InfoPath       string
	SubtitlePath   []string
	Path           string
	StoryboardPath *string
	ChannelID      pgtype.UUID
	CreatedAt      pgtype.Timestamptz
	UpdatedAt      pgtype.Timestamptz
}

func (q *Queries) GetVideosFilter(ctx context.Context, arg GetVideosFilterParams) ([]GetVideosFilterRow, error) {
	rows, err := q.db.Query(ctx, getVideosFilter, arg.ChannelID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetVideosFilterRow
	for rows.Next() {
		var i GetVideosFilterRow
		if err := rows.Scan(
			&i.ID,
			&i.ExtID,
			&i.Title,
			&i.Description,
			&i.UploadDate,
			&i.Uploader,
			&i.Duration,
			&i.ViewCount,
			&i.LikeCount,
			&i.DislikeCount,
			&i.Format,
			&i.Height,
			&i.Width,
			&i.Resolution,
			&i.Fps,
			&i.VideoCodec,
			&i.Vbr,
			&i.AudioCodec,
			&i.Abr,
			&i.CommentCount,
			&i.VideoPath,
			&i.ThumbnailPath,
			&i.InfoPath,
			&i.SubtitlePath,
			&i.Path,
			&i.StoryboardPath,
			&i.ChannelID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertVideo = `-- name: InsertVideo :one
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
) RETURNING id, ext_id, title, description, upload_date, uploader, duration, view_count, like_count, dislike_count, format, height, width, resolution, fps, video_codec, vbr, audio_codec, abr, comment_count, video_path, thumbnail_path, info_path, subtitle_path, path, storyboard_path, created_at, updated_at, channel_id, title_fts_en
`

type InsertVideoParams struct {
	ID             pgtype.UUID
	ExtID          *string
	Title          string
	Description    *string
	UploadDate     pgtype.Timestamptz
	Uploader       *string
	Duration       int32
	ViewCount      int64
	LikeCount      *int64
	DislikeCount   *int64
	Format         *string
	Height         *int32
	Width          *int32
	Resolution     *string
	Fps            *float32
	VideoCodec     *string
	Vbr            *float32
	AudioCodec     *string
	Abr            *float32
	CommentCount   *int32
	VideoPath      string
	ThumbnailPath  string
	InfoPath       string
	SubtitlePath   []string
	Path           string
	StoryboardPath *string
	ChannelID      pgtype.UUID
}

func (q *Queries) InsertVideo(ctx context.Context, arg InsertVideoParams) (Video, error) {
	row := q.db.QueryRow(ctx, insertVideo,
		arg.ID,
		arg.ExtID,
		arg.Title,
		arg.Description,
		arg.UploadDate,
		arg.Uploader,
		arg.Duration,
		arg.ViewCount,
		arg.LikeCount,
		arg.DislikeCount,
		arg.Format,
		arg.Height,
		arg.Width,
		arg.Resolution,
		arg.Fps,
		arg.VideoCodec,
		arg.Vbr,
		arg.AudioCodec,
		arg.Abr,
		arg.CommentCount,
		arg.VideoPath,
		arg.ThumbnailPath,
		arg.InfoPath,
		arg.SubtitlePath,
		arg.Path,
		arg.StoryboardPath,
		arg.ChannelID,
	)
	var i Video
	err := row.Scan(
		&i.ID,
		&i.ExtID,
		&i.Title,
		&i.Description,
		&i.UploadDate,
		&i.Uploader,
		&i.Duration,
		&i.ViewCount,
		&i.LikeCount,
		&i.DislikeCount,
		&i.Format,
		&i.Height,
		&i.Width,
		&i.Resolution,
		&i.Fps,
		&i.VideoCodec,
		&i.Vbr,
		&i.AudioCodec,
		&i.Abr,
		&i.CommentCount,
		&i.VideoPath,
		&i.ThumbnailPath,
		&i.InfoPath,
		&i.SubtitlePath,
		&i.Path,
		&i.StoryboardPath,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ChannelID,
		&i.TitleFtsEn,
	)
	return i, err
}
