package services

import (
	"bytes"
	"context"
	"log"
	"mime"
	"path/filepath"
	"photo-service/src/interfaces"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Uploader struct {
	client *s3.Client
	bucket string
}

func NewS3Uploader(client *s3.Client, bucket string) *S3Uploader {
	return &S3Uploader{
		client: client,
		bucket: bucket,
	}
}

func (u *S3Uploader) Upload(ctx context.Context, request interfaces.UploadFileRequest) (string, error) {
	contentType := mime.TypeByExtension(filepath.Ext(request.FileName))
	if contentType == "" {
		contentType = "application/octet-stream" // Default to binary if unknown
	}
	path := request.UserID + "/" + request.Id + "--" + request.FileName
	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(u.bucket),
		Key:         aws.String(path),
		Body:        bytes.NewReader(request.FileData),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Println("Failed to upload file to S3:", err)
		return "", err
	}

	return path, nil
}
