package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/zibbp/eos/internal/views/pages"
)

func (h *Handler) HandleLandingIndex(c echo.Context) error {
	return render(c, pages.Index())
}
