package channel

import (
	"time"

	"github.com/labstack/echo/v4"
)

type Store interface {
	CreateChannel(c echo.Context, cha Channel) (Channel, error)
	GetChannels(c echo.Context) ([]Channel, error)
	GetChannel(c echo.Context, id string) (Channel, error)
	GetChannelByName(c echo.Context, channelName string) (Channel, error)
}

type Service struct {
	Store Store
}

func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

type Channel struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	ChannelImagePath string    `gorm:"not null" json:"channel_image_path"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// func (s *Service) CreateChannel(c echo.Context, cha *db.Channel) (*db.Channel, error) {
// 	cha, err := s.Store.CreateChannel(c, cha)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return cha, nil
// }

func (s *Service) GetChannels(c echo.Context) ([]Channel, error) {
	channels, err := s.Store.GetChannels(c)
	if err != nil {
		return nil, err
	}
	return channels, nil
}

func (s *Service) CreateChannel(c echo.Context, cha Channel) (Channel, error) {
	cha, err := s.Store.CreateChannel(c, cha)
	if err != nil {
		return cha, err
	}
	return cha, nil
}

func (s *Service) GetChannel(c echo.Context, id string) (Channel, error) {
	cha, err := s.Store.GetChannel(c, id)
	if err != nil {
		return cha, err
	}
	return cha, nil
}

func (s *Service) GetChannelByName(c echo.Context, channelName string) (Channel, error) {
	cha, err := s.Store.GetChannelByName(c, channelName)
	if err != nil {
		return cha, err
	}
	return cha, nil
}
