package jobs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
	"github.com/zibbp/eos/internal/blocked_paths"
	"github.com/zibbp/eos/internal/utils"
	// "github.com/zibbp/eos/internal/yt"
)

type VideoScanArgs struct {
	VideosDir string `json:"videos_dir"`
}

func (VideoScanArgs) Kind() string { return QueueScanVideo }

func (VideoScanArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		MaxAttempts: 2,
		Queue:       QueueScanVideo,
		Priority:    1,
	}
}

func (VideoScanArgs) Timeout(job *river.Job[VideoScanArgs]) time.Duration {
	return -1
}

type VideoScanWorker struct {
	river.WorkerDefaults[VideoScanArgs]
}

func (w *VideoScanWorker) Work(ctx context.Context, job *river.Job[VideoScanArgs]) error {
	path := job.Args.VideosDir

	log := log.With().Str("job_id", fmt.Sprintf("%d", job.ID)).Str("attempt", fmt.Sprintf("%d", job.Attempt)).Str("attempted_by", job.AttemptedBy[job.Attempt-1]).Logger()

	log.Info().Str("path", path).Msg("scanning for videos")

	// get store from context
	store, err := StoreFromContext(ctx)
	if err != nil {
		return err
	}

	// get blocked paths
	blockedPathsService := blocked_paths.NewService(store)
	blockedPaths, err := blockedPathsService.GetBlockedPaths(context.Background())
	if err != nil {
		return err
	}

	// check if path exists
	_, err = os.ReadDir(path)
	if err != nil {
		return err
	}

	// get paths to existing videos
	videoInfoPaths, err := store.GetVideoInfoPaths(context.Background())
	if err != nil {
		return err
	}

	// walk the path looking for json files
	err = filepath.Walk(path, func(subPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		// check if file is json
		if filepath.Ext(info.Name()) != ".json" {
			return nil
		}

		// check if path is already processed
		if utils.StringInSlice(videoInfoPaths, subPath) {
			return nil
		}

		// check if path is blocked
		for _, blockedPath := range blockedPaths {
			cleanedSubPath := filepath.Clean(subPath)
			cleanedBlockedPath := filepath.Clean(blockedPath.Path)

			if cleanedSubPath == cleanedBlockedPath && blockedPath.ErrorCount >= 5 {
				log.Debug().Str("path", subPath).Msg("skipping blocked path")
				return nil
			}
		}

		// insert queue item to process video
		client := river.ClientFromContext[pgx.Tx](ctx)
		_, jobErr := client.Insert(ctx, VideoImportArgs{
			VideoInfoPath: subPath,
		}, nil)
		if jobErr != nil {
			return fmt.Errorf("error inserting queue item: %v", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// ///////////////////
// Import Video Job //
// //////////////////
type VideoImportArgs struct {
	VideoInfoPath string `json:"video_info_path"`
}

func (VideoImportArgs) Kind() string { return QueueImportVideo }

func (VideoImportArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		MaxAttempts: 2,
		Queue:       QueueImportVideo,
		Priority:    1,
	}
}

func (VideoImportArgs) Timeout(job *river.Job[VideoImportArgs]) time.Duration {
	return 30 * time.Second
}

type VideoImportWorker struct {
	river.WorkerDefaults[VideoImportArgs]
}

func (w *VideoImportWorker) Work(ctx context.Context, job *river.Job[VideoImportArgs]) error {

	log := log.With().Str("job_id", fmt.Sprintf("%d", job.ID)).Str("attempt", fmt.Sprintf("%d", job.Attempt)).Str("attempted_by", job.AttemptedBy[job.Attempt-1]).Logger()

	log.Info().Msg("processing video")

	// get store from context
	store, err := StoreFromContext(ctx)
	if err != nil {
		return err
	}

	// parser from context
	parser, err := ParserFromContext(ctx)
	if err != nil {
		return err
	}

	blockedPathsService := blocked_paths.NewService(store)

	// import video
	video, err := parser.ImportVideo(ctx, job.Args.VideoInfoPath)
	if err != nil {

		blockedPathErr := blockedPathsService.CreateOrIncrementBlockedPath(ctx, job.Args.VideoInfoPath)
		if blockedPathErr != nil {
			return err
		}

		return fmt.Errorf("error importing video: %s", err)
	}

	// insert task to import comments
	client := river.ClientFromContext[pgx.Tx](ctx)
	_, err = client.Insert(ctx, &VideoImportCommentsArgs{
		VideoInfoPath: video.InfoPath,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

// /////////////////////////////
// Import Video Comments Job //
// /////////////////////////////
type VideoImportCommentsArgs struct {
	VideoInfoPath string `json:"video_info_path"`
}

func (VideoImportCommentsArgs) Kind() string { return QueueImportVideoComments }

func (VideoImportCommentsArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		MaxAttempts: 2,
		Queue:       QueueImportVideo,
		Priority:    2,
	}
}

func (VideoImportCommentsArgs) Timeout(job *river.Job[VideoImportCommentsArgs]) time.Duration {
	return 5 * time.Minute
}

type VideoImportCommentsWorker struct {
	river.WorkerDefaults[VideoImportCommentsArgs]
}

func (w *VideoImportCommentsWorker) Work(ctx context.Context, job *river.Job[VideoImportCommentsArgs]) error {

	// parser from context
	parser, err := ParserFromContext(ctx)
	if err != nil {
		return err
	}

	// import comments
	if err := parser.ImportComments(ctx, job.Args.VideoInfoPath); err != nil {
		return err
	}

	return nil
}
