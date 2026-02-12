package repository

import (
	"context"
	"database/sql"
	"time"
)

// AdminRepository defines the interface for admin-specific statistics and operations
type AdminRepository interface {
	GetSystemStatistics(ctx context.Context) (*AdminStatistics, error)
	GetBookingsByStatus(ctx context.Context) ([]BookingStatusCount, error)
	GetResourcesByCategory(ctx context.Context) ([]CategoryResourceCount, error)
	GetRevenueByMonth(ctx context.Context, months int) ([]MonthlyRevenue, error)
	GetBookingsByDay(ctx context.Context, days int) ([]DailyBookings, error)
}

// AdminStatistics represents system-wide statistics
type AdminStatistics struct {
	TotalUsers         int     `json:"total_users"`
	TotalResources     int     `json:"total_resources"`
	TotalBookings      int     `json:"total_bookings"`
	ActiveBookings     int     `json:"active_bookings"`
	CancelledBookings  int     `json:"cancelled_bookings"`
	TotalRevenue       float64 `json:"total_revenue"`
	TotalReviews       int     `json:"total_reviews"`
	AverageRating      float64 `json:"average_rating"`
	TotalCategories    int     `json:"total_categories"`
}

// BookingStatusCount represents booking count by status
type BookingStatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

// CategoryResourceCount represents resource count by category
type CategoryResourceCount struct {
	CategoryName string `json:"category_name"`
	Count        int    `json:"count"`
}

// MonthlyRevenue represents revenue aggregated by month
type MonthlyRevenue struct {
	Month   string  `json:"month"`
	Revenue float64 `json:"revenue"`
}

// DailyBookings represents booking count aggregated by day
type DailyBookings struct {
	Day   string `json:"day"`
	Count int    `json:"count"`
}

// adminRepository implements AdminRepository interface with PostgreSQL storage
type adminRepository struct {
	db *sql.DB
}

// NewAdminRepository creates a new instance of AdminRepository
func NewAdminRepository(db *sql.DB) AdminRepository {
	return &adminRepository{
		db: db,
	}
}

// GetSystemStatistics retrieves system-wide statistics
func (r *adminRepository) GetSystemStatistics(ctx context.Context) (*AdminStatistics, error) {
	stats := &AdminStatistics{}

	// Get counts in parallel
	query := `
		SELECT
			(SELECT COUNT(*) FROM users) as total_users,
			(SELECT COUNT(*) FROM resources) as total_resources,
			(SELECT COUNT(*) FROM bookings) as total_bookings,
			(SELECT COUNT(*) FROM bookings WHERE status IN ('pending', 'confirmed')) as active_bookings,
			(SELECT COUNT(*) FROM bookings WHERE status = 'cancelled') as cancelled_bookings,
			(SELECT COALESCE(SUM(total_price), 0) FROM bookings WHERE status IN ('pending', 'confirmed')) as total_revenue,
			(SELECT COUNT(*) FROM reviews) as total_reviews,
			(SELECT COALESCE(AVG(rating), 0) FROM reviews) as average_rating,
			(SELECT COUNT(*) FROM resource_categories) as total_categories
	`

	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalUsers,
		&stats.TotalResources,
		&stats.TotalBookings,
		&stats.ActiveBookings,
		&stats.CancelledBookings,
		&stats.TotalRevenue,
		&stats.TotalReviews,
		&stats.AverageRating,
		&stats.TotalCategories,
	)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// GetBookingsByStatus retrieves booking count grouped by status
func (r *adminRepository) GetBookingsByStatus(ctx context.Context) ([]BookingStatusCount, error) {
	query := `
		SELECT status, COUNT(*) as count
		FROM bookings
		GROUP BY status
		ORDER BY count DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []BookingStatusCount
	for rows.Next() {
		var item BookingStatusCount
		if err := rows.Scan(&item.Status, &item.Count); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, rows.Err()
}

// GetResourcesByCategory retrieves resource count grouped by category
func (r *adminRepository) GetResourcesByCategory(ctx context.Context) ([]CategoryResourceCount, error) {
	query := `
		SELECT
			COALESCE(c.name, 'Uncategorized') as category_name,
			COUNT(r.id) as count
		FROM resources r
		LEFT JOIN resource_categories c ON r.category_id = c.id
		GROUP BY c.name
		ORDER BY count DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []CategoryResourceCount
	for rows.Next() {
		var item CategoryResourceCount
		if err := rows.Scan(&item.CategoryName, &item.Count); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, rows.Err()
}

// GetRevenueByMonth retrieves revenue aggregated by month for the last N months
func (r *adminRepository) GetRevenueByMonth(ctx context.Context, months int) ([]MonthlyRevenue, error) {
	query := `
		SELECT
			TO_CHAR(DATE_TRUNC('month', created_at), 'Mon YYYY') as month,
			COALESCE(SUM(total_price), 0) as revenue
		FROM bookings
		WHERE status IN ('pending', 'confirmed')
			AND created_at >= NOW() - INTERVAL '1 month' * $1
		GROUP BY DATE_TRUNC('month', created_at)
		ORDER BY DATE_TRUNC('month', created_at)
	`

	rows, err := r.db.QueryContext(ctx, query, months)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []MonthlyRevenue
	for rows.Next() {
		var item MonthlyRevenue
		if err := rows.Scan(&item.Month, &item.Revenue); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	// If no data, fill with zeros for the last N months
	if len(result) == 0 {
		now := time.Now()
		for i := months - 1; i >= 0; i-- {
			month := now.AddDate(0, -i, 0)
			result = append(result, MonthlyRevenue{
				Month:   month.Format("Jan 2006"),
				Revenue: 0,
			})
		}
	}

	return result, rows.Err()
}

// GetBookingsByDay retrieves booking count aggregated by day for the last N days
func (r *adminRepository) GetBookingsByDay(ctx context.Context, days int) ([]DailyBookings, error) {
	query := `
		SELECT
			TO_CHAR(DATE_TRUNC('day', created_at), 'Mon DD') as day,
			COUNT(*) as count
		FROM bookings
		WHERE created_at >= NOW() - INTERVAL '1 day' * $1
		GROUP BY DATE_TRUNC('day', created_at)
		ORDER BY DATE_TRUNC('day', created_at)
	`

	rows, err := r.db.QueryContext(ctx, query, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []DailyBookings
	for rows.Next() {
		var item DailyBookings
		if err := rows.Scan(&item.Day, &item.Count); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	// If no data, fill with zeros for the last N days
	if len(result) == 0 {
		now := time.Now()
		for i := days - 1; i >= 0; i-- {
			day := now.AddDate(0, 0, -i)
			result = append(result, DailyBookings{
				Day:   day.Format("Mon 02"),
				Count: 0,
			})
		}
	}

	return result, rows.Err()
}
