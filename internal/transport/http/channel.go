package http

import (
	"net/http"
	"net/url"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/zibbp/eos/internal/channel"
)

type ChannelService interface {
	CreateChannel(c echo.Context, cha channel.Channel) (channel.Channel, error)
	GetChannels(c echo.Context) ([]channel.Channel, error)
	GetChannel(c echo.Context, id string) (channel.Channel, error)
	GetChannelByName(c echo.Context, channelName string) (channel.Channel, error)
}

type CreateChannel struct {
	ID               string `json:"id" validate:"required"`
	Name             string `json:"name" validate:"required"`
	ChannelImagePath string `json:"channel_image_path" validate:"required"`
}

func channelPostToChannel(c CreateChannel) channel.Channel {
	return channel.Channel{
		ID:               c.ID,
		Name:             c.Name,
		ChannelImagePath: c.ChannelImagePath,
	}
}

func (h *Handler) GetChannels(c echo.Context) error {

	channels, err := h.Service.ChannelService.GetChannels(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, channels)
}

func (h *Handler) CreateChannel(c echo.Context) error {
	var createChannel CreateChannel
	if err := c.Bind(&createChannel); err != nil {
		return err
	}

	validate := validator.New()
	err := validate.Struct(createChannel)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	cha := channelPostToChannel(createChannel)
	cha, err = h.Service.ChannelService.CreateChannel(c, cha)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, cha)
}

func (h *Handler) GetChannel(c echo.Context) error {
	id := c.Param("id")
	cha, err := h.Service.ChannelService.GetChannel(c, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, cha)
}

func (h *Handler) GetChannelByName(c echo.Context) error {
	name := c.Param("name")
	// Channel name needs to be deocded as it's used in the frontend
	decodedName, err := url.QueryUnescape(name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Failed to decode channel name")
	}
	cha, err := h.Service.ChannelService.GetChannelByName(c, decodedName)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, cha)
}
