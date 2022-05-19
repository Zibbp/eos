package database

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
	"github.com/zibbp/eos/internal/channel"
)

type ChannelRow struct {
	ID               sql.NullString `json:"id"`
	Name             sql.NullString `json:"name"`
	ChannelImagePath sql.NullString `json:"channel_image_path"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

func convertChannelRowToChannel(c ChannelRow) channel.Channel {
	return channel.Channel{
		ID:               c.ID.String,
		Name:             c.Name.String,
		ChannelImagePath: c.ChannelImagePath.String,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}
}

func (d *Database) CreateChannel(c echo.Context, cha channel.Channel) (channel.Channel, error) {
	channelRow := ChannelRow{
		ID:               sql.NullString{String: cha.ID, Valid: true},
		Name:             sql.NullString{String: cha.Name, Valid: true},
		ChannelImagePath: sql.NullString{String: cha.ChannelImagePath, Valid: true},
	}

	rows, err := d.Client.NamedQueryContext(c.Request().Context(), "INSERT INTO channels (id, name, channel_image_path) VALUES (:id, :name, :channelimagepath)", &channelRow)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			switch err.Code {
			case "23505":
				return channel.Channel{}, echo.NewHTTPError(http.StatusConflict, "Channel already exists")
			default:
				return channel.Channel{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
			}
		}
	}
	if err := rows.Close(); err != nil {
		return channel.Channel{}, fmt.Errorf("failed to close rows: %w", err)
	}
	return cha, nil
}

func (d *Database) GetChannels(c echo.Context) ([]channel.Channel, error) {
	// get all channels
	rows, err := d.Client.QueryContext(c.Request().Context(), "SELECT id, name, channel_image_path, created_at FROM channels")
	if err != nil {
		return nil, fmt.Errorf("failed to get channels: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error("failed to close rows: %s", err)
		}
	}()

	channels := []channel.Channel{}
	for rows.Next() {
		var channelRow ChannelRow
		if err := rows.Scan(&channelRow.ID, &channelRow.Name, &channelRow.ChannelImagePath, &channelRow.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan channel: %w", err)
		}
		channels = append(channels, channel.Channel{
			ID:               channelRow.ID.String,
			Name:             channelRow.Name.String,
			ChannelImagePath: channelRow.ChannelImagePath.String,
			CreatedAt:        channelRow.CreatedAt,
		})
	}
	return channels, nil
}

func (d *Database) GetChannel(c echo.Context, channelID string) (channel.Channel, error) {
	// get channel by id
	rows, err := d.Client.QueryContext(c.Request().Context(), "SELECT id, name, channel_image_path, created_at FROM channels WHERE id=$1", channelID)
	if err != nil {
		return channel.Channel{}, fmt.Errorf("failed to get channel: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error("failed to close rows: %s", err)
		}
	}()

	if !rows.Next() {
		return channel.Channel{}, echo.NewHTTPError(http.StatusNotFound, "Channel not found")
	}
	var channelRow ChannelRow
	if err := rows.Scan(&channelRow.ID, &channelRow.Name, &channelRow.ChannelImagePath, &channelRow.CreatedAt); err != nil {
		return channel.Channel{}, fmt.Errorf("failed to scan channel: %w", err)
	}
	return channel.Channel{
		ID:               channelRow.ID.String,
		Name:             channelRow.Name.String,
		ChannelImagePath: channelRow.ChannelImagePath.String,
		CreatedAt:        channelRow.CreatedAt,
	}, nil
}

func (d *Database) GetChannelByName(c echo.Context, channelName string) (channel.Channel, error) {
	// get channel by id
	rows, err := d.Client.QueryContext(c.Request().Context(), "SELECT id, name, channel_image_path, created_at FROM channels WHERE name=$1", channelName)
	if err != nil {
		return channel.Channel{}, fmt.Errorf("failed to get channel: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error("failed to close rows: %s", err)
		}
	}()

	if !rows.Next() {
		return channel.Channel{}, echo.NewHTTPError(http.StatusNotFound, "Channel not found")
	}
	var channelRow ChannelRow
	if err := rows.Scan(&channelRow.ID, &channelRow.Name, &channelRow.ChannelImagePath, &channelRow.CreatedAt); err != nil {
		return channel.Channel{}, fmt.Errorf("failed to scan channel: %w", err)
	}
	return channel.Channel{
		ID:               channelRow.ID.String,
		Name:             channelRow.Name.String,
		ChannelImagePath: channelRow.ChannelImagePath.String,
		CreatedAt:        channelRow.CreatedAt,
	}, nil
}

func (d *Database) ScannerGetChannel(channelID string) error {
	// get channel by id
	rows, err := d.Client.Query("SELECT id, name, channel_image_path, created_at FROM channels WHERE id=$1", channelID)
	if err != nil {
		return fmt.Errorf("failed to get channel: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error("failed to close rows: %s", err)
		}
	}()

	if !rows.Next() {
		return echo.NewHTTPError(http.StatusNotFound, "Channel not found")
	}
	var channelRow ChannelRow
	if err := rows.Scan(&channelRow.ID, &channelRow.Name, &channelRow.ChannelImagePath, &channelRow.CreatedAt); err != nil {
		return fmt.Errorf("failed to scan channel: %w", err)
	}
	return nil
}

func (d *Database) ScannerCreateChannel(cha channel.Channel) (channel.Channel, error) {
	channelRow := ChannelRow{
		ID:               sql.NullString{String: cha.ID, Valid: true},
		Name:             sql.NullString{String: cha.Name, Valid: true},
		ChannelImagePath: sql.NullString{String: cha.ChannelImagePath, Valid: true},
	}

	rows, err := d.Client.NamedQuery("INSERT INTO channels (id, name, channel_image_path) VALUES (:id, :name, :channelimagepath)", &channelRow)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			switch err.Code {
			case "23505":
				return channel.Channel{}, echo.NewHTTPError(http.StatusConflict, "Channel already exists")
			default:
				return channel.Channel{}, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
			}
		}
	}
	if err := rows.Close(); err != nil {
		return channel.Channel{}, fmt.Errorf("failed to close rows: %w", err)
	}
	return cha, nil
}
