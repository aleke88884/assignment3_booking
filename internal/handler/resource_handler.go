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
// @Summary Create a new resource
// @Description Create a new bookable resource (room, apartment, facility)
// @Tags resources
// @Accept json
// @Produce json
// @Param request body CreateResourceRequest true "Resource details"
// @Success 201 {object} models.Resource
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router /resources [post]
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
// @Summary Get resource by ID
// @Description Get details of a specific resource
// @Tags resources
// @Produce json
// @Param id path int true "Resource ID"
// @Success 200 {object} models.Resource
// @Failure 400 {string} string "Invalid resource ID"
// @Failure 404 {string} string "Resource not found"
// @Failure 500 {string} string "Internal server error"
// @Router /resources/{id} [get]
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
// @Summary List all resources
// @Description Get a list of all available resources
// @Tags resources
// @Produce json
// @Success 200 {array} models.Resource
// @Failure 500 {string} string "Internal server error"
// @Router /resources [get]
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
// @Summary Delete a resource
// @Description Delete a resource by ID
// @Tags resources
// @Param id path int true "Resource ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid resource ID"
// @Failure 500 {string} string "Internal server error"
// @Router /resources/{id} [delete]
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
