package services_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/raffreitas/codeflix-video-encoder/application/services"
	"github.com/stretchr/testify/require"
)

func TestVideoServiceUpload(t *testing.T) {

	video, repo := prepare()

	videoService := services.NewVideoService()
	videoService.VideoRepository = repo
	videoService.Video = video

	err := videoService.Download("codeflix-videos")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	videoUpload := services.NewVideoUpload()
	videoUpload.OutputBucket = "codeflix-videos"
	videoUpload.VideoPath = filepath.Join(os.TempDir(), video.ID)

	doneUpload := make(chan string)
	go videoUpload.ProcessUpload(50, doneUpload)

	result := <-doneUpload

	require.Equal(t, "upload completed", result)
}
