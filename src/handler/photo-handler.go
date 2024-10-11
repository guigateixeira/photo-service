package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"photo-service/src/interfaces"

	"github.com/google/uuid"
)

type PhotoHandler struct {
	photoService interfaces.IPhotoService
}

func NewPhotoHandler(photoService interfaces.IPhotoService) *PhotoHandler {
	return &PhotoHandler{photoService: photoService}
}

type CreatePhotoRequest struct {
	UserID      string `json:"user_id"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

func (h *PhotoHandler) CreatePhoto(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil { // Limit to 10MB
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Extract fields from the form data
	userIDStr := r.FormValue("userId")
	description := r.FormValue("description")

	// Validate and parse the user ID as a UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Retrieve the file from the form data
	file, handler, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	originalFilename := handler.Filename
	fileBytes, err := h.fileToBytes(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	serviceRequest := interfaces.CreatePhotoRequest{
		UserID:      userID,
		Description: description,
		FileName:    originalFilename,
		FileData:    fileBytes,
	}

	photoID, err := h.photoService.CreatePhoto(r.Context(), serviceRequest)
	if err != nil {
		http.Error(w, "Error creating photo", http.StatusInternalServerError)
		return
	}

	// Return a success response with the photo ID
	response := map[string]string{"photo_id": photoID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *PhotoHandler) fileToBytes(file multipart.File) ([]byte, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	return fileBytes, nil
}
