package video

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zibbp/eos/internal/channel"
)

type Store interface {
	CreateVideo(c echo.Context, vid Video) (Video, error)
	GetChannelVideos(c echo.Context, channelID string, vidList VideoList) (VideoList, error)
	GetVideo(c echo.Context, vidID string) (SingleVideoList, error)
	GetRandomVideos(c echo.Context, count string) ([]Video, error)
	SearchVideos(c echo.Context, query string, vidList VideoList) (VideoList, error)
}

type Service struct {
	Store Store
}

func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

type Video struct {
	ID            string    `json:"id"`
	Channel       string    `json:"channel"`
	ChannelID     string    `json:"channel_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	UploadDate    time.Time `json:"upload_date"`
	Uploader      string    `json:"uploader"`
	Duration      uint64    `json:"duration"`
	ViewCount     uint64    `json:"view_count"`
	LikeCount     int64     `json:"like_count"`
	DislikeCount  int64     `json:"dislike_count"`
	Format        string    `json:"format"`
	Width         int64     `json:"width"`
	Height        int64     `json:"height"`
	Resolution    string    `json:"resolution"`
	FPS           float64   `json:"fps"`
	VideoCodec    string    `json:"vcodec"`
	VBR           float64   `json:"vbr"`
	AudioCodec    string    `json:"acodec"`
	ABR           float64   `json:"abr"`
	Epoch         int64     `json:"epoch"`
	CommentCount  uint64    `json:"comment_count"`
	Tags          string    `json:"tags"`
	Categories    string    `json:"categories"`
	VideoPath     string    `json:"video_path"`
	ThumbnailPath string    `json:"thumbnail_path"`
	JsonPath      string    `json:"json_path"`
	SubtitlePath  string    `json:"subtitle_path"`
	Chapters      []Chapter `json:"chapters"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Chapter struct {
	ID        string  `json:"id"`
	StartTime float64 `json:"start_time"`
	Title     string  `json:"title"`
	EndTime   float64 `json:"end_time"`
	VideoID   string  `json:"video_id"`
}

type SingleVideoList struct {
	Video   Video           `json:"video"`
	Channel channel.Channel `json:"channel"`
}

type VideoList struct {
	Limit      int     `json:"limit"`
	Page       int     `json:"page"`
	PrevPage   int     `json:"prev_page"`
	NextPage   int     `json:"next_page"`
	LastPage   int     `json:"last_page"`
	TotalItems int     `json:"total_items"`
	Items      []Video `json:"items"`
}

func (s *Service) GetVideo(c echo.Context, vidID string) (SingleVideoList, error) {
	vid, err := s.Store.GetVideo(c, vidID)
	if err != nil {
		return vid, err
	}
	return vid, nil
}

func (s *Service) CreateVideo(c echo.Context, vid Video) (Video, error) {
	vid, err := s.Store.CreateVideo(c, vid)
	if err != nil {
		return vid, err
	}
	return vid, nil
}

func (s *Service) GetChannelVideos(c echo.Context, channelID string, vidList VideoList) (VideoList, error) {
	videos, err := s.Store.GetChannelVideos(c, channelID, vidList)
	if err != nil {
		return VideoList{}, err
	}
	return videos, nil
}
func (s *Service) GetRandomVideos(c echo.Context, count string) ([]Video, error) {
	videos, err := s.Store.GetRandomVideos(c, count)
	if err != nil {
		return []Video{}, err
	}
	return videos, nil
}
func (s *Service) SearchVideos(c echo.Context, query string, vidList VideoList) (VideoList, error) {
	videos, err := s.Store.SearchVideos(c, query, vidList)
	if err != nil {
		return VideoList{}, err
	}
	return videos, nil
}
