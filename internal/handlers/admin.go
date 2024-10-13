package handlers

import (
	"github.com/labstack/echo/v4"
	pages "github.com/zibbp/eos/internal/views/pages/admin"
)

func (h *Handler) HandleBlockedPathsPage(c echo.Context) error {

	blockedPaths, err := h.Services.BlockedPaths.GetBlockedPaths(c.Request().Context())
	if err != nil {
		return c.JSON(500, err)
	}

	return render(c, pages.AdminBlockedPaths(blockedPaths))

}
