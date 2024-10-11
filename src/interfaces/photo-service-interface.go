package interfaces

import (
	"context"

	"github.com/google/uuid"
)

type CreatePhotoRequest struct {
	UserID      uuid.UUID
	Description string
	FileName    string
	FileData    []byte
}

type IPhotoService interface {
	CreatePhoto(ctx context.Context, request CreatePhotoRequest) (string, error)
}
