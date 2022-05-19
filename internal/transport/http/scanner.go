package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ScannerService interface {
	StartScanner(c echo.Context) error
}

func (h *Handler) StartScanner(c echo.Context) error {
	err := h.Service.ScannerService.StartScanner(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "Scanner started")
}
