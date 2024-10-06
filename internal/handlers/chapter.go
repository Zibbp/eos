package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/zibbp/eos/internal/chapter"
)

type ChapterService interface {
	GetChaptersForVideo(ctx context.Context, videoId uuid.UUID) ([]chapter.Chapter, error)
	GetChaptersForVidstackPlayer(ctx context.Context, videoId uuid.UUID) ([]chapter.VidstackPlayerChapter, error)
}

func (h *Handler) GetChaptersForVideo(c echo.Context) error {
	videoExtId := c.Param("video_id")

	video, err := h.Services.VideoService.GetVideoByExtId(c.Request().Context(), videoExtId)
	if err != nil {
		return c.JSON(500, err)
	}

	chapters, err := h.Services.ChapterService.GetChaptersForVidstackPlayer(c.Request().Context(), video.ID)
	if err != nil {
		return c.JSON(500, err)
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(200, chapters)
}
