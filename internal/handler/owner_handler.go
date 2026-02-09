package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"smartbooking/internal/service"
)

// OwnerHandler handles HTTP requests for owner operations
type OwnerHandler struct {
	ownerService service.OwnerService
}

// NewOwnerHandler creates a new OwnerHandler
func NewOwnerHandler(ownerService service.OwnerService) *OwnerHandler {
	return &OwnerHandler{
		ownerService: ownerService,
	}
}

// GetOwnerResources godoc
// @Summary Get owner's resources
// @Description Retrieves all resources owned by a specific owner
// @Tags owners
// @Produce json
// @Param id path int true "Owner ID"
// @Success 200 {array} models.Resource
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/owners/{id}/resources [get]
func (h *OwnerHandler) GetOwnerResources(w http.ResponseWriter, r *http.Request) {
	ownerIDStr := r.PathValue("id")
	ownerID, err := strconv.ParseInt(ownerIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "Invalid owner ID"}`, http.StatusBadRequest)
		return
	}

	resources, err := h.ownerService.GetOwnerResources(r.Context(), ownerID)
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch owner resources"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resources)
}

// GetOwnerBookings godoc
// @Summary Get bookings for owner's resources
// @Description Retrieves all bookings for resources owned by a specific owner
// @Tags owners
// @Produce json
// @Param id path int true "Owner ID"
// @Success 200 {array} models.Booking
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/owners/{id}/bookings [get]
func (h *OwnerHandler) GetOwnerBookings(w http.ResponseWriter, r *http.Request) {
	ownerIDStr := r.PathValue("id")
	ownerID, err := strconv.ParseInt(ownerIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "Invalid owner ID"}`, http.StatusBadRequest)
		return
	}

	bookings, err := h.ownerService.GetOwnerBookings(r.Context(), ownerID)
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch owner bookings"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

// GetOwnerStatistics godoc
// @Summary Get owner statistics
// @Description Retrieves aggregated statistics for an owner (total resources, bookings, revenue, ratings, etc.)
// @Tags owners
// @Produce json
// @Param id path int true "Owner ID"
// @Success 200 {object} repository.OwnerStatistics
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/owners/{id}/statistics [get]
func (h *OwnerHandler) GetOwnerStatistics(w http.ResponseWriter, r *http.Request) {
	ownerIDStr := r.PathValue("id")
	ownerID, err := strconv.ParseInt(ownerIDStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "Invalid owner ID"}`, http.StatusBadRequest)
		return
	}

	stats, err := h.ownerService.GetOwnerStatistics(r.Context(), ownerID)
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch owner statistics"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
