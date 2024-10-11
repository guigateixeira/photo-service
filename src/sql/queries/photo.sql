-- name: CreatePhoto :one
INSERT INTO photos (owner_id, description, photo_url)
VALUES ($1, $2, $3)
RETURNING *;
