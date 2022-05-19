package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/zibbp/eos/internal/comment"
)

type CommentService interface {
	CreateComment(c echo.Context, cmt comment.Comment) (comment.Comment, error)
	GetVideoComments(c echo.Context, vidID string, cmtList comment.CommentList) (comment.CommentList, error)
}

type CreateComment struct {
	ID               string `json:"id" validate:"required"`
	Text             string `json:"text" validate:"required"`
	Timestamp        uint64 `json:"timestamp" validate:"required"`
	LikeCount        uint64 `json:"like_count" validate:"required"`
	IsFavorited      bool   `json:"is_favorited"`
	Author           string `json:"author" validate:"required"`
	AuthorID         string `json:"author_id" validate:"required"`
	AuthorThumbnail  string `json:"author_thumbnail"`
	AuthorIsUploader bool   `json:"author_is_uploader"`
	Parent           string `json:"parent" validate:"required"`
	VideoID          string `json:"video_id" validate:"required"`
}

func postCommentToComment(c CreateComment) comment.Comment {
	return comment.Comment{
		ID:               c.ID,
		Text:             c.Text,
		Timestamp:        c.Timestamp,
		LikeCount:        c.LikeCount,
		IsFavorited:      c.IsFavorited,
		Author:           c.Author,
		AuthorID:         c.AuthorID,
		AuthorThumbnail:  c.AuthorThumbnail,
		AuthorIsUploader: c.AuthorIsUploader,
		Parent:           c.Parent,
		VideoID:          c.VideoID,
	}
}

func (h *Handler) GetVideoComments(c echo.Context) error {
	videoID := c.Param("vid_id")

	// ToDo: Move to middleware
	var getLimit, getPage string
	var limit, page int
	if c.QueryParam("limit") != "" {
		getLimit = c.QueryParam("limit")
	} else {
		fmt.Println("missing limit")
		getLimit = "10"
	}
	if c.QueryParam("page") != "" {
		getPage = c.QueryParam("page")
	} else {
		fmt.Println("missing page")
		getPage = "1"
	}
	limit, err := strconv.Atoi(getLimit)
	if err != nil {
		fmt.Println("limit is not a number")
	}
	page, err = strconv.Atoi(getPage)
	if err != nil {
		fmt.Println("page is not a number")
		page = 1
	}

	cmtList := comment.CommentList{
		Limit: limit,
		Page:  page,
	}

	cmts, err := h.Service.CommentService.GetVideoComments(c, videoID, cmtList)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, cmts)
}

func (h *Handler) CreateComment(c echo.Context) error {
	var createComment CreateComment
	if err := c.Bind(&createComment); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	err := validate.Struct(createComment)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	cmt := postCommentToComment(createComment)
	cmt, err = h.Service.CommentService.CreateComment(c, cmt)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, cmt)
}

// func (h *Handler) CreateComment(c echo.Context) error {
// 	var postCmt CreateComment
// 	if err := c.Bind(&postCmt); err != nil {
// 		return err
// 	}

// 	validate := validator.New()
// 	err := validate.Struct(postCmt)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, err.Error())
// 	}

// 	cmt := postCommentToComment(&postCmt)
// 	cmt, err = h.Service.CommentService.CreateComment(c, cmt)
// 	if err != nil {
// 		return err
// 	}

// 	return c.JSON(http.StatusOK, cmt)
// }
