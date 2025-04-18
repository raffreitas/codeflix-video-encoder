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

func TestJobRepositoryDbInsert(t *testing.T) {
	db := database.NewDbTest()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewString()
	video.FilePath = "fake_path"
	video.CreatedAt = time.Now()

	job, err := domain.NewJob("output_path", "Pending", video)

	require.Nil(t, err)
	require.NotEmpty(t, job.ID)

	repoJob := repositories.JobRepositoryDb{Db: db}
	repoJob.Insert(job)

	j, err := repoJob.Find(job.ID)
	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
	require.Equal(t, j.ID, job.ID)
	require.Equal(t, j.VideoID, video.ID)
}

func TestJobRepositoryDbUpdate(t *testing.T) {
	db := database.NewDbTest()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewString()
	video.FilePath = "fake_path"
	video.CreatedAt = time.Now()

	job, err := domain.NewJob("output_path", "Pending", video)

	require.Nil(t, err)
	require.NotEmpty(t, job.ID)
	require.Equal(t, job.OutputBucketPath, "output_path")

	repoJob := repositories.JobRepositoryDb{Db: db}
	repoJob.Insert(job)

	job.Status = "Complete"

	repoJob.Update(job)

	j, err := repoJob.Find(job.ID)
	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
	require.Equal(t, j.Status, job.Status)
}
