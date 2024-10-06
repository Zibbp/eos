package scanner

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/zibbp/avalon/internal/channel"
	db "github.com/zibbp/avalon/internal/db/sqlc"
	"github.com/zibbp/avalon/internal/jobs"
	"github.com/zibbp/avalon/internal/video"
)

type ScannerService struct {
	RiverClient    *river.Client[pgx.Tx]
	Store          db.Store
	ChannelService channel.ChannelService
	VideoService   video.VideoService
	VideosDir      string
}

func NewScannerService(riverClient *river.Client[pgx.Tx], store db.Store, channelService channel.ChannelService, videoService video.VideoService, videosDir string) *ScannerService {
	scannerService := &ScannerService{
		RiverClient:    riverClient,
		Store:          store,
		ChannelService: channelService,
		VideoService:   videoService,
		VideosDir:      videosDir,
	}

	return scannerService
}

func (s *ScannerService) Scan() error {

	// use new ctx as the request will be canceled if the ctx is canceled
	ctx := context.Background()
	_, err := s.RiverClient.Insert(ctx, &jobs.VideoScanArgs{VideosDir: s.VideosDir}, nil)
	if err != nil {
		return err
	}

	return nil
}
