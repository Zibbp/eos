package jobs

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	db "github.com/zibbp/avalon/internal/db/sqlc"
	"github.com/zibbp/avalon/internal/parser"
)

type contextKey string

const StoreKey contextKey = "store"
const ParserKey contextKey = "parser"

var (
	QueueImportVideo         = "import_video"
	QueueImportVideoComments = "import_video_comments"
	QueueScanVideo           = "scan_video"
)

func StoreFromContext(ctx context.Context) (db.Store, error) {
	store, exists := ctx.Value(StoreKey).(db.Store)
	if !exists || store == nil {
		return nil, errors.New("store not found in context")
	}

	return store, nil
}

func ParserFromContext(ctx context.Context) (parser.Parser, error) {
	parser, exists := ctx.Value(ParserKey).(parser.Parser)
	if !exists || parser == nil {
		return nil, errors.New("parser not found in context")
	}

	return parser, nil
}

type CustomErrorHandler struct{}

func (*CustomErrorHandler) HandleError(ctx context.Context, job *rivertype.JobRow, err error) *river.ErrorHandlerResult {
	log.Error().Str("job_id", fmt.Sprintf("%d", job.ID)).Str("attempt", fmt.Sprintf("%d", job.Attempt)).Str("attempted_by", job.AttemptedBy[job.Attempt-1]).Str("args", string(job.EncodedArgs)).Err(err).Msg("task error")

	return nil
}

func (*CustomErrorHandler) HandlePanic(ctx context.Context, job *rivertype.JobRow, panicVal any, trace string) *river.ErrorHandlerResult {
	log.Error().Str("job_id", fmt.Sprintf("%d", job.ID)).Str("attempt", fmt.Sprintf("%d", job.Attempt)).Str("attempted_by", job.AttemptedBy[job.Attempt-1]).Str("args", string(job.EncodedArgs)).Str("panic_val", fmt.Sprintf("%v", panicVal)).Str("trace", trace).Msg("task error")

	return nil
}
