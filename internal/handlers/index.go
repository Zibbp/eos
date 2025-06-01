package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/zibbp/eos/internal/views/pages"
)

func (h *Handler) HandleLandingIndex(c echo.Context) error {
	channels, err := h.Services.ChannelService.GetChannels(c.Request().Context())
	if err != nil {
		return c.JSON(500, err)
	}

	return render(c, pages.Index(channels))
}
