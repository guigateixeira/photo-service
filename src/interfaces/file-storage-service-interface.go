package interfaces

import "context"

type UploadFileRequest struct {
	UserID   string
	Id       string
	FileName string
	FileData []byte
}

type IFileUpload interface {
	Upload(ctx context.Context, request UploadFileRequest) (string, error)
}
