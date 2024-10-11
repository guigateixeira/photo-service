-- +goose Up
CREATE EXTENSION postgis;

CREATE TABLE photo_metadata (
    id UUID PRIMARY KEY,
    location GEOGRAPHY(Point, 4326),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_photo
        FOREIGN KEY (id) 
        REFERENCES photo (id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE photo_metadata;
DROP EXTENSION IF EXISTS postgis;