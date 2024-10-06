package chapter

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/zibbp/avalon/internal/db/sqlc"
)

type ChapterService interface {
	BatchInsertChapters(ctx context.Context, chapters []Chapter) error
	GetChaptersForVideo(ctx context.Context, videoId uuid.UUID) ([]Chapter, error)
	GetChaptersForVidstackPlayer(ctx context.Context, videoId uuid.UUID) ([]VidstackPlayerChapter, error)
}

type Service struct {
	Store db.Store
}

func NewService(store db.Store) ChapterService {
	return &Service{
		Store: store,
	}
}

type Chapter struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	StartTime int       `json:"start_time"`
	EndTime   int       `json:"end_time"`
	VideoID   uuid.UUID `json:"video_id"`
}

// expected format for Vidstack player
type VidstackPlayerChapter struct {
	Text      string `json:"text"`
	StartTime int    `json:"startTime"`
	EndTime   int    `json:"endTime"`
}

func (s *Service) BatchInsertChapters(ctx context.Context, chapters []Chapter) error {
	var insertChapterParams []db.BatchInsertChaptersParams

	for _, comment := range chapters {
		insertChapterParams = append(insertChapterParams, chapterToDBBatchChapter(comment))
	}

	_, err := s.Store.BatchInsertChapters(ctx, insertChapterParams)
	if err != nil {
		return fmt.Errorf("error inserting comments: %v", err)
	}

	return nil

}

func (s *Service) GetChaptersForVideo(ctx context.Context, videoId uuid.UUID) ([]Chapter, error) {
	dbChapters, err := s.Store.GetChaptersByVideoId(ctx, pgtype.UUID{Bytes: videoId, Valid: true})
	if err != nil {
		return nil, err
	}

	var chapters []Chapter
	for _, dbChapter := range dbChapters {
		chapters = append(chapters, dbChapterToChapter(dbChapter))
	}

	return chapters, nil
}

func (s *Service) GetChaptersForVidstackPlayer(ctx context.Context, videoId uuid.UUID) ([]VidstackPlayerChapter, error) {
	dbChapters, err := s.Store.GetChaptersByVideoId(ctx, pgtype.UUID{Bytes: videoId, Valid: true})
	if err != nil {
		return nil, err
	}

	var chapters []VidstackPlayerChapter
	for _, dbChapter := range dbChapters {
		chapters = append(chapters, VidstackPlayerChapter{
			Text:      dbChapter.Title,
			StartTime: int(dbChapter.StartTime),
			EndTime:   int(dbChapter.EndTime),
		})
	}

	return chapters, nil
}

func dbChapterToChapter(dbChapter db.Chapter) Chapter {
	return Chapter{
		ID:        dbChapter.ID.Bytes,
		Title:     dbChapter.Title,
		StartTime: int(dbChapter.StartTime),
		EndTime:   int(dbChapter.EndTime),
		VideoID:   dbChapter.VideoID.Bytes,
	}
}

func chapterToDBBatchChapter(chapter Chapter) db.BatchInsertChaptersParams {
	return db.BatchInsertChaptersParams{
		ID:        pgtype.UUID{Bytes: chapter.ID, Valid: true},
		Title:     chapter.Title,
		StartTime: int32(chapter.StartTime),
		EndTime:   int32(chapter.EndTime),
		VideoID:   pgtype.UUID{Bytes: chapter.VideoID, Valid: true},
	}
}
