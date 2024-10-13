package parser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/zibbp/eos/internal/channel"
	"github.com/zibbp/eos/internal/chapter"
	"github.com/zibbp/eos/internal/comment"
	db "github.com/zibbp/eos/internal/db/sqlc"
	"github.com/zibbp/eos/internal/utils"
	"github.com/zibbp/eos/internal/video"
)

type YoutubeParser struct {
	Store          db.Store
	ChannelService channel.ChannelService
	VideoService   video.VideoService
	CommentService comment.CommentService
	ChapterService chapter.ChapterService
}

func (p *YoutubeParser) ImportVideo(ctx context.Context, jsonPath string) (*db.Video, error) {
	path := filepath.Dir(jsonPath)

	// parse youtube video info
	info, err := getYoutubeVideoInfo(jsonPath)
	if err != nil {
		return nil, err
	}

	// sanity check some required fields
	if info.ID == "" || info.Channel == "" || info.Title == "" || info.UploadDate == "" {
		return nil, fmt.Errorf("invalid youtube video info - missing id, channel, title or upload date")
	}

	videoDirPath := filepath.Dir(jsonPath)

	videoPath := ""
	thumbnailPath := ""

	// get all files in the directory
	files, err := utils.GetFilesInDirectory(videoDirPath)
	if err != nil {
		return nil, err
	}

	var subtitlePaths []string

	// find required files
	for _, file := range files {
		// find video file
		if filepath.Ext(file.Name()) == ".mkv" || filepath.Ext(file.Name()) == ".mp4" || filepath.Ext(file.Name()) == ".webm" {
			videoPath = filepath.Join(videoDirPath, file.Name())
		}
		// find thumbnail file
		// TODO: remove hardcoded thumbnails.jpg. Delete in all directories and move to sub-dir??
		if (filepath.Ext(file.Name()) == ".webp" || filepath.Ext(file.Name()) == ".png" || filepath.Ext(file.Name()) == ".jpg") && file.Name() != "thumbnails.jpg" {
			thumbnailPath = filepath.Join(videoDirPath, file.Name())
		}
		if filepath.Ext(file.Name()) == ".vtt" {
			subtitlePaths = append(subtitlePaths, filepath.Join(videoDirPath, file.Name()))
		}
	}

	if videoPath == "" || thumbnailPath == "" {
		return nil, fmt.Errorf("missing video or thumbnail path: %s", path)
	}

	// get or create channel
	c, err := p.Store.GetChannelByExternalId(ctx, info.ChannelID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// create channel
			channelPath := upDirectories(jsonPath, 2)
			newChannel, err := p.ChannelService.CreateChannel(ctx, channel.CreateChannelInput{
				ExternalID:         info.ChannelID,
				Name:               info.Channel,
				Path:               channelPath,
				GenerateThumbnails: false,
			})
			if err != nil {
				return nil, err
			}
			c = *newChannel
		} else {
			return nil, err
		}
	}

	// parse upload date
	parsedUploadDate, err := time.Parse("20060102", info.UploadDate)
	if err != nil {
		return nil, err
	}

	// fmt.Println(parsedUploadDate)
	// fmt.Println(c.ID.String())

	videoInsert := db.InsertVideoParams{
		ID:            pgtype.UUID{Bytes: uuid.New(), Valid: true},
		ExtID:         &info.ID,
		Title:         info.Title,
		Description:   &info.Description,
		UploadDate:    pgtype.Timestamptz{Time: parsedUploadDate, Valid: true},
		Uploader:      &info.Uploader,
		Duration:      int32(info.Duration),
		ViewCount:     info.ViewCount,
		LikeCount:     &info.LikeCount,
		DislikeCount:  &info.DislikeCount,
		Format:        &info.Format,
		Height:        &info.Height,
		Width:         &info.Width,
		Resolution:    &info.Resolution,
		Fps:           &info.FPS,
		VideoCodec:    &info.VideoCodec,
		Vbr:           &info.VBR,
		AudioCodec:    &info.AudioCodec,
		Abr:           &info.ABR,
		CommentCount:  &info.CommentCount,
		VideoPath:     videoPath,
		ThumbnailPath: thumbnailPath,
		InfoPath:      jsonPath,
		SubtitlePath:  subtitlePaths,
		Path:          path,
		// StoryboardPath: pgtype.Text{String: "", Valid: true},
		ChannelID: c.ID,
	}

	video, err := p.Store.InsertVideo(ctx, videoInsert)
	if err != nil {
		return nil, err
	}

	// import chapters
	if len(info.Chapters) > 0 {
		var chapters []chapter.Chapter

		for _, c := range info.Chapters {

			chapters = append(chapters, chapter.Chapter{
				ID:        uuid.New(),
				Title:     c.Title,
				StartTime: int(c.StartTime),
				EndTime:   int(c.EndTime),
				VideoID:   video.ID.Bytes,
			})
		}

		if err := p.ChapterService.BatchInsertChapters(ctx, chapters); err != nil {
			return nil, err
		}

		log.Info().Str("video_id", info.ID).Msgf("imported %d chapters", len(chapters))
	}

	log.Info().Str("video_id", info.ID).Msg("imported video")

	return &video, nil
}

