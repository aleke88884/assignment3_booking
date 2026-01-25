package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"smartbooking/internal/service"
)

// BookingHandler handles booking-related HTTP requests
type BookingHandler struct {
	bookingService service.BookingService
}

// NewBookingHandler creates a new BookingHandler instance
func NewBookingHandler(bookingService service.BookingService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
	}
}

// CreateBookingRequest represents the create booking request body
type CreateBookingRequest struct {
	UserID     int64  `json:"user_id"`
	ResourceID int64  `json:"resource_id"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}

// Create handles POST /bookings
func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		http.Error(w, "Invalid start_time format", http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		http.Error(w, "Invalid end_time format", http.StatusBadRequest)
		return
	}

	booking, err := h.bookingService.Create(r.Context(), req.UserID, req.ResourceID, startTime, endTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

// GetByID handles GET /bookings/{id}
func (h *BookingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	booking, err := h.bookingService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if booking == nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booking)
}

// Cancel handles POST /bookings/{id}/cancel
func (h *BookingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	if err := h.bookingService.Cancel(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Booking cancelled successfully"}`))
}

// ListAll handles GET /bookings
func (h *BookingHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	bookings, err := h.bookingService.ListAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

// ListByUser handles GET /users/{id}/bookings
func (h *BookingHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	bookings, err := h.bookingService.ListByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}
