package parser

import (
	"context"

	db "github.com/zibbp/avalon/internal/db/sqlc"
)

type Parser interface {
	ImportVideo(ctx context.Context, jsonPath string) (*db.Video, error)
	ImportComments(ctx context.Context, jsonPath string) error
}
