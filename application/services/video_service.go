package services

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/raffreitas/codeflix-video-encoder/application/repositories"
	"github.com/raffreitas/codeflix-video-encoder/domain"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (v *VideoService) Download(bucketName string) error {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	s3Options := s3.Options{
		Region:      cfg.Region,
		Credentials: cfg.Credentials,
	}

	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	if awsEndpoint != "" {
		s3Options.BaseEndpoint = aws.String(awsEndpoint)
		s3Options.UsePathStyle = true
	}

	client := s3.New(s3Options)

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(v.Video.FilePath),
	}

	result, err := client.GetObject(ctx, input)
	if err != nil {
		return err
	}
	defer result.Body.Close()

	tempDir := os.TempDir()
	filePath := filepath.Join(tempDir, v.Video.ID+".mp4")

	output, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, result.Body)
	if err != nil {
		return err
	}

	return nil
}

func (v *VideoService) Fragment() error {
	tempDir := os.TempDir()
	filePath := filepath.Join(tempDir, v.Video.ID)
	err := os.Mkdir(filePath, os.ModePerm)
	if err != nil {
		return err
	}
	source := filepath.Join(filePath + ".mp4")
	target := filepath.Join(filePath + ".frag")

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	printOutput(output)
	return nil
}

func (v *VideoService) Encode() error {
	tempDir := os.TempDir()
	filePath := filepath.Join(tempDir, v.Video.ID)

	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, filePath+".frag")
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, filePath)
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")
	cmd := exec.Command("mp4dash", cmdArgs...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Finish() error {
	tempDir := os.TempDir()
	filePath := filepath.Join(tempDir, v.Video.ID)
	err := os.Remove(filePath + ".mp4")
	if err != nil {
		log.Println("error removing mp4 " + v.Video.ID + ".mp4")
		return err
	}
	err = os.Remove(filePath + ".frag")
	if err != nil {
		log.Println("error removing frag " + v.Video.ID + ".frag")
		return err
	}

	err = os.RemoveAll(filePath)
	if err != nil {
		log.Println("error removing dir: " + filePath)
	}

	log.Println("files have been removed: " + v.Video.ID)

	return nil
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("======> Output: %s\n", string(out))
	}
}
