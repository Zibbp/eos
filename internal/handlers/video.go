package handlers

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/zibbp/eos/internal/video"
	"github.com/zibbp/eos/internal/views/pages"
)

type VideoService interface {
	GetVideosFilter(ctx context.Context, filter video.VideoFilter) ([]video.Video, int, error)
	GetVideoByExtId(ctx context.Context, extVideoId string) (*video.Video, error)
}

func (h *Handler) HandelVideoPage(c echo.Context) error {
	videoId := c.Param("video_id")

	video, err := h.Services.VideoService.GetVideoByExtId(c.Request().Context(), videoId)
	if err != nil {
		return c.JSON(500, err)
	}

	channel, err := h.Services.ChannelService.GetChannelByID(c.Request().Context(), video.ChannelID)
	if err != nil {
		return c.JSON(500, err)
	}

	return render(c, pages.VideoPage(h.Config.CDN_URL, *video, *channel))
}
