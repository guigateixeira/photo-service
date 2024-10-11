package repositories

import (
	"context"
	"database/sql"
	"log"

	"photo-service/src/interfaces"
	"photo-service/src/internal/database"
)

type PhotoRepo struct {
	db *database.Queries
}

// Constructor creates a new instance of PhotoRepo.
func NewPhotoRepo(db *database.Queries) *PhotoRepo {
	return &PhotoRepo{db: db}
}

// CreatePhoto creates a new photo entry in the database.
func (r *PhotoRepo) CreatePhoto(ctx context.Context, request interfaces.CreatePhotoRepoRequest) (string, error) {
	photo, err := r.db.CreatePhoto(ctx, database.CreatePhotoParams{
		OwnerID: request.UserID,
		Description: sql.NullString{
			String: request.Description,
			Valid:  request.Description != "", // Set Valid to true if the description is not empty
		},
		PhotoUrl: request.URL,
	})
	if err != nil {
		log.Printf("Error creating photo: %v", err)
		return "", err
	}
	return photo.ID.String(), nil
}
