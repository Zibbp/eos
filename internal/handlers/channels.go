package handlers

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	db "github.com/zibbp/avalon/internal/db/sqlc"
	"github.com/zibbp/avalon/internal/video"
	components "github.com/zibbp/avalon/internal/views/components/video"
	"github.com/zibbp/avalon/internal/views/pages"
)

type ChannelService interface {
	GetChannels(ctx context.Context) ([]db.Channel, error)
	GetChannelByName(ctx context.Context, name string) (*db.Channel, error)
	GetChannelByID(ctx context.Context, id uuid.UUID) (*db.Channel, error)
}

func (h *Handler) HandleChannelsPage(c echo.Context) error {

	channels, err := h.Services.ChannelService.GetChannels(c.Request().Context())
	if err != nil {
		return c.JSON(500, err)
	}

	return render(c, pages.Channels(channels))
}

func (h *Handler) HandleChannelPage(c echo.Context) error {
	name := c.Param("name")

	cleanName, err := url.QueryUnescape(name)
	if err != nil {
		return c.JSON(500, fmt.Errorf("error unescaping name"))
	}

	channel, err := h.Services.ChannelService.GetChannelByName(c.Request().Context(), cleanName)
	if err != nil {
		return c.JSON(500, err)
	}

	limit := 20

	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	videos, totalVideos, err := h.Services.VideoService.GetVideosFilter(c.Request().Context(), video.VideoFilter{
		ChannelID: channel.ID.Bytes,
		Limit:     limit,
		Offset:    offset,
	})

	if err != nil {
		return c.JSON(500, err)
	}

	totalPages := (totalVideos + limit - 1) / limit

	if c.Request().Header.Get("HX-Request") == "true" && c.Request().Header.Get("HX-Boosted") == "" {
		return render(c, components.VideoList(channel, videos, page, totalPages, true))
	} else {
		return render(c, pages.ChannelName(channel, videos, page, totalPages))
	}
}
