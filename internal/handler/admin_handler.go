package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"smartbooking/internal/service"
)

// AdminHandler handles HTTP requests for admin operations
type AdminHandler struct {
	adminService service.AdminService
}

// NewAdminHandler creates a new AdminHandler
func NewAdminHandler(adminService service.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// GetSystemStatistics godoc
// @Summary Get system-wide statistics
// @Description Retrieves aggregated statistics for the entire system (users, resources, bookings, revenue, etc.)
// @Tags admin
// @Produce json
// @Success 200 {object} repository.AdminStatistics
// @Failure 500 {object} map[string]string
// @Router /api/admin/statistics [get]
func (h *AdminHandler) GetSystemStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := h.adminService.GetSystemStatistics(r.Context())
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch system statistics"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetBookingsByStatus godoc
// @Summary Get bookings grouped by status
// @Description Retrieves booking counts grouped by status (pending, confirmed, cancelled)
// @Tags admin
// @Produce json
// @Success 200 {array} repository.BookingStatusCount
// @Failure 500 {object} map[string]string
// @Router /api/admin/bookings/by-status [get]
func (h *AdminHandler) GetBookingsByStatus(w http.ResponseWriter, r *http.Request) {
	data, err := h.adminService.GetBookingsByStatus(r.Context())
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch bookings by status"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetResourcesByCategory godoc
// @Summary Get resources grouped by category
// @Description Retrieves resource counts grouped by category
// @Tags admin
// @Produce json
// @Success 200 {array} repository.CategoryResourceCount
// @Failure 500 {object} map[string]string
// @Router /api/admin/resources/by-category [get]
func (h *AdminHandler) GetResourcesByCategory(w http.ResponseWriter, r *http.Request) {
	data, err := h.adminService.GetResourcesByCategory(r.Context())
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch resources by category"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetRevenueByMonth godoc
// @Summary Get revenue by month
// @Description Retrieves revenue aggregated by month for the last N months
// @Tags admin
// @Produce json
// @Param months query int false "Number of months" default(6)
// @Success 200 {array} repository.MonthlyRevenue
// @Failure 500 {object} map[string]string
// @Router /api/admin/revenue/by-month [get]
func (h *AdminHandler) GetRevenueByMonth(w http.ResponseWriter, r *http.Request) {
	months := 6 // default
	if monthsParam := r.URL.Query().Get("months"); monthsParam != "" {
		if m, err := strconv.Atoi(monthsParam); err == nil && m > 0 {
			months = m
		}
	}

	data, err := h.adminService.GetRevenueByMonth(r.Context(), months)
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch revenue by month"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetBookingsByDay godoc
// @Summary Get bookings by day
// @Description Retrieves booking counts aggregated by day for the last N days
// @Tags admin
// @Produce json
// @Param days query int false "Number of days" default(7)
// @Success 200 {array} repository.DailyBookings
// @Failure 500 {object} map[string]string
// @Router /api/admin/bookings/by-day [get]
func (h *AdminHandler) GetBookingsByDay(w http.ResponseWriter, r *http.Request) {
	days := 7 // default
	if daysParam := r.URL.Query().Get("days"); daysParam != "" {
		if d, err := strconv.Atoi(daysParam); err == nil && d > 0 {
			days = d
		}
	}

	data, err := h.adminService.GetBookingsByDay(r.Context(), days)
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch bookings by day"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
