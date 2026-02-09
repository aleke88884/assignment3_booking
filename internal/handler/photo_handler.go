package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"smartbooking/internal/service"
)

type PhotoHandler struct {
	photoService service.PhotoService
}

func NewPhotoHandler(photoService service.PhotoService) *PhotoHandler {
	return &PhotoHandler{
		photoService: photoService,
	}
}

func (h *PhotoHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseMultipartForm(10 << 20)

	resourceIDStr := r.FormValue("resource_id")
	resourceID, err := strconv.ParseInt(resourceIDStr, 10, 64)
	if err != nil || resourceID == 0 {
		http.Error(w, "Invalid resource_id", http.StatusBadRequest)
		return
	}

	isPrimaryStr := r.FormValue("is_primary")
	isPrimary := isPrimaryStr == "true"

	file, handler, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	photo, err := h.photoService.UploadPhoto(r.Context(), resourceID, file, handler.Filename, isPrimary)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error uploading photo: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(photo)
}

func (h *PhotoHandler) GetResourcePhotos(w http.ResponseWriter, r *http.Request) {
	resourceIDStr := r.PathValue("resource_id")
	resourceID, err := strconv.ParseInt(resourceIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid resource_id", http.StatusBadRequest)
		return
	}

	photos, err := h.photoService.GetResourcePhotos(r.Context(), resourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(photos)
}

func (h *PhotoHandler) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid photo id", http.StatusBadRequest)
		return
	}

	err = h.photoService.DeletePhoto(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PhotoHandler) SetPrimaryPhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid photo id", http.StatusBadRequest)
		return
	}

	var req struct {
		ResourceID int64 `json:"resource_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.photoService.SetPrimaryPhoto(r.Context(), id, req.ResourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Primary photo updated"})
}
