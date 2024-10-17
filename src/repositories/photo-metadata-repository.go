package repositories

import (
	"context"
	"database/sql"
	"log"
	"photo-service/src/interfaces"
	"photo-service/src/internal/database"
)

type PhotoMetadataRepo struct {
	db *database.Queries
}

func NewPhotoMetadataRepo(db *database.Queries) *PhotoMetadataRepo {
	return &PhotoMetadataRepo{db: db}
}

func (r *PhotoMetadataRepo) CreatePhotoMetadata(ctx context.Context, request interfaces.CreatePhotoMetadataRepoRequest) (string, error) {
	metadata, err := r.db.CreatePhotoMetadata(ctx, database.CreatePhotoMetadataParams{
		ID:        request.Id,
		Column2:   *request.Latitude,
		Column3:   *request.Longitude,
		CreatedAt: sql.NullTime{Time: *request.CreatedAt, Valid: true},
	})
	if err != nil {
		log.Printf("Error creating photo metadata: %v", err)
		return "", err
	}
	return metadata.ID.String(), nil
}
