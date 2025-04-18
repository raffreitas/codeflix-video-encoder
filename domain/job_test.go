package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/raffreitas/codeflix-video-encoder/domain"
	"github.com/stretchr/testify/require"
)

func TestNewJob(t *testing.T) {
	video := domain.NewVideo()
	video.ID = uuid.NewString()
	video.FilePath = "fake_path"
	video.CreatedAt = time.Now()

	job, err := domain.NewJob("fake_path", "Converted", video)

	require.NotNil(t, job)
	require.Nil(t, err)
}
