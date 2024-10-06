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

	return c.JSON(200, "queued scan task")
}
