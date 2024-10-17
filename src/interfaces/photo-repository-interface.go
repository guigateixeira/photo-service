package interfaces

import (
	"context"

	"github.com/google/uuid"
)

type CreatePhotoRepoRequest struct {
	UserID      uuid.UUID
	Description string
	URL         string
}

type IPhotoRepository interface {
	CreatePhoto(ctx context.Context, req CreatePhotoRepoRequest) (string, error)
}
