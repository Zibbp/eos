package handlers

import (
	"github.com/labstack/echo/v4"
)

type ScannerService interface {
	Scan() error
}

func (h *Handler) StartScanner(c echo.Context) error {

	err := h.Services.ScannerService.Scan()
	if err != nil {
		return c.JSON(500, err)
	}

	if c.Request().Header.Get("HX-Request") == "true" && c.Request().Header.Get("HX-Boosted") == "" {
		return c.String(200, "✅ Queued Scan Task")
	} else {
		return c.JSON(200, "queued scan task")
	}
}
