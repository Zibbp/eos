package database

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
	"github.com/zibbp/eos/internal/comment"
)

type CommentRow struct {
	ID               string         `json:"id"`
	Text             string         `json:"text"`
	Timestamp        uint64         `json:"timestamp"`
	LikeCount        uint64         `json:"like_count"`
	IsFavorited      sql.NullBool   `json:"is_favorited"`
	Author           string         `json:"author"`
	AuthorID         string         `json:"author_id"`
	AuthorThumbnail  sql.NullString `json:"author_thumbnail"`
	AuthorIsUploader sql.NullBool   `json:"author_is_uploader"`
	Parent           string         `json:"parent"`
	VideoID          string         `json:"video_id"`
}

func paginate(x []comment.Comment, cmtList *comment.CommentList) []comment.Comment {

	skip := (cmtList.Page - 1) * cmtList.Limit

	// Set last page
	cmtList.LastPage = int(math.Ceil(float64(len(x)) / float64(cmtList.Limit)))

	// Set next page
	if cmtList.Page < cmtList.LastPage {
		cmtList.NextPage = cmtList.Page + 1
	} else {
		cmtList.NextPage = cmtList.LastPage
	}

	// Set previous page
	if cmtList.Page > 1 {
		cmtList.PrevPage = cmtList.Page - 1
	} else {
		cmtList.PrevPage = 1
	}

	if skip > len(x) {
		skip = len(x)
	}

	end := skip + cmtList.Limit
	if end > len(x) {
		end = len(x)
	}
	return x[skip:end]
}

func (d *Database) GetVideoComments(c echo.Context, vidID string, cmtList comment.CommentList) (comment.CommentList, error) {
	// Fetch video to ensure it exists
	_, err := d.GetVideo(c, vidID)
	if err != nil {
		return cmtList, err
	}

	rows, err := d.Client.QueryContext(c.Request().Context(), "SELECT id, text, timestamp, like_count, is_favorited, author, author_id, author_thumbnail, author_is_uploader, parent, video_id FROM comments WHERE video_id = $1 ORDER BY like_count DESC", vidID)
	if err != nil {
		log.Errorf("Failed to get video comments: %v", err)
		return cmtList, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Errorf("Failed to close rows: %v", err)
		}
	}()

	var comments []comment.Comment
	for rows.Next() {
		var comment comment.Comment
		if err := rows.Scan(&comment.ID, &comment.Text, &comment.Timestamp, &comment.LikeCount, &comment.IsFavorited, &comment.Author, &comment.AuthorID, &comment.AuthorThumbnail, &comment.AuthorIsUploader, &comment.Parent, &comment.VideoID); err != nil {
			log.Errorf("Failed to scan comment: %v", err)
			return cmtList, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
		comments = append(comments, comment)
	}

	// add child comment to parent comment
	for i := range comments {
		if comments[i].Parent != "root" {
			for j := range comments {
				if comments[j].ID == comments[i].Parent {
					comments[j].Replies = append(comments[j].Replies, comments[i])
					break
				}
			}
		}
	}

	// Cannot paginate with sql as the child comments need to be added
	// Use a less efficient method to paginate (but it works!)
	cmtList.Items = paginate(comments, &cmtList)
	cmtList.TotalItems = len(comments)

	return cmtList, nil

}

func (d *Database) CreateComment(c echo.Context, cmt comment.Comment) (comment.Comment, error) {
	cmtRow := CommentRow{
		ID:               cmt.ID,
		Text:             cmt.Text,
		Timestamp:        cmt.Timestamp,
		LikeCount:        cmt.LikeCount,
		IsFavorited:      sql.NullBool{Bool: cmt.IsFavorited, Valid: true},
		Author:           cmt.Author,
		AuthorID:         cmt.AuthorID,
		AuthorThumbnail:  sql.NullString{String: cmt.AuthorThumbnail, Valid: true},
		AuthorIsUploader: sql.NullBool{Bool: cmt.AuthorIsUploader, Valid: true},
		Parent:           cmt.Parent,
		VideoID:          cmt.VideoID,
	}

	rows, err := d.Client.NamedQueryContext(c.Request().Context(), "INSERT INTO comments (id, text, timestamp, like_count, is_favorited, author, author_id, author_thumbnail, author_is_uploader, parent, video_id) VALUES (:id, :text, :timestamp, :likecount, :isfavorited, :author, :authorid, :authorthumbnail, :authorisuploader, :parent, :videoid)", &cmtRow)

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			switch err.Code {
			case "23505":
				return comment.Comment{}, echo.NewHTTPError(http.StatusConflict, "Comment already exists")
			case "23503":
				return comment.Comment{}, echo.NewHTTPError(http.StatusNotFound, "Video ID does not exist")
			default:
				return comment.Comment{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
			}
		}
	}

	if err := rows.Close(); err != nil {
		return comment.Comment{}, fmt.Errorf("Failed to close rows: %v", err)
	}

	return cmt, nil

}

func (d *Database) ScannerCreateComments(cmts []comment.Comment) error {
	for _, cmt := range cmts {
		cmtRow := CommentRow{
			ID:               cmt.ID,
			Text:             cmt.Text,
			Timestamp:        cmt.Timestamp,
			LikeCount:        cmt.LikeCount,
			IsFavorited:      sql.NullBool{Bool: cmt.IsFavorited, Valid: true},
			Author:           cmt.Author,
			AuthorID:         cmt.AuthorID,
			AuthorThumbnail:  sql.NullString{String: cmt.AuthorThumbnail, Valid: true},
			AuthorIsUploader: sql.NullBool{Bool: cmt.AuthorIsUploader, Valid: true},
			Parent:           cmt.Parent,
			VideoID:          cmt.VideoID,
		}

		rows, err := d.Client.NamedQuery("INSERT INTO comments (id, text, timestamp, like_count, is_favorited, author, author_id, author_thumbnail, author_is_uploader, parent, video_id) VALUES (:id, :text, :timestamp, :likecount, :isfavorited, :author, :authorid, :authorthumbnail, :authorisuploader, :parent, :videoid)", &cmtRow)

		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				switch err.Code {
				case "23505":
					return echo.NewHTTPError(http.StatusConflict, "Comment already exists")
				case "23503":
					return echo.NewHTTPError(http.StatusNotFound, "Video ID does not exist")
				default:
					log.Errorf("Failed to insert comment: %v", err)
					return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
				}
			}
		}

		if err := rows.Close(); err != nil {
			return fmt.Errorf("Failed to close rows: %v", err)
		}

	}
	return nil
}
