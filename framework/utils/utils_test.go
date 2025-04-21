package utils_test

import (
	"testing"

	"github.com/raffreitas/codeflix-video-encoder/framework/utils"
	"github.com/stretchr/testify/require"
)

func TestIsJsonValid(t *testing.T) {
	json := `
		{
			"resource_id": "52d8ef70-e795-41ab-9846-2523bf34c30e",
			"file_path": "file.mp4",
			"status": "pending"
		}
	`

	err := utils.IsJson(json)
	require.NoError(t, err)
}

func TestIsJsonInvalid(t *testing.T) {
	json := `some invalid json`

	err := utils.IsJson(json)
	require.Error(t, err)
}
