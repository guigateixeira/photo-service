package interfaces

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CreatePhotoMetadataRepoRequest struct {
	Id        uuid.UUID
	Longitude *float64
	Latitude  *float64
	CreatedAt *time.Time
}

type IPhotoMetadataRepository interface {
	CreatePhotoMetadata(ctx context.Context, req CreatePhotoMetadataRepoRequest) (string, error)
}
