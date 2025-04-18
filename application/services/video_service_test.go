package services_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/raffreitas/codeflix-video-encoder/application/repositories"
	"github.com/raffreitas/codeflix-video-encoder/application/services"
	"github.com/raffreitas/codeflix-video-encoder/domain"
	"github.com/raffreitas/codeflix-video-encoder/framework/database"
	"github.com/stretchr/testify/require"
)

func prepare() (*domain.Video, repositories.VideoRepositoryDb) {
	db := database.NewDbTest()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewString()
	video.FilePath = "tears_for_fears.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	return video, repo
}

func TestVideoServiceDownload(t *testing.T) {

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

	err = videoService.Finish()
	require.Nil(t, err)
}
