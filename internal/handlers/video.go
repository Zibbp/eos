package handlers

import (
	"context"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zibbp/eos/internal/video"
	components "github.com/zibbp/eos/internal/views/components/video"
	"github.com/zibbp/eos/internal/views/pages"
)

type VideoService interface {
	GetVideosFilter(ctx context.Context, filter video.VideoFilter) ([]video.Video, int, error)
	GetVideoByExtId(ctx context.Context, extVideoId string) (*video.Video, error)
	FtsVideosFilter(ctx context.Context, filter video.FtsVideoFilter) ([]video.VideoSearchResult, int, error)
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

	return render(c, pages.VideoPage(*video, *channel))
}

func (h *Handler) HandleVideoSearchPage(c echo.Context) error {
	searchQuery := c.QueryParam("q")

	if searchQuery == "" {
		return c.JSON(500, "Search parameter 'q' required")
	}

	limit := 20

	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	videos, totalVideos, err := h.Services.VideoService.FtsVideosFilter(c.Request().Context(), video.FtsVideoFilter{
		Term:   searchQuery,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return c.JSON(500, err)
	}

	totalPages := (totalVideos + limit - 1) / limit

	if c.Request().Header.Get("HX-Request") == "true" && c.Request().Header.Get("HX-Boosted") == "" {
		return render(c, components.VideoSearchList(searchQuery, videos, page, totalPages, true))
	} else {
		return render(c, pages.SearchPage(searchQuery, videos, page, totalPages))
	}

}