func (p *YoutubeParser) ImportComments(ctx context.Context, jsonPath string) error {
	info, err := getYoutubeVideoInfo(jsonPath)
	if err != nil {
		return err
	}

	dbVideo, err := p.VideoService.GetVideoByExtId(ctx, info.ID)
	if err != nil {
		return err
	}

	if len(info.Comments) == 0 {
		log.Info().Str("video_id", info.ID).Msg("imported video")
		return nil
	}

	var comments []comment.Comment

	for _, c := range info.Comments {
		timestamp := time.Unix(c.Timestamp, 0)
		// if parent is root then it should be nil in the database
		var parent *string
		if c.Parent != "root" {
			parentCopy := c.Parent
			parent = &parentCopy
		}
		comments = append(comments, comment.Comment{
			ID:               c.ID,
			Text:             c.Text,
			Timestamp:        timestamp,
			LikeCount:        int(c.LikeCount),
			IsFavorited:      c.IsFavorited,
			Author:           c.Author,
			AuthorID:         c.AuthorID,
			AuthorThumbnail:  c.AuthorThumbnail,
			AuthorIsUploader: c.AuthorIsUploader,
			Parent:           parent,
			VideoID:          dbVideo.ID,
		})
	}

	if err := p.CommentService.BatchInsertComments(ctx, comments); err != nil {
		return err
	}

	log.Info().Str("video_id", info.ID).Msgf("imported %d comments", len(comments))

	return nil
}

// struct for the contents of the info.json file saved by yt-dlp
type YoutubeVideoInfo struct {
	ID             string    `json:"id"`
	Channel        string    `json:"channel"`
	ChannelID      string    `json:"channel_id"`
	Title          string    `json:"title"`
	Formats        []Format  `json:"formats"`
	Description    string    `json:"description"`
	Uploader       string    `json:"uploader"`
	Duration       int32     `json:"duration"`
	ViewCount      int64     `json:"view_count"`
	LikeCount      int64     `json:"like_count"`
	DislikeCount   int64     `json:"dislike_count"`
	Format         string    `json:"format"`
	Width          int32     `json:"width"`
	Height         int32     `json:"height"`
	Resolution     string    `json:"resolution"`
	FPS            float32   `json:"fps"`
	VideoCodec     string    `json:"vcodec"`
	VBR            float32   `json:"vbr"`
	AudioCodec     string    `json:"acodec"`
	ABR            float32   `json:"abr"`
	Epoch          int32     `json:"epoch"`
	CommentCount   int32     `json:"comment_count"`
	VideoPath      string    `json:"video_path"`
	ThumbnailPath  string    `json:"thumbnail_path"`
	JsonPath       string    `json:"json_path"`
	SubtitlePath   string    `json:"subtitle_path"`
	UploadDate     string    `json:"upload_date"`
	TempUploadDate time.Time `json:"temp_upload_date"`
	Path           string    `json:"path"`
	Comments       []Comment `json:"comments"`
	Type           string    `json:"_type"`
	Chapters       []Chapter `json:"chapters"`
}

type Comment struct {
	ID               string    `json:"id"`
	Text             string    `json:"text"`
	Timestamp        int64     `json:"timestamp"`
	LikeCount        int64     `json:"like_count"`
	IsFavorited      bool      `json:"is_favorited"`
	Author           string    `json:"author"`
	AuthorID         string    `json:"author_id"`
	AuthorThumbnail  string    `json:"author_thumbnail"`
	AuthorIsUploader bool      `json:"author_is_uploader"`
	Parent           string    `json:"parent"`
	VideoID          string    `json:"video_id"`
	Replies          []Comment `json:"replies"`
}

type Format struct {
	FormatID   string     `json:"format_id"`
	FormatNote string     `json:"format_note"`
	Width      *int64     `json:"width,omitempty"`
	Height     *int64     `json:"height,omitempty"`
	FPS        *float64   `json:"fps,omitempty"`
	Rows       *int64     `json:"rows,omitempty"`
	Columns    *int64     `json:"columns,omitempty"`
	Fragments  []Fragment `json:"fragments,omitempty"`
}

type Fragment struct {
	URL      string  `json:"url"`
	Duration float64 `json:"duration"`
}

type Chapter struct {
	StartTime float64 `json:"start_time"`
	EndTime   float64 `json:"end_time"`
	Title     string  `json:"title"`
}

// getVideoInfo reads a JSON file from a given path and decodes it into a VideoInfo struct.
func getYoutubeVideoInfo(path string) (YoutubeVideoInfo, error) {
	var info YoutubeVideoInfo
	file, err := os.Open(path)
	if err != nil {
		return info, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&info); err != nil {
		return info, err
	}

	return info, nil
}

func upDirectories(path string, levels int) string {
	for i := 0; i < levels; i++ {
		path = filepath.Dir(path)
	}
	return path
}
