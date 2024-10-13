package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/zibbp/eos/internal/blocked_paths"
)

type BlockedPaths interface {
	GetBlockedPaths(ctx context.Context) ([]blocked_paths.BlockedPath, error)
	DeleteBlockedPathById(ctx context.Context, id uuid.UUID) error
}

func (h *Handler) DeleteBlockedPathById(c echo.Context) error {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(500, err)
	}

	err = h.Services.BlockedPaths.DeleteBlockedPathById(c.Request().Context(), uuid)
	if err != nil {
		return c.JSON(500, err)
	}

	return c.NoContent(200)

}
