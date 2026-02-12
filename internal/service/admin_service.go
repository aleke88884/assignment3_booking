package service

import (
	"context"

	"smartbooking/internal/repository"
)

// AdminService defines the interface for admin-specific business logic
type AdminService interface {
	GetSystemStatistics(ctx context.Context) (*repository.AdminStatistics, error)
	GetBookingsByStatus(ctx context.Context) ([]repository.BookingStatusCount, error)
	GetResourcesByCategory(ctx context.Context) ([]repository.CategoryResourceCount, error)
	GetRevenueByMonth(ctx context.Context, months int) ([]repository.MonthlyRevenue, error)
	GetBookingsByDay(ctx context.Context, days int) ([]repository.DailyBookings, error)
}

// adminService implements AdminService interface
type adminService struct {
	adminRepo repository.AdminRepository
}

// NewAdminService creates a new instance of AdminService
func NewAdminService(adminRepo repository.AdminRepository) AdminService {
	return &adminService{
		adminRepo: adminRepo,
	}
}

// GetSystemStatistics retrieves system-wide statistics
func (s *adminService) GetSystemStatistics(ctx context.Context) (*repository.AdminStatistics, error) {
	return s.adminRepo.GetSystemStatistics(ctx)
}

// GetBookingsByStatus retrieves booking count grouped by status
func (s *adminService) GetBookingsByStatus(ctx context.Context) ([]repository.BookingStatusCount, error) {
	return s.adminRepo.GetBookingsByStatus(ctx)
}

// GetResourcesByCategory retrieves resource count grouped by category
func (s *adminService) GetResourcesByCategory(ctx context.Context) ([]repository.CategoryResourceCount, error) {
	return s.adminRepo.GetResourcesByCategory(ctx)
}

// GetRevenueByMonth retrieves revenue aggregated by month
func (s *adminService) GetRevenueByMonth(ctx context.Context, months int) ([]repository.MonthlyRevenue, error) {
	if months <= 0 {
		months = 6 // default to 6 months
	}
	return s.adminRepo.GetRevenueByMonth(ctx, months)
}

// GetBookingsByDay retrieves booking count aggregated by day
func (s *adminService) GetBookingsByDay(ctx context.Context, days int) ([]repository.DailyBookings, error) {
	if days <= 0 {
		days = 7 // default to 7 days
	}
	return s.adminRepo.GetBookingsByDay(ctx, days)
}
