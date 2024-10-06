package channel

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/zibbp/eos/internal/db/sqlc"
)

type ChannelService interface {
	CreateChannel(ctx context.Context, input CreateChannelInput) (*db.Channel, error)
	GetChannels(ctx context.Context) ([]db.Channel, error)
	GetChannelByName(ctx context.Context, name string) (*db.Channel, error)
	GetChannelByID(ctx context.Context, id uuid.UUID) (*db.Channel, error)
}

type CreateChannelInput struct {
	ExternalID         string
	Name               string
	Path               string
	GenerateThumbnails bool
}

type Service struct {
	Store db.Store
}

func NewService(store db.Store) ChannelService {
	return &Service{
		Store: store,
	}
}

var ErrChannelNotFound = errors.New("channel not found")

func (s *Service) CreateChannel(ctx context.Context, input CreateChannelInput) (*db.Channel, error) {
	imagePath, err := getChannelImageInPath(input.Path)
	if err != nil {
		return nil, err
	}

	c, err := s.Store.InsertChannel(ctx, db.InsertChannelParams{
		ID:                 pgtype.UUID{Bytes: uuid.New(), Valid: true},
		ExtID:              input.ExternalID,
		Name:               input.Name,
		ImagePath:          &imagePath,
		GenerateThumbnails: input.GenerateThumbnails,
	})
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (s *Service) GetChannels(ctx context.Context) ([]db.Channel, error) {

	dbChannels, err := s.Store.GetChannels(ctx)
	if err != nil {
		return nil, err
	}

	return dbChannels, nil
}

func (s *Service) GetChannelByName(ctx context.Context, name string) (*db.Channel, error) {

	dbChannel, err := s.Store.GetChannelByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrChannelNotFound
		}
		return nil, err
	}

	return &dbChannel, nil
}

func (s *Service) GetChannelByID(ctx context.Context, id uuid.UUID) (*db.Channel, error) {
	channel, err := s.Store.GetChannelById(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}

	return &channel, nil
}

// getChannelImageInPath returns the path to the channel image if found in the channel directory
func getChannelImageInPath(path string) (string, error) {
	validExtensions := map[string]bool{
		".webp": true,
		".png":  true,
		".jpg":  true,
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		ext := filepath.Ext(name)
		if validExtensions[ext] && name != "thumbnails.jpg" {
			return filepath.Join(path, name), nil
		}
	}

	return "", nil
}
