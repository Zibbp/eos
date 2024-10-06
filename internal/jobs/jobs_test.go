package jobs_test

import (
	"context"
	"testing"

	"github.com/riverqueue/river"
	"github.com/stretchr/testify/assert"
	"github.com/zibbp/avalon/internal/jobs"
)

func TestVideoProcessWorker(t *testing.T) {
	ctx := context.Background()

	err := (&jobs.VideoProcessWorker{}).Work(ctx, &river.Job[jobs.VideoProcessArgs]{Args: jobs.VideoProcessArgs{
		VideoInfoPath: "/data/videos/info.json",
	}})

	assert.NoError(t, err)
}
