package config

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	DB_HOST string `env:"DB_HOST,required"`
	DB_PORT string `env:"DB_PORT,required"`
	DB_USER string `env:"DB_USER,required"`
	DB_PASS string `env:"DB_PASS,required"`
	DB_NAME string `env:"DB_NAME,required"`

	VIDEOS_DIR string `env:"VIDEOS_DIR,required"`

	CDN_URL string `env:"CDN_URL,required"`

	MaxVideoImportWorkers int `env:"MAX_VIDEO_IMPORT_WORKERS,required"`

	DB_URL string
}

func GetConfig() Config {
	ctx := context.Background()

	var c Config
	if err := envconfig.Process(ctx, &c); err != nil {
		log.Panic().Err(err).Msg("error parsing config")
	}

	c.DB_URL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.DB_USER, c.DB_PASS, c.DB_HOST, c.DB_PORT, c.DB_NAME)

	return c
}
