package interfaces

import (
	"context"

	"github.com/google/uuid"
)

// CreatePhotoRequest represents the data needed to create a photo.
type CreatePhotoRepoRequest struct {
	UserID      uuid.UUID
	Description string
	URL         string
}

// IPhotoRepository defines methods for managing photos.
type IPhotoRepository interface {
	CreatePhoto(ctx context.Context, req CreatePhotoRepoRequest) (string, error)
}
