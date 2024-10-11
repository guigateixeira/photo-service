package services

import (
	"context"
	"log"
	"photo-service/src/interfaces"

	"github.com/google/uuid"
)

type PhotoService struct {
	repo                interfaces.IPhotoRepository
	fileUploaderService interfaces.IFileUpload
}

func NewPhotoService(repo interfaces.IPhotoRepository, fileUploaderService interfaces.IFileUpload) *PhotoService {
	return &PhotoService{repo: repo, fileUploaderService: fileUploaderService}
}

func (s *PhotoService) CreatePhoto(ctx context.Context, request interfaces.CreatePhotoRequest) (string, error) {
	// Upload the file to the file storage service
	photoId := uuid.New()
	uploadRequest := interfaces.UploadFileRequest{
		UserID:   request.UserID.String(),
		Id:       photoId.String(),
		FileName: request.FileName,
		FileData: request.FileData,
	}
	url, err := s.fileUploaderService.Upload(ctx, uploadRequest)
	if err != nil {
		log.Printf("Error uploading file to S3: %v", err)
		return "", err
	}
	req := interfaces.CreatePhotoRepoRequest{
		UserID:      request.UserID,
		Description: request.Description,
		URL:         url,
	}
	return s.repo.CreatePhoto(ctx, req)
}
