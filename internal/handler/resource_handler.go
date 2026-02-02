package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"smartbooking/internal/models"
	"smartbooking/internal/service"
)

// ResourceHandler handles resource-related HTTP requests
type ResourceHandler struct {
	resourceService service.ResourceService
}

// NewResourceHandler creates a new ResourceHandler instance
func NewResourceHandler(resourceService service.ResourceService) *ResourceHandler {
	return &ResourceHandler{
		resourceService: resourceService,
	}
}

// CreateResourceRequest represents the create resource request body
type CreateResourceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Capacity    int    `json:"capacity"`
}

// Create handles POST /resources
func (h *ResourceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resource := &models.Resource{
		Name:        req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.resourceService.Create(r.Context(), resource); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resource)
}

// GetByID handles GET /resources/{id}
func (h *ResourceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid resource ID", http.StatusBadRequest)
		return
	}

	resource, err := h.resourceService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if resource == nil {
		http.Error(w, "Resource not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resource)
}

// List handles GET /resources
func (h *ResourceHandler) List(w http.ResponseWriter, r *http.Request) {
	resources, err := h.resourceService.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resources)
}

// Delete handles DELETE /resources/{id}
func (h *ResourceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid resource ID", http.StatusBadRequest)
		return
	}

	if err := h.resourceService.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
