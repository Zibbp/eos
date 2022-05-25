package database

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
	"github.com/zibbp/eos/internal/video"
)

type VideoRow struct {
	ID            string          `json:"id"`
	Channel       string          `json:"channel"`
	ChannelID     string          `json:"channel_id"`
	Title         string          `json:"title"`
	Description   sql.NullString  `json:"description"`
	UploadDate    time.Time       `json:"upload_date"`
	Uploader      sql.NullString  `json:"uploader"`
	Duration      uint64          `json:"duration"`
	ViewCount     uint64          `json:"view_count"`
	LikeCount     sql.NullInt64   `json:"like_count"`
	DislikeCount  sql.NullInt64   `json:"dislike_count"`
	Format        sql.NullString  `json:"format"`
	Width         sql.NullInt64   `json:"width"`
	Height        sql.NullInt64   `json:"height"`
	Resolution    sql.NullString  `json:"resolution"`
	FPS           sql.NullFloat64 `json:"fps"`
	VideoCodec    sql.NullString  `json:"vcodec"`
	VBR           sql.NullFloat64 `json:"vbr"`
	AudioCodec    sql.NullString  `json:"acodec"`
	ABR           sql.NullFloat64 `json:"abr"`
	Epoch         sql.NullInt64   `json:"epoch"`
	CommentCount  uint64          `json:"comment_count"`
	Tags          sql.NullString  `json:"tags"`
	Categories    sql.NullString  `json:"categories"`
	VideoPath     string          `json:"video_path"`
	ThumbnailPath string          `json:"thumbnail_path"`
	JsonPath      sql.NullString  `json:"json_path"`
	SubtitlePath  sql.NullString  `json:"subtitle_path"`
	Chapters      []video.Chapter `json:"chapters"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type ChapterRow struct {
	ID        string  `json:"id"`
	StartTime float64 `json:"start_time"`
	Title     string  `json:"title"`
	EndTime   float64 `json:"end_time"`
	VideoID   string  `json:"video_id"`
}

func convertVideoRowToVideo(v VideoRow) video.Video {
	return video.Video{
		ID:            v.ID,
		Channel:       v.Channel,
		ChannelID:     v.ChannelID,
		Title:         v.Title,
		Description:   v.Description.String,
		UploadDate:    v.UploadDate,
		Uploader:      v.Uploader.String,
		Duration:      v.Duration,
		ViewCount:     v.ViewCount,
		LikeCount:     v.LikeCount.Int64,
		DislikeCount:  v.DislikeCount.Int64,
		Format:        v.Format.String,
		Width:         v.Width.Int64,
		Height:        v.Height.Int64,
		Resolution:    v.Resolution.String,
		FPS:           v.FPS.Float64,
		VideoCodec:    v.VideoCodec.String,
		VBR:           v.VBR.Float64,
		AudioCodec:    v.AudioCodec.String,
		ABR:           v.ABR.Float64,
		Epoch:         v.Epoch.Int64,
		CommentCount:  v.CommentCount,
		Tags:          v.Tags.String,
		Categories:    v.Categories.String,
		VideoPath:     v.VideoPath,
		ThumbnailPath: v.ThumbnailPath,
		JsonPath:      v.JsonPath.String,
		SubtitlePath:  v.SubtitlePath.String,
		Chapters:      v.Chapters,
		CreatedAt:     v.CreatedAt,
		UpdatedAt:     v.UpdatedAt,
	}
}

func (d *Database) GetVideo(c echo.Context, vidID string) (video.SingleVideoList, error) {
	var vid VideoRow
	var cha ChannelRow
	var chap video.Chapter
	// Video and Channel query
	row := d.Client.QueryRowContext(c.Request().Context(), "SELECT * FROM videos JOIN channels ON videos.channel_id = channels.id WHERE videos.id = $1", vidID)
	err := row.Scan(&vid.ID, &vid.Channel, &vid.ChannelID, &vid.Title, &vid.Description, &vid.UploadDate, &vid.Uploader, &vid.Duration, &vid.ViewCount, &vid.LikeCount, &vid.DislikeCount, &vid.Format, &vid.Width, &vid.Height, &vid.Resolution, &vid.FPS, &vid.VideoCodec, &vid.VBR, &vid.AudioCodec, &vid.ABR, &vid.Epoch, &vid.CommentCount, &vid.Tags, &vid.Categories, &vid.VideoPath, &vid.ThumbnailPath, &vid.JsonPath, &vid.SubtitlePath, &vid.CreatedAt, &vid.UpdatedAt, &cha.ID, &cha.Name, &cha.ChannelImagePath, &cha.CreatedAt, &cha.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return video.SingleVideoList{}, echo.ErrNotFound
		}
		log.Errorf("Failed to get video: %s", err)

		return video.SingleVideoList{}, fmt.Errorf("Failed to get video: %w", err)
	}
	// Chapter query
	rows, err := d.Client.QueryContext(c.Request().Context(), "SELECT * FROM chapters WHERE video_id = $1", vidID)
	if err != nil {
		log.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&chap.ID, &chap.StartTime, &chap.Title, &chap.EndTime, &chap.VideoID)
		if err != nil {
			log.Errorf("Failed to get video chapters: %s", err)
		}
		vid.Chapters = append(vid.Chapters, chap)
	}
	if err := rows.Err(); err != nil {
		log.Errorf("Failed to get video chapters: %s", err)
	}

	return video.SingleVideoList{
		Video:   convertVideoRowToVideo(vid),
		Channel: convertChannelRowToChannel(cha),
	}, nil
}

func (d *Database) CreateVideo(c echo.Context, vid video.Video) (video.Video, error) {
	vidRow := VideoRow{
		ID:            vid.ID,
		Channel:       vid.Channel,
		ChannelID:     vid.ChannelID,
		Title:         vid.Title,
		Description:   sql.NullString{String: vid.Description, Valid: true},
		UploadDate:    vid.UploadDate,
		Uploader:      sql.NullString{String: vid.Uploader, Valid: true},
		Duration:      vid.Duration,
		ViewCount:     vid.ViewCount,
		LikeCount:     sql.NullInt64{Int64: vid.LikeCount, Valid: true},
		DislikeCount:  sql.NullInt64{Int64: vid.DislikeCount, Valid: true},
		Format:        sql.NullString{String: vid.Format, Valid: true},
		Width:         sql.NullInt64{Int64: vid.Width, Valid: true},
		Height:        sql.NullInt64{Int64: vid.Height, Valid: true},
		Resolution:    sql.NullString{String: vid.Resolution, Valid: true},
		FPS:           sql.NullFloat64{Float64: vid.FPS, Valid: true},
		VideoCodec:    sql.NullString{String: vid.VideoCodec, Valid: true},
		VBR:           sql.NullFloat64{Float64: vid.VBR, Valid: true},
		AudioCodec:    sql.NullString{String: vid.AudioCodec, Valid: true},
		ABR:           sql.NullFloat64{Float64: vid.ABR, Valid: true},
		Epoch:         sql.NullInt64{Int64: vid.Epoch, Valid: true},
		CommentCount:  vid.CommentCount,
		Tags:          sql.NullString{String: vid.Tags, Valid: true},
		Categories:    sql.NullString{String: vid.Categories, Valid: true},
		VideoPath:     vid.VideoPath,
		ThumbnailPath: vid.ThumbnailPath,
		JsonPath:      sql.NullString{String: vid.JsonPath, Valid: true},
		SubtitlePath:  sql.NullString{String: vid.SubtitlePath, Valid: true},
		CreatedAt:     vid.CreatedAt,
		UpdatedAt:     vid.UpdatedAt,
	}
	// Set the default values
	vidRow.CreatedAt = time.Now()
	vidRow.UpdatedAt = time.Now()

	rows, err := d.Client.NamedQueryContext(c.Request().Context(), "INSERT INTO videos (id, channel, channel_id, title, description, upload_date, uploader, duration, view_count, like_count, dislike_count, format, width, height, resolution, fps, video_codec, vbr, audio_codec, abr, epoch, comment_count, tags, categories, video_path, thumbnail_path, json_path, subtitle_path, created_at, updated_at) VALUES (:id, :channel, :channelid, :title, :description, :uploaddate, :uploader, :duration, :viewcount, :likecount, :dislikecount, :format, :width, :height, :resolution, :fps, :videocodec, :vbr, :audiocodec, :abr, :epoch, :commentcount, :tags, :categories, :videopath, :thumbnailpath, :jsonpath, :subtitlepath, :createdat, :updatedat)", vidRow)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Errorf("Error creating video: %s: %s", err.Code, err.Message)
			switch err.Code {
			case "23505":
				return video.Video{}, echo.NewHTTPError(http.StatusConflict, "Video already exists")
			case "23503":
				return video.Video{}, echo.NewHTTPError(http.StatusNotFound, "Channel ID does not exist")
			default:
				log.Error(err)
				return video.Video{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
			}
		}
	}
	if err := rows.Close(); err != nil {
		return video.Video{}, fmt.Errorf("Failed to close rows: %w", err)
	}

	return vid, nil

}

func (d *Database) GetChannelVideos(c echo.Context, channelID string, vidList video.VideoList) (video.VideoList, error) {
	// Fetch channel first to check if it exists
	_, err := d.GetChannel(c, channelID)
	if err != nil {
		return video.VideoList{}, err
	}

	skip := (vidList.Page - 1) * vidList.Limit

	rows, err := d.Client.QueryContext(c.Request().Context(), "SELECT *, COUNT(*) OVER() AS TotalCount FROM videos WHERE channel_id = $1 ORDER BY upload_date DESC LIMIT $2 OFFSET $3", channelID, vidList.Limit, skip)
	if err != nil {
		log.Error(err)
		return video.VideoList{}, fmt.Errorf("Failed to get channel videos: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error(err)
		}
	}()
	var rowCount int
	var videos []video.Video
	for rows.Next() {
		var vid video.Video
		if err := rows.Scan(&vid.ID, &vid.Channel, &vid.ChannelID, &vid.Title, &vid.Description, &vid.UploadDate, &vid.Uploader, &vid.Duration, &vid.ViewCount, &vid.LikeCount, &vid.DislikeCount, &vid.Format, &vid.Width, &vid.Height, &vid.Resolution, &vid.FPS, &vid.VideoCodec, &vid.VBR, &vid.AudioCodec, &vid.ABR, &vid.Epoch, &vid.CommentCount, &vid.Tags, &vid.Categories, &vid.VideoPath, &vid.ThumbnailPath, &vid.JsonPath, &vid.SubtitlePath, &vid.CreatedAt, &vid.UpdatedAt, &rowCount); err != nil {
			log.Error(err)
			return video.VideoList{}, fmt.Errorf("Failed to scan video: %w", err)
		}

		videos = append(videos, vid)
	}

	// Set pagination info
	vidList.LastPage = int(math.Ceil(float64(rowCount) / float64(vidList.Limit)))
	vidList.TotalItems = rowCount

	// Set next page
	if vidList.Page < vidList.LastPage {
		vidList.NextPage = vidList.Page + 1
	} else {
		vidList.NextPage = vidList.LastPage
	}

	// Set previous page
	if vidList.Page > 1 {
		vidList.PrevPage = vidList.Page - 1
	} else {
		vidList.PrevPage = 1
	}

	vidList.Items = videos

	return vidList, nil
}

func (d *Database) GetVideoIDs() ([]string, error) {
	rows, err := d.Client.Query("SELECT id FROM videos")
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("Failed to get video IDs: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error(err)
		}
	}()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Error(err)
			return nil, fmt.Errorf("Failed to scan video ID: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (d *Database) ScannerCreateVideo(vid video.Video) (video.Video, error) {
	vidRow := &VideoRow{
		ID:            vid.ID,
		Channel:       vid.Channel,
		ChannelID:     vid.ChannelID,
		Title:         vid.Title,
		Description:   sql.NullString{String: vid.Description, Valid: true},
		UploadDate:    vid.UploadDate,
		Uploader:      sql.NullString{String: vid.Uploader, Valid: true},
		Duration:      vid.Duration,
		ViewCount:     vid.ViewCount,
		LikeCount:     sql.NullInt64{Int64: vid.LikeCount, Valid: true},
		DislikeCount:  sql.NullInt64{Int64: vid.DislikeCount, Valid: true},
		Format:        sql.NullString{String: vid.Format, Valid: true},
		Width:         sql.NullInt64{Int64: vid.Width, Valid: true},
		Height:        sql.NullInt64{Int64: vid.Height, Valid: true},
		Resolution:    sql.NullString{String: vid.Resolution, Valid: true},
		FPS:           sql.NullFloat64{Float64: vid.FPS, Valid: true},
		VideoCodec:    sql.NullString{String: vid.VideoCodec, Valid: true},
		VBR:           sql.NullFloat64{Float64: vid.VBR, Valid: true},
		AudioCodec:    sql.NullString{String: vid.AudioCodec, Valid: true},
		ABR:           sql.NullFloat64{Float64: vid.ABR, Valid: true},
		Epoch:         sql.NullInt64{Int64: vid.Epoch, Valid: true},
		CommentCount:  vid.CommentCount,
		Tags:          sql.NullString{String: vid.Tags, Valid: true},
		Categories:    sql.NullString{String: vid.Categories, Valid: true},
		VideoPath:     vid.VideoPath,
		ThumbnailPath: vid.ThumbnailPath,
		JsonPath:      sql.NullString{String: vid.JsonPath, Valid: true},
		SubtitlePath:  sql.NullString{String: vid.SubtitlePath, Valid: true},
		CreatedAt:     vid.CreatedAt,
		UpdatedAt:     vid.UpdatedAt,
	}
	// Set the default values
	vidRow.CreatedAt = time.Now()
	vidRow.UpdatedAt = time.Now()

	rows, err := d.Client.NamedQuery("INSERT INTO videos (id, channel, channel_id, title, description, upload_date, uploader, duration, view_count, like_count, dislike_count, format, width, height, resolution, fps, video_codec, vbr, audio_codec, abr, epoch, comment_count, tags, categories, video_path, thumbnail_path, json_path, subtitle_path, created_at, updated_at) VALUES (:id, :channel, :channelid, :title, :description, :uploaddate, :uploader, :duration, :viewcount, :likecount, :dislikecount, :format, :width, :height, :resolution, :fps, :videocodec, :vbr, :audiocodec, :abr, :epoch, :commentcount, :tags, :categories, :videopath, :thumbnailpath, :jsonpath, :subtitlepath, :createdat, :updatedat)", vidRow)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Errorf("Error creating video: %s: %s", err.Code, err.Message)
			switch err.Code {
			case "23505":
				return video.Video{}, echo.NewHTTPError(http.StatusConflict, "Video already exists")
			case "23503":
				return video.Video{}, echo.NewHTTPError(http.StatusNotFound, "Channel ID does not exist")
			default:
				log.Error(err)
				return video.Video{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
			}
		}
	}
	if err := rows.Close(); err != nil {
		return video.Video{}, fmt.Errorf("Failed to close rows: %w", err)
	}

	return vid, nil

}

func (d *Database) ScannerCreateChapters(chaps []video.Chapter) error {
	for _, chap := range chaps {
		chapRow := ChapterRow{
			ID:        chap.ID,
			StartTime: chap.StartTime,
			Title:     chap.Title,
			EndTime:   chap.EndTime,
			VideoID:   chap.VideoID,
		}

		chapRow.ID = uuid.New().String()

		rows, err := d.Client.NamedQuery("INSERT INTO chapters (id, start_time, title, end_time, video_id) VALUES (:id, :starttime, :title, :endtime, :videoid)", &chapRow)

		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				switch err.Code {
				case "23505":
					return echo.NewHTTPError(http.StatusConflict, "Chapter already exists")
				case "23503":
					return echo.NewHTTPError(http.StatusNotFound, "Video ID does not exist")
				default:
					log.Errorf("Failed to insert comment: %v", err)
					return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
				}
			}
		}

		if err := rows.Close(); err != nil {
			return fmt.Errorf("Failed to close rows: %v", err)
		}

	}
	return nil
}

func (d *Database) GetRandomVideos(c echo.Context, count string) ([]video.Video, error) {
	rows, err := d.Client.QueryContext(c.Request().Context(), "SELECT * FROM videos OFFSET floor(random() * ( SELECT COUNT(*) FROM videos)) LIMIT $1", count)
	if err != nil {
		log.Error(err)
		return []video.Video{}, fmt.Errorf("Failed to get videos: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error(err)
		}
	}()
	var videos []video.Video
	for rows.Next() {
		var vid video.Video
		if err := rows.Scan(&vid.ID, &vid.Channel, &vid.ChannelID, &vid.Title, &vid.Description, &vid.UploadDate, &vid.Uploader, &vid.Duration, &vid.ViewCount, &vid.LikeCount, &vid.DislikeCount, &vid.Format, &vid.Width, &vid.Height, &vid.Resolution, &vid.FPS, &vid.VideoCodec, &vid.VBR, &vid.AudioCodec, &vid.ABR, &vid.Epoch, &vid.CommentCount, &vid.Tags, &vid.Categories, &vid.VideoPath, &vid.ThumbnailPath, &vid.JsonPath, &vid.SubtitlePath, &vid.CreatedAt, &vid.UpdatedAt); err != nil {
			log.Error(err)
			return []video.Video{}, fmt.Errorf("Failed to scan video: %w", err)
		}

		videos = append(videos, vid)
	}

	return videos, nil
}
func (d *Database) SearchVideos(c echo.Context, query string, vidList video.VideoList) (video.VideoList, error) {
	skip := (vidList.Page - 1) * vidList.Limit

	rows, err := d.Client.QueryContext(c.Request().Context(), "SELECT *, COUNT(*) OVER() AS TotalCount FROM videos WHERE title ILIKE $1 LIMIT $2 OFFSET $3", "%"+query+"%", vidList.Limit, skip)
	if err != nil {
		log.Error(err)
		return video.VideoList{}, fmt.Errorf("Failed to search videos: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error(err)
		}
	}()

	var rowCount int
	var videos []video.Video
	for rows.Next() {
		var vid video.Video
		if err := rows.Scan(&vid.ID, &vid.Channel, &vid.ChannelID, &vid.Title, &vid.Description, &vid.UploadDate, &vid.Uploader, &vid.Duration, &vid.ViewCount, &vid.LikeCount, &vid.DislikeCount, &vid.Format, &vid.Width, &vid.Height, &vid.Resolution, &vid.FPS, &vid.VideoCodec, &vid.VBR, &vid.AudioCodec, &vid.ABR, &vid.Epoch, &vid.CommentCount, &vid.Tags, &vid.Categories, &vid.VideoPath, &vid.ThumbnailPath, &vid.JsonPath, &vid.SubtitlePath, &vid.CreatedAt, &vid.UpdatedAt, &rowCount); err != nil {
			log.Error(err)
			return video.VideoList{}, fmt.Errorf("Failed to scan video: %w", err)
		}

		videos = append(videos, vid)
	}

	// Set pagination info
	vidList.LastPage = int(math.Ceil(float64(rowCount) / float64(vidList.Limit)))
	vidList.TotalItems = rowCount

	// Set next page
	if vidList.Page < vidList.LastPage {
		vidList.NextPage = vidList.Page + 1
	} else {
		vidList.NextPage = vidList.LastPage
	}

	// Set previous page
	if vidList.Page > 1 {
		vidList.PrevPage = vidList.Page - 1
	} else {
		vidList.PrevPage = 1
	}

	vidList.Items = videos

	return vidList, nil
}

func (d *Database) ScannerGetChannelVideoCount(id string) (int, error) {
	var count int
	err := d.Client.QueryRow("SELECT COUNT(*) FROM videos WHERE channel_id = $1", id).Scan(&count)
	if err != nil {
		log.Error(err)
		return 0, fmt.Errorf("Failed to get video count: %w", err)
	}

	return count, nil
}
