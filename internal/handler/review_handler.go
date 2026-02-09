package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"smartbooking/internal/service"
)

type ReviewHandler struct {
	reviewService service.ReviewService
}

func NewReviewHandler(reviewService service.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
	}
}

type CreateReviewRequest struct {
	UserID     int64  `json:"user_id"`
	ResourceID int64  `json:"resource_id"`
	BookingID  *int64 `json:"booking_id,omitempty"`
	Rating     int    `json:"rating"`
	Comment    string `json:"comment"`
}

type UpdateReviewRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

// Create handles POST /reviews
// @Summary Create a new review
// @Description Create a new review for a resource
// @Tags reviews
// @Accept json
// @Produce json
// @Param request body CreateReviewRequest true "Review details"
// @Success 201 {object} models.Review
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router /reviews [post]
func (h *ReviewHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	review, err := h.reviewService.Create(r.Context(), req.UserID, req.ResourceID, req.BookingID, req.Rating, req.Comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}

// GetByID handles GET /reviews/{id}
// @Summary Get review by ID
// @Description Get details of a specific review
// @Tags reviews
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} models.Review
// @Failure 400 {string} string "Invalid review ID"
// @Failure 404 {string} string "Review not found"
// @Failure 500 {string} string "Internal server error"
// @Router /reviews/{id} [get]
func (h *ReviewHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	review, err := h.reviewService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if review == nil {
		http.Error(w, "Review not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

// GetByResource handles GET /resources/{resource_id}/reviews
// @Summary Get reviews for a resource
// @Description Get all reviews for a specific resource
// @Tags reviews
// @Produce json
// @Param resource_id path int true "Resource ID"
// @Success 200 {array} models.Review
// @Failure 400 {string} string "Invalid resource ID"
// @Failure 500 {string} string "Internal server error"
// @Router /resources/{resource_id}/reviews [get]
func (h *ReviewHandler) GetByResource(w http.ResponseWriter, r *http.Request) {
	resourceIDStr := r.PathValue("resource_id")
	resourceID, err := strconv.ParseInt(resourceIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid resource ID", http.StatusBadRequest)
		return
	}

	reviews, err := h.reviewService.GetByResource(r.Context(), resourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

// GetByUser handles GET /users/{user_id}/reviews
// @Summary Get reviews by user
// @Description Get all reviews written by a specific user
// @Tags reviews
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {array} models.Review
// @Failure 400 {string} string "Invalid user ID"
// @Failure 500 {string} string "Internal server error"
// @Router /users/{user_id}/reviews [get]
func (h *ReviewHandler) GetByUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	reviews, err := h.reviewService.GetByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

// Update handles PUT /reviews/{id}
// @Summary Update a review
// @Description Update an existing review
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Param request body UpdateReviewRequest true "Updated review details"
// @Success 200 {object} models.Review
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /reviews/{id} [put]
func (h *ReviewHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	var req UpdateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	review, err := h.reviewService.Update(r.Context(), id, req.Rating, req.Comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

// Delete handles DELETE /reviews/{id}
// @Summary Delete a review
// @Description Delete a review by ID
// @Tags reviews
// @Param id path int true "Review ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid review ID"
// @Failure 500 {string} string "Internal server error"
// @Router /reviews/{id} [delete]
func (h *ReviewHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	if err := h.reviewService.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetResourceAverageRating handles GET /resources/{resource_id}/rating
// @Summary Get average rating for a resource
// @Description Get the average rating for a specific resource
// @Tags reviews
// @Produce json
// @Param resource_id path int true "Resource ID"
// @Success 200 {object} map[string]float64
// @Failure 400 {string} string "Invalid resource ID"
// @Failure 500 {string} string "Internal server error"
// @Router /resources/{resource_id}/rating [get]
func (h *ReviewHandler) GetResourceAverageRating(w http.ResponseWriter, r *http.Request) {
	resourceIDStr := r.PathValue("resource_id")
	resourceID, err := strconv.ParseInt(resourceIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid resource ID", http.StatusBadRequest)
		return
	}

	avgRating, err := h.reviewService.GetResourceAverageRating(r.Context(), resourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"average_rating": avgRating})
}
