package comment

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/zibbp/avalon/internal/db/sqlc"
)

type CommentService interface {
	BatchInsertComments(ctx context.Context, comments []Comment) error
	GetRootCommentsForVideo(ctx context.Context, videoId uuid.UUID, limit int, offset int) ([]Comment, error)
	GetCommentReplies(ctx context.Context, commentId string) ([]Comment, error)
}

type Service struct {
	Store db.Store
}

func NewService(store db.Store) CommentService {
	return &Service{
		Store: store,
	}
}

type Comment struct {
	ID               string    `json:"id"`
	Text             string    `json:"test"`
	Timestamp        time.Time `json:"timestamp"`
	LikeCount        int       `json:"like_count"`
	IsFavorited      bool      `json:"is_favorited"`
	Author           string    `json:"author"`
	AuthorID         string    `json:"author_id"`
	AuthorThumbnail  string    `json:"author_thumbnail"`
	AuthorIsUploader bool      `json:"author_is_uploader"`
	Parent           *string   `json:"parent"`
	VideoID          uuid.UUID `json:"video_id"`
	Replies          int64     `json:"replies"`
}

func (s *Service) BatchInsertComments(ctx context.Context, comments []Comment) error {
	var insertCommentParams []db.BatchInsertCommentsParams

	for _, comment := range comments {
		insertCommentParams = append(insertCommentParams, commentToDBBatchComment(comment))
	}

	_, err := s.Store.BatchInsertComments(ctx, insertCommentParams)
	if err != nil {
		return fmt.Errorf("error inserting comments: %v", err)
	}

	return nil

}

func (s *Service) GetRootCommentsForVideo(ctx context.Context, videoId uuid.UUID, limit int, offset int) ([]Comment, error) {

	dbComments, err := s.Store.GetRootCommentsForVideoId(ctx, db.GetRootCommentsForVideoIdParams{
		VideoID: pgtype.UUID{Bytes: videoId, Valid: true},
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}

	if len(dbComments) == 0 {
		return nil, nil
	}

	var comments []Comment
	for _, dbComment := range dbComments {
		newComment := Comment{
			ID:               dbComment.ID,
			Text:             dbComment.Text,
			Timestamp:        dbComment.Timestamp.Time,
			LikeCount:        int(*dbComment.LikeCount),
			IsFavorited:      *dbComment.IsFavorited,
			Author:           dbComment.Author,
			AuthorID:         dbComment.ID,
			AuthorThumbnail:  *dbComment.AuthorThumbnail,
			AuthorIsUploader: *dbComment.AuthorIsUploader,
			Parent:           dbComment.Parent,
			VideoID:          dbComment.VideoID.Bytes,
			Replies:          dbComment.ReplyCount,
		}

		comments = append(comments, newComment)
	}

	return comments, nil
}

func (s *Service) GetCommentReplies(ctx context.Context, commentId string) ([]Comment, error) {

	dbComments, err := s.Store.GetCommentRepliesByCommentId(ctx, &commentId)
	if err != nil {
		return nil, err
	}

	if len(dbComments) == 0 {
		return nil, nil
	}

	var comments []Comment
	for _, dbComment := range dbComments {
		comments = append(comments, dbCommentToComment(dbComment))
	}

	return comments, nil
}

// func (s *Service) GetCommentsForVideoSorted(ctx context.Context, videoId uuid.UUID, limit int, offset int) (*Pagination, error) {

// 	// get all comments
// 	dbComments, err := s.Store.GetAllCommentsByVideoId(ctx, pgtype.UUID{Bytes: videoId, Valid: true})
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(dbComments) == 0 {
// 		return nil, nil
// 	}

// 	var comments []Comment
// 	for _, dbComment := range dbComments {
// 		comments = append(comments, dbCommentToComment(dbComment))
// 	}

// 	// add child comments to parent comments
// 	for i, comment := range comments {
// 		if comment.Parent != "" {
// 			for j, parentComment := range comments {
// 				if parentComment.ID == comment.Parent {
// 					comments[j].Replies = append(comments[j].Replies, comments[i])
// 				}
// 			}
// 		}
// 	}

// 	if offset+limit > len(comments) {
// 		limit = len(comments) - offset
// 	}

// 	// pagination
// 	pagination := Pagination{}
// 	pagination.Limit = limit
// 	pagination.Offset = offset
// 	pagination.TotalCount = len(comments)
// 	pagination.Pages = int(math.Ceil(float64(len(comments)) / float64(limit)))
// 	pagination.Data = comments[offset : offset+limit]

// 	return &pagination, nil
// }

func dbCommentToComment(dbComment db.Comment) Comment {
	return Comment{
		ID:               dbComment.ID,
		Text:             dbComment.Text,
		Timestamp:        dbComment.Timestamp.Time,
		LikeCount:        int(*dbComment.LikeCount),
		IsFavorited:      *dbComment.IsFavorited,
		Author:           dbComment.Author,
		AuthorID:         dbComment.ID,
		AuthorThumbnail:  *dbComment.AuthorThumbnail,
		AuthorIsUploader: *dbComment.AuthorIsUploader,
		Parent:           dbComment.Parent,
		VideoID:          dbComment.VideoID.Bytes,
	}
}

func commentToDBBatchComment(comment Comment) db.BatchInsertCommentsParams {
	lCount := int32(comment.LikeCount)
	return db.BatchInsertCommentsParams{
		ID:               comment.ID,
		Text:             comment.Text,
		Timestamp:        pgtype.Timestamptz{Time: comment.Timestamp, Valid: true},
		LikeCount:        &lCount,
		IsFavorited:      &comment.IsFavorited,
		Author:           comment.Author,
		AuthorID:         comment.AuthorID,
		AuthorThumbnail:  &comment.AuthorThumbnail,
		AuthorIsUploader: &comment.AuthorIsUploader,
		Parent:           comment.Parent,
		VideoID:          pgtype.UUID{Bytes: comment.VideoID, Valid: true},
	}
}
