package services

import (
	"context"
	"fmt"
	"log"
	"photo-service/src/interfaces"
	"strconv"
	"strings"
	"time"

	"github.com/dsoprea/go-exif/v3"
	"github.com/google/uuid"
)

type PhotoService struct {
	repo                interfaces.IPhotoRepository
	fileUploaderService interfaces.IFileUpload
	photoMetadataRepo   interfaces.IPhotoMetadataRepository
}

func NewPhotoService(
	repo interfaces.IPhotoRepository,
	fileUploaderService interfaces.IFileUpload,
	photoMetadataRepo interfaces.IPhotoMetadataRepository,
) *PhotoService {
	return &PhotoService{repo: repo, fileUploaderService: fileUploaderService, photoMetadataRepo: photoMetadataRepo}
}

func (s *PhotoService) CreatePhoto(ctx context.Context, request interfaces.CreatePhotoRequest) (string, error) {
	uniqueId := uuid.New()
	uploadRequest := interfaces.UploadFileRequest{
		UserID:   request.UserID.String(),
		Id:       uniqueId.String(),
		FileName: request.FileName,
		FileData: request.FileData,
	}
	url, err := s.fileUploaderService.Upload(ctx, uploadRequest)
	if err != nil {
		log.Printf("Error uploading file to S3: %v", err)
		return "", err
	}
	req := interfaces.CreatePhotoRepoRequest{
		UserID:      request.UserID,
		Description: request.Description,
		URL:         url,
	}
	photoId, err := s.repo.CreatePhoto(ctx, req)
	if err != nil {
		return "", err
	}
	lat, long, time, err2 := extractExifData(request.FileData)
	if err2 != nil {
		log.Printf("Error extracting EXIF data: %v", err2)
	}
	if lat != 0 || long != 0 {
		log.Printf("EXIF data found. Creating photo metadata... lat %v, long %v, time %v", lat, long, time)
		photoUUID, err3 := uuid.Parse(photoId)
		if err3 != nil {
			log.Printf("Error parsing photo UUID: %v", err3)
		}
		req := interfaces.CreatePhotoMetadataRepoRequest{
			Id:        photoUUID,
			Latitude:  &lat,
			Longitude: &long,
			CreatedAt: &time,
		}
		_, err = s.photoMetadataRepo.CreatePhotoMetadata(ctx, req)
		if err != nil {
			log.Printf("Error creating photo metadata: %v", err)
			return "", err
		}
	}
	return photoId, err
}

// Extract EXIF data from the image file bytes
func extractExifData(fileBytes []byte) (float64, float64, time.Time, error) {
	var latitude, longitude float64
	var createdAt time.Time
	var latSign = 1
	var longSign = 1

	rawExif, err := exif.SearchAndExtractExif(fileBytes)
	if err != nil {
		if err.Error() == "no exif data" {
			log.Printf("No EXIF data found in the image")
			return latitude, longitude, createdAt, err
		}
		log.Printf("Error extracting EXIF: %v", err)
		return latitude, longitude, createdAt, fmt.Errorf("error extracting EXIF: %v", err)
	}

	log.Printf("EXIF data found. Parsing...")

	entries, _, err := exif.GetFlatExifDataUniversalSearch(rawExif, nil, true)
	if err != nil {
		log.Printf("Error getting flat EXIF data: %v", err)
		return latitude, longitude, createdAt, fmt.Errorf("error getting flat EXIF data: %v", err)
	}

	for _, entry := range entries {
		switch entry.TagName {
		case "GPSLatitude":
			latitude, err = parseGPSCoordinate(entry.Formatted)
			if err != nil {
				log.Printf("Error parsing latitude: %v", err)
			}
		case "GPSLongitude":
			longitude, err = parseGPSCoordinate(entry.Formatted)
			if err != nil {
				log.Printf("Error parsing longitude: %v", err)
			}
		case "DateTimeOriginal":
			createdAt, err = time.Parse("2006:01:02 15:04:05", entry.FormattedFirst)
			if err != nil {
				log.Printf("Error parsing creation date: %v", err)
			}
		case "GPSLongitudeRef":
			if entry.Value == "W" {
				longSign = -1
			}
		case "GPSLatitudeRef":
			if entry.Value == "S" {
				latSign = -1
			}
		}
	}

	if latitude == 0 && longitude == 0 {
		log.Printf("GPS coordinates not found in EXIF data")
	}
	if createdAt.IsZero() {
		log.Printf("Creation date not found in EXIF data")
	}

	latitude *= float64(latSign)
	longitude *= float64(longSign)
	return latitude, longitude, createdAt, nil
}

// Helper function to parse GPS coordinates
func parseGPSCoordinate(latStr string) (float64, error) {
	// Remove the square brackets and split the string
	latStr = strings.Trim(latStr, "[]")
	parts := strings.Split(latStr, " ")

	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid latitude format")
	}

	// Parse degrees, minutes, and seconds
	degrees, err := parseFraction(parts[0])
	if err != nil {
		return 0, err
	}

	minutes, err := parseFraction(parts[1])
	if err != nil {
		return 0, err
	}

	seconds, err := parseFraction(parts[2])
	if err != nil {
		return 0, err
	}

	// Convert to decimal degrees
	decimalDegrees := degrees + (minutes / 60) + (seconds / 3600)
	return decimalDegrees, nil
}

// Helper function to parse fractions (e.g., "20/1")
func parseFraction(fraction string) (float64, error) {
	parts := strings.Split(fraction, "/")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid fraction format")
	}

	numerator, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, err
	}

	denominator, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, err
	}

	return numerator / denominator, nil
}
