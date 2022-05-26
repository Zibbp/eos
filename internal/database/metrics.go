package database

import (
	"fmt"

	"github.com/labstack/gommon/log"
	"github.com/zibbp/eos/internal/metrics"
)

func (d *Database) GetChannelVideoCount() ([]metrics.ChannelVideoCount, error) {
	// Get all channels and the number of videos they have
	rows, err := d.Client.Query("SELECT channels.name, COUNT(videos.id) AS video_count FROM channels LEFT JOIN videos ON channels.id = videos.channel_id GROUP BY channels.name")
	if err != nil {
		return nil, fmt.Errorf("failed to get channels: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error("failed to close rows: %s", err)
		}
	}()

	channelVideoCounts := []metrics.ChannelVideoCount{}
	for rows.Next() {
		var channelVideoCount metrics.ChannelVideoCount
		if err := rows.Scan(&channelVideoCount.Channel, &channelVideoCount.Count); err != nil {
			return nil, fmt.Errorf("failed to scan channel: %w", err)
		}
		channelVideoCounts = append(channelVideoCounts, metrics.ChannelVideoCount{
			Channel: channelVideoCount.Channel,
			Count:   channelVideoCount.Count,
		})
	}
	return channelVideoCounts, nil

}

func (d *Database) GetVideoCount() (int64, error) {
	// Get total videos
	var count int64
	if err := d.Client.Get(&count, "SELECT COUNT(*) FROM videos"); err != nil {
		return 0, fmt.Errorf("failed to get total videos: %w", err)
	}
	return count, nil
}

func (d *Database) GetChannelCount() (int64, error) {
	// Get total channels
	var count int64
	if err := d.Client.Get(&count, "SELECT COUNT(*) FROM channels"); err != nil {
		return 0, fmt.Errorf("failed to get total channels: %w", err)
	}
	return count, nil
}
