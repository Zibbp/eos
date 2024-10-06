package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/zibbp/avalon/internal/comment"
	components "github.com/zibbp/avalon/internal/views/components/video"
)

type CommentService interface {
	GetRootCommentsForVideo(ctx context.Context, videoId uuid.UUID, limit int, offset int) ([]comment.Comment, error)
	GetCommentReplies(ctx context.Context, commentId string) ([]comment.Comment, error)
}

func (h *Handler) HandleVideoCommentsPage(c echo.Context) error {

	videoIdStr := c.Param("video_id")

	video, err := h.Services.VideoService.GetVideoByExtId(c.Request().Context(), videoIdStr)
	if err != nil {
		return c.JSON(500, err)
	}

	comments, err := h.Services.CommentService.GetRootCommentsForVideo(c.Request().Context(), video.ID, 25, 0)
	if err != nil {
		return c.JSON(500, err)
	}

	return render(c, components.VideoComments(comments))
}

func (h *Handler) HandleVideoCommentReplies(c echo.Context) error {

	commentId := c.Param("comment_id")

	comments, err := h.Services.CommentService.GetCommentReplies(c.Request().Context(), commentId)
	if err != nil {
		return c.JSON(500, err)
	}

	return render(c, components.VideoComments(comments))
}
