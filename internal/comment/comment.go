package comment

import "github.com/labstack/echo/v4"

type Store interface {
	CreateComment(c echo.Context, cmt Comment) (Comment, error)
	GetVideoComments(c echo.Context, vidID string, cmtList CommentList) (CommentList, error)
}

type Service struct {
	Store Store
}

func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

type Comment struct {
	ID               string    `json:"id"`
	Text             string    `json:"text"`
	Timestamp        uint64    `json:"timestamp"`
	LikeCount        uint64    `json:"like_count"`
	IsFavorited      bool      `json:"is_favorited"`
	Author           string    `json:"author"`
	AuthorID         string    `json:"author_id"`
	AuthorThumbnail  string    `json:"author_thumbnail"`
	AuthorIsUploader bool      `json:"author_is_uploader"`
	Parent           string    `json:"parent"`
	VideoID          string    `json:"video_id"`
	Replies          []Comment `json:"replies"`
}

type CommentList struct {
	Limit      int       `json:"limit"`
	Page       int       `json:"page"`
	PrevPage   int       `json:"prev_page"`
	NextPage   int       `json:"next_page"`
	LastPage   int       `json:"last_page"`
	TotalItems int       `json:"total_items"`
	Items      []Comment `json:"items"`
}

func (s *Service) GetVideoComments(c echo.Context, vidID string, cmtList CommentList) (CommentList, error) {
	cmts, err := s.Store.GetVideoComments(c, vidID, cmtList)
	if err != nil {
		return CommentList{}, err
	}
	return cmts, nil
}

func (s *Service) CreateComment(c echo.Context, cmt Comment) (Comment, error) {
	cmt, err := s.Store.CreateComment(c, cmt)
	if err != nil {
		return cmt, err
	}
	return cmt, nil
}
