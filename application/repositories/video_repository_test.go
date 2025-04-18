package repositories_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/raffreitas/codeflix-video-encoder/application/repositories"
	"github.com/raffreitas/codeflix-video-encoder/domain"
	"github.com/raffreitas/codeflix-video-encoder/framework/database"
	"github.com/stretchr/testify/require"
)

func TestVideoRepositoryDbInsert(t *testing.T) {
	db := database.NewDbTest()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewString()
	video.FilePath = "fake_path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	v, err := repo.Find(video.ID)

	require.NotEmpty(t, v.ID)
	require.Nil(t, err)
	require.Equal(t, v.ID, video.ID)
}
