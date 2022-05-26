package scanner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	channelService "github.com/zibbp/eos/internal/channel"
	"github.com/zibbp/eos/internal/comment"
	"github.com/zibbp/eos/internal/util"
	"github.com/zibbp/eos/internal/video"
)

type Store interface {
	GetVideoIDs() ([]string, error)
	ScannerGetChannel(id string) error
	ScannerCreateChannel(channel channelService.Channel) (channelService.Channel, error)
	ScannerCreateVideo(video video.Video) (video.Video, error)
	ScannerCreateComments(comments []comment.Comment) error
	ScannerCreateChapters(chapters []video.Chapter) error
	ScannerGetChannelVideoCount(id string) (int, error)
}

type Service struct {
	Store Store
}

func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

type VideoInfo struct {
	ID             string            `json:"id"`
	Channel        string            `json:"channel"`
	ChannelID      string            `json:"channel_id"`
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	Uploader       string            `json:"uploader"`
	Duration       uint64            `json:"duration"`
	ViewCount      uint64            `json:"view_count"`
	LikeCount      int64             `json:"like_count"`
	DislikeCount   int64             `json:"dislike_count"`
	Format         string            `json:"format"`
	Width          int64             `json:"width"`
	Height         int64             `json:"height"`
	Resolution     string            `json:"resolution"`
	FPS            float64           `json:"fps"`
	VideoCodec     string            `json:"vcodec"`
	VBR            float64           `json:"vbr"`
	AudioCodec     string            `json:"acodec"`
	ABR            float64           `json:"abr"`
	Epoch          int64             `json:"epoch"`
	CommentCount   uint64            `json:"comment_count"`
	VideoPath      string            `json:"video_path"`
	ThumbnailPath  string            `json:"thumbnail_path"`
	JsonPath       string            `json:"json_path"`
	SubtitlePath   string            `json:"subtitle_path"`
	UploadDate     string            `json:"upload_date"`
	TempUploadDate time.Time         `json:"temp_upload_date"`
	Comments       []comment.Comment `json:"comments"`
	Chapters       []video.Chapter   `json:"chapters"`
}

func convertScannerVideoToVideo(v VideoInfo) video.Video {
	return video.Video{
		ID:            v.ID,
		Channel:       v.Channel,
		ChannelID:     v.ChannelID,
		Title:         v.Title,
		Description:   v.Description,
		Uploader:      v.Uploader,
		Duration:      v.Duration,
		ViewCount:     v.ViewCount,
		LikeCount:     v.LikeCount,
		DislikeCount:  v.DislikeCount,
		Format:        v.Format,
		Width:         v.Width,
		Height:        v.Height,
		Resolution:    v.Resolution,
		FPS:           v.FPS,
		VideoCodec:    v.VideoCodec,
		VBR:           v.VBR,
		AudioCodec:    v.AudioCodec,
		ABR:           v.ABR,
		Epoch:         v.Epoch,
		CommentCount:  v.CommentCount,
		VideoPath:     v.VideoPath,
		ThumbnailPath: v.ThumbnailPath,
		JsonPath:      v.JsonPath,
		SubtitlePath:  v.SubtitlePath,
		UploadDate:    v.TempUploadDate,
	}
}

func convertScannerCommentToComment(c comment.Comment, vidID string) comment.Comment {
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
		VideoID:          vidID,
	}
}

func convertScannerChapterToChapter(c video.Chapter, vidID string) video.Chapter {
	return video.Chapter{
		ID:        c.ID,
		StartTime: c.StartTime,
		Title:     c.Title,
		EndTime:   c.EndTime,
		VideoID:   vidID,
	}
}

func (s *Service) StartScanner(c echo.Context) error {
	log.Info("Starting scanner")

	// Get root channel folders to process
	channelDirs, err := util.GetDirFolders("/videos")
	if err != nil {
		log.Errorf("Failed to get channel folders: %s", err)
		return c.JSON(http.StatusInternalServerError, "Failed to get channel folders")
	}

	// Fetch already imported videos
	log.Debug("Fetching imported video IDs")
	importedVideoIds, err := s.Store.GetVideoIDs()
	if err != nil {
		log.Errorf("Failed to get imported video IDs: %s", err)
		return c.JSON(http.StatusInternalServerError, "Failed to get imported video IDs")
	}

	// Create goroutines for each channel folder
	channel := make(chan string)

	for _, channelDir := range channelDirs {
		go s.processChannelDir(channelDir, channel, importedVideoIds)
	}

	return nil
}

