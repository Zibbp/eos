package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/zibbp/eos/internal/video"
)

type VideoService interface {
	CreateVideo(c echo.Context, vid video.Video) (video.Video, error)
	GetChannelVideos(c echo.Context, channelID string, vidList video.VideoList) (video.VideoList, error)
	GetVideo(c echo.Context, vidID string) (video.SingleVideoList, error)
	GetRandomVideos(c echo.Context, count string) ([]video.Video, error)
	SearchVideos(c echo.Context, query string, vidList video.VideoList) (video.VideoList, error)
}

type CreateVideo struct {
	ID            string    `json:"id" validate:"required"`
	Title         string    `json:"title" validate:"required"`
	Description   string    `json:"description"`
	UploadDate    time.Time `json:"upload_date" validate:"required"`
	Uploader      string    `json:"uploader"`
	Duration      uint64    `json:"duration" validate:"required"`
	ViewCount     uint64    `json:"view_count" validate:"required"`
	LikeCount     int64     `json:"like_count" validate:"required"`
	DislikeCount  int64     `json:"dislike_count"`
	Channel       string    `json:"channel" validate:"required"`
	ChannelID     string    `json:"channel_id" validate:"required"`
	Format        string    `json:"format"`
	Width         int64     `json:"width"`
	Height        int64     `json:"height"`
	Resolution    string    `json:"resolution"`
	FPS           float64   `json:"fps"`
	VideoCodec    string    `json:"video_codec"`
	VBR           float64   `json:"vbr"`
	AudioCodec    string    `json:"audio_codec"`
	ABR           float64   `json:"abr"`
	Epoch         int64     `json:"epoch"`
	CommentCount  uint64    `json:"comment_count" validate:"required"`
	Categories    string    `json:"categories"`
	Tags          string    `json:"tags"`
	VideoPath     string    `json:"video_path" validate:"required"`
	ThumbnailPath string    `json:"thumbnail_path" validate:"required"`
	JsonPath      string    `json:"json_path"`
	SubtitlePath  string    `json:"subtitle_path"`
}

func videoPostToVideo(v CreateVideo) video.Video {
	return video.Video{
		ID:            v.ID,
		Title:         v.Title,
		Description:   v.Description,
		UploadDate:    v.UploadDate,
		Uploader:      v.Uploader,
		Duration:      v.Duration,
		ViewCount:     v.ViewCount,
		LikeCount:     v.LikeCount,
		DislikeCount:  v.DislikeCount,
		Channel:       v.Channel,
		ChannelID:     v.ChannelID,
		Format:        v.Format,
		Width:         v.Width,
		Height:        v.Height,
		Resolution:    v.Resolution,
		FPS:           v.FPS,
		VideoCodec:    v.VideoCodec,
		VBR:           v.VBR,
		AudioCodec:    v.AudioCodec,
		ABR:           v.ABR,
		Epoch:         v.Epoch,
		CommentCount:  v.CommentCount,
		Tags:          v.Tags,
		Categories:    v.Categories,
		VideoPath:     v.VideoPath,
		ThumbnailPath: v.ThumbnailPath,
		JsonPath:      v.JsonPath,
		SubtitlePath:  v.SubtitlePath,
	}
}

func (h *Handler) GetVideo(c echo.Context) error {
	vidID := c.Param("vid_id")
	vid, err := h.Service.VideoService.GetVideo(c, vidID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, vid)
}

func (h *Handler) CreateVideo(c echo.Context) error {
	var createVideo CreateVideo
	if err := c.Bind(&createVideo); err != nil {
		log.Error("CreateVideo", err)
		return c.JSON(http.StatusBadRequest, "invalid request")
	}

	validate := validator.New()
	err := validate.Struct(createVideo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	vid := videoPostToVideo(createVideo)
	vid, err = h.Service.VideoService.CreateVideo(c, vid)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, vid)
}

func (h *Handler) GetChannelVideos(c echo.Context) error {

	// time.Sleep(2 * time.Second)

	// ToDo: Move to middleware
	var getLimit, getPage string
	var limit, page int
	if c.QueryParam("limit") != "" {
		getLimit = c.QueryParam("limit")
	} else {
		getLimit = "10"
	}
	if c.QueryParam("page") != "" {
		getPage = c.QueryParam("page")
	} else {
		getPage = "1"
	}
	limit, err := strconv.Atoi(getLimit)
	if err != nil {
	}
	page, err = strconv.Atoi(getPage)
	if err != nil {
		page = 1
	}

	vidList := video.VideoList{
		Limit: limit,
		Page:  page,
	}

	channelID := c.Param("channel_id")
	videos, err := h.Service.VideoService.GetChannelVideos(c, channelID, vidList)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, videos)
}

func (h *Handler) GetRandomVideos(c echo.Context) error {
	count := c.QueryParam("count")
	if count == "" {
		count = "10"
	}
	videos, err := h.Service.VideoService.GetRandomVideos(c, count)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, videos)
}

func (h *Handler) SearchVideos(c echo.Context) error {

	var getLimit, getPage string
	var limit, page int
	if c.QueryParam("limit") != "" {
		getLimit = c.QueryParam("limit")
	} else {
		getLimit = "10"
	}
	if c.QueryParam("page") != "" {
		getPage = c.QueryParam("page")
	} else {
		getPage = "1"
	}
	limit, err := strconv.Atoi(getLimit)
	if err != nil {
	}
	page, err = strconv.Atoi(getPage)
	if err != nil {
		page = 1
	}

	vidList := video.VideoList{
		Limit: limit,
		Page:  page,
	}

	searchQuery := c.QueryParam("query")
	if searchQuery == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request. Query parameter is required.")
	}
	videos, err := h.Service.VideoService.SearchVideos(c, searchQuery, vidList)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, videos)
}
