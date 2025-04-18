package domain_test

import (
	"testing"
	"time"

	"github.com/raffreitas/codeflix-video-encoder/domain"
	"github.com/stretchr/testify/require"

	uuid "github.com/google/uuid"
)

func TestValidateIfVideoIsEmpty(t *testing.T) {
	video := domain.NewVideo()
	err := video.Validate()

	require.Error(t, err)
}

func TestVideoIdIsNotAUuid(t *testing.T) {
	video := domain.NewVideo()

	video.ID = "not_a_uuid"
	video.ResourceId = "fake_resource_id"
	video.FilePath = "fake_path"
	video.CreatedAt = time.Now()

	err := video.Validate()

	require.Error(t, err)
}

func TestVideoValidation(t *testing.T) {
	video := domain.NewVideo()

	video.ID = uuid.NewString()
	video.ResourceId = "fake_resource_id"
	video.FilePath = "fake_path"
	video.CreatedAt = time.Now()

	err := video.Validate()

	require.Nil(t, err)
}