// Process each channel directory and it's videos
func (s *Service) processChannelDir(channelDir string, channel chan string, importedVideoIds []string) {
	log.Debugf("Processing channel directory: %s", channelDir)
	channelDirPath := fmt.Sprintf("/videos/%s", channelDir)
	// Get each video folder within the channel directory
	videoDirs, err := util.GetDirFolders(channelDirPath)
	if err != nil {
		log.Errorf("Failed to get video folders: %s", err)
		return
	}
	// Note number of videoDirs to process
	numVideoDirs := len(videoDirs)
	// YT-DLP script default folder (need to skip this)
	ytDlpChannelFolder := fmt.Sprintf("%v-NA-%v-Videos", channelDir, channelDir)
	// Skip YT-Archive generate playlist folder
	ytdlpPlaylistFolder := fmt.Sprintf("%v-NA-YT-Archive", channelDir)
	// process each video folder
	for i, videoDir := range videoDirs {
		// Skip the default yt-dlp folder
		if videoDir == ytDlpChannelFolder || videoDir == ytdlpPlaylistFolder {
			continue
		}
		// Create video item
		var vid VideoInfo
		// Get the files in the video directory
		videoFiles, err := util.GetFilesinDir(fmt.Sprintf("/videos/%s/%s", channelDir, videoDir))
		if err != nil {
			log.Errorf("Failed to get video files: %s", err)
			continue
		}
		// Process each video file
		for _, videoFile := range videoFiles {
			// Get file extension
			sliceVideoFileString := strings.LastIndex(videoFile, ".")
			ext := videoFile[sliceVideoFileString+1:]
			// Switch for each file
			switch ext {
			case "json":
				// Ensure we have the correct file
				if strings.Contains(videoFile, "info.json") {
					// Get JSON data
					readVideoJson(fmt.Sprintf("/videos/%s/%s/%s", channelDir, videoDir, videoFile), &vid)
					// Format date
					vid.TempUploadDate = util.StringDateToDateTime(vid.UploadDate)
					vid.JsonPath = fmt.Sprintf("/%s/%s/%s", channelDir, videoDir, videoFile)
				}
			case "mkv":
				vid.VideoPath = fmt.Sprintf("/%s/%s/%s", channelDir, videoDir, videoFile)
			case "webp":
				vid.ThumbnailPath = fmt.Sprintf("/%s/%s/%s", channelDir, videoDir, videoFile)
			case "jpg":
				vid.ThumbnailPath = fmt.Sprintf("/%s/%s/%s", channelDir, videoDir, videoFile)
			case "vtt":
				vid.SubtitlePath = fmt.Sprintf("/%s/%s/%s", channelDir, videoDir, videoFile)
			case "srt":
				vid.SubtitlePath = fmt.Sprintf("/%s/%s/%s", channelDir, videoDir, videoFile)
			default:
			}
		}

		// First video in the channel directory - do some processing
		// Check if channel exists
		if i == 0 {

			err = s.Store.ScannerGetChannel(vid.ChannelID)
			if err != nil {
				log.Errorf("Failed to get channel: %s", err)
				// This failing means the channel doesn't exist, so we need to create it
				newChannel := channelService.Channel{
					ID:               vid.ChannelID,
					Name:             vid.Channel,
					ChannelImagePath: fmt.Sprintf("/%s/%s.jpg", channelDir, channelDir),
				}
				cha, err := s.Store.ScannerCreateChannel(newChannel)
				if err != nil {
					log.Errorf("Failed to create channel: %s", err)
					return
				}
				log.Infof("Scanner created channel %s", cha.Name)
			}

			// Compare number of videoDirs to number of video db entries
			vidDBRows, err := s.Store.ScannerGetChannelVideoCount(vid.ChannelID)
			if err != nil {
				log.Errorf("Failed to get channel video count: %s", err)
				return
			}
			vidDBRows = vidDBRows + 1
			if vidDBRows == numVideoDirs || vidDBRows > numVideoDirs {
				log.Infof("All videos in channel directory %s have already been processed - skipping.", channelDir)
				return
			}
		}

		// Check if video is already imported
		if util.StringInSlice(vid.ID, importedVideoIds) {
			log.Debugf("Video already imported: %s", vid.ID)
			continue
		}

		// Insert video
		convertedVid := convertScannerVideoToVideo(vid)
		insertedVid, err := s.Store.ScannerCreateVideo(convertedVid)
		if err != nil {
			log.Errorf("Failed to create video: %s", err)
			continue
		}
		log.Infof("Video %s imported", insertedVid.ID)
		// Insert chapters
		if len(vid.Chapters) > 0 {
			var convertedChapters []video.Chapter
			for _, chap := range vid.Chapters {
				tmpChapter := convertScannerChapterToChapter(chap, vid.ID)
				convertedChapters = append(convertedChapters, tmpChapter)
			}

			err := s.Store.ScannerCreateChapters(convertedChapters)
			if err != nil {
				log.Errorf("Failed to create chapters: %s", err)
				continue
			}
			log.Infof("Chapters imported for video %s", insertedVid.ID)
		}
		// Inset Comments
		if len(vid.Comments) > 0 {
			var convertedCmts []comment.Comment

			for _, cmt := range vid.Comments {
				tmpCmt := convertScannerCommentToComment(cmt, vid.ID)
				convertedCmts = append(convertedCmts, tmpCmt)
			}
			err := s.Store.ScannerCreateComments(convertedCmts)
			if err != nil {
				log.Errorf("Failed to create comments: %s", err)
				continue
			}
			log.Infof("Comments imported for video %s", vid.ID)
		}
	}
	log.Infof("Channel %s processed", channelDir)
	channel <- "Channel Processed"
}

func readVideoJson(videoJsonPath string, vid *VideoInfo) {
	content, err := ioutil.ReadFile(videoJsonPath)
	if err != nil {
		fmt.Println("Error reading file", err)
	}

	err = json.Unmarshal(content, &vid)
	if err != nil {
		fmt.Println("Error unmarshalling json", err)
	}
}
