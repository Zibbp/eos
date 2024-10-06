package blocked_paths

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/zibbp/eos/internal/db/sqlc"
)

type BlockedPath struct {
	ID         uuid.UUID `json:"id"`
	Path       string    `json:"path"`
	ErrorCount int       `json:"error_count"`
	IsBlocked  bool      `json:"is_blocked"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type BlockedPathsService interface {
	GetBlockedPaths(ctx context.Context) ([]BlockedPath, error)
	CreateOrIncrementBlockedPath(ctx context.Context, path string) error
}

type Service struct {
	Store db.Store
}

func NewService(store db.Store) BlockedPathsService {
	return &Service{
		Store: store,
	}
}

func (s *Service) GetBlockedPaths(ctx context.Context) ([]BlockedPath, error) {

	dbBlockedPaths, err := s.Store.GetBlockedPaths(ctx)
	if err != nil {
		return nil, err
	}

	var blockedPaths []BlockedPath
	for _, dbBlockedPath := range dbBlockedPaths {
		blockedPaths = append(blockedPaths, convertToBlockedPath(dbBlockedPath))
	}

	return blockedPaths, nil
}

func (s *Service) CreateOrIncrementBlockedPath(ctx context.Context, path string) error {
	_, err := s.Store.InsertBlockedPath(ctx, db.InsertBlockedPathParams{
		ID:   pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Path: path,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// duplicate key; increment error count
			err = s.Store.IncrementBlockedPathErrorCount(ctx, path)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func convertToBlockedPath(dbBlockedPath db.BlockedPath) BlockedPath {

	return BlockedPath{
		ID:         dbBlockedPath.ID.Bytes,
		Path:       dbBlockedPath.Path,
		ErrorCount: int(dbBlockedPath.ErrorCount),
		IsBlocked:  dbBlockedPath.IsBlocked,
		CreatedAt:  dbBlockedPath.CreatedAt.Time,
		UpdatedAt:  dbBlockedPath.UpdatedAt.Time,
	}
}
