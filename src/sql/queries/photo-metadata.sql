---- name: CreatePhotoMetadata :one
INSERT INTO photo_metadata (id, location, created_at)
VALUES ($1, ST_MakePoint($2::double precision, $3::double precision), $4)
RETURNING *;
