package video

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/zibbp/avalon/internal/db/sqlc"
)

type VideoFilter struct {
	ChannelID uuid.UUID
	Limit     int
	Offset    int
}

type VideoService interface {
	GetVideosFilter(ctx context.Context, filter VideoFilter) ([]Video, int, error)
	GetVideoByExtId(ctx context.Context, extVideoId string) (*Video, error)
}

type Service struct {
	Store db.Store
}

func NewService(store db.Store) VideoService {
	return &Service{
		Store: store,
	}
}

type Video struct {
	ID             uuid.UUID
	ExtID          string
	Title          string
	Description    string
	UploadDate     time.Time
	Uploader       string
	Duration       int
	ViewCount      int
	LikeCount      int
	DislikeCount   int
	Format         string
	Height         int
	Width          int
	Resolution     string
	Fps            float32
	VideoCodec     string
	Vbr            float32
	AudioCodec     string
	Abr            float32
	CommentCount   int32
	VideoPath      string
	ThumbnailPath  string
	InfoPath       string
	SubtitlePath   []string
	Path           string
	StoryboardPath string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ChannelID      uuid.UUID
}

func (s *Service) GetVideosFilter(ctx context.Context, filter VideoFilter) ([]Video, int, error) {

	dbVideos, err := s.Store.GetVideosFilter(ctx, db.GetVideosFilterParams{
		ChannelID: pgtype.UUID{Bytes: filter.ChannelID, Valid: true},
		Limit:     int32(filter.Limit),
		Offset:    int32(filter.Offset),
	})
	if err != nil {
		return nil, 0, err
	}

	videoCount, err := s.Store.GetVideosCount(ctx, pgtype.UUID{Bytes: filter.ChannelID, Valid: true})
	if err != nil {
		return nil, 0, err
	}

	var videos []Video
	for _, dbVideo := range dbVideos {
		videos = append(videos, dbVideoToVideo(dbVideo))
	}

	return videos, int(videoCount), nil
}

func (s *Service) GetVideoByExtId(ctx context.Context, extVideoId string) (*Video, error) {
	dbVideo, err := s.Store.GetVideoByExternalID(ctx, &extVideoId)
	if err != nil {
		return nil, err
	}
	video := dbVideoToVideo(dbVideo)
	return &video, nil
}

func dbVideoToVideo(dbVideo db.Video) Video {
	return Video{
		ID:            dbVideo.ID.Bytes,
		ExtID:         *dbVideo.ExtID,
		Title:         dbVideo.Title,
		Description:   *dbVideo.Description,
		UploadDate:    dbVideo.UploadDate.Time,
		Uploader:      *dbVideo.Uploader,
		Duration:      int(dbVideo.Duration),
		ViewCount:     int(dbVideo.ViewCount),
		LikeCount:     int(*dbVideo.LikeCount),
		DislikeCount:  int(*dbVideo.DislikeCount),
		Format:        *dbVideo.Format,
		Height:        int(*dbVideo.Height),
		Width:         int(*dbVideo.Width),
		Resolution:    *dbVideo.Resolution,
		Fps:           *dbVideo.Fps,
		VideoCodec:    *dbVideo.VideoCodec,
		Vbr:           *dbVideo.Vbr,
		AudioCodec:    *dbVideo.AudioCodec,
		Abr:           *dbVideo.Abr,
		CommentCount:  *dbVideo.CommentCount,
		VideoPath:     dbVideo.VideoPath,
		ThumbnailPath: dbVideo.ThumbnailPath,
		InfoPath:      dbVideo.InfoPath,
		SubtitlePath:  dbVideo.SubtitlePath,
		Path:          dbVideo.Path,
		CreatedAt:     dbVideo.CreatedAt.Time,
		UpdatedAt:     dbVideo.UpdatedAt.Time,
		ChannelID:     dbVideo.ChannelID.Bytes,
	}
}
