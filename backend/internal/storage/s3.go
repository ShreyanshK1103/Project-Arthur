package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3Client() (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}

func TestS3Connection() error {
	client, err := NewS3Client()
	if err != nil {
		return err
	}

	result, err := client.ListBuckets(
		context.Background(),
		&s3.ListBucketsInput{},
	)
	if err != nil {
		return err
	}

	for _, bucket := range result.Buckets {
		fmt.Println(*bucket.Name)
	}

	return nil
}

func UploadDirectory(
	deploymentID string,
	localPath string,
) error {

	client, err := NewS3Client()
	if err != nil {
		return err
	}

	uploader := manager.NewUploader(client)

	bucketName := os.Getenv(
		"AWS_BUCKET_NAME",
	)

	return filepath.Walk(
		localPath,
		func(
			path string,
			info os.FileInfo,
			err error,
		) error {

			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			relativePath, err := filepath.Rel(
				localPath,
				path,
			)
			if err != nil {
				return err
			}

			s3Key := filepath.Join(
				deploymentID,
				relativePath,
			)

			contentType := "application/octet-stream"
				switch filepath.Ext(path) {
				case ".html":
					contentType = "text/html"
				case ".css":
					contentType = "text/css"
				case ".js":
					contentType = "application/javascript"
				case ".json":
					contentType = "application/json"
				case ".svg":
					contentType = "image/svg+xml"
			}

			_, err = uploader.Upload(
				context.Background(),
				&s3.PutObjectInput{
					Bucket:      &bucketName,
					Key:         &s3Key,
					Body:        file,
					ContentType: &contentType,
				},
			)

			return err
		},
	)
}

func GetObject (key string,) (*s3.GetObjectOutput, error) {
	client, err := NewS3Client()
	if err != nil {
		return nil, err
	}

	bucketName := os.Getenv(
		"AWS_BUCKET_NAME",
	)

	return client.GetObject(
		context.Background(),
		&s3.GetObjectInput{
			Bucket: &bucketName,
			Key: &key,
		},
	)
}

func GetDeploymentFile(deploymentID string, requestPath string,) (*s3.GetObjectOutput, error) {
	if requestPath == "/" {
		requestPath = "/index.html"
	}

	key := deploymentID + requestPath

	return GetObject(key)
}