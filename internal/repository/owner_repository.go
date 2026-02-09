package repository

import (
	"context"
	"database/sql"

	"smartbooking/internal/models"
)

// OwnerRepository defines the interface for owner-specific data operations
type OwnerRepository interface {
	GetOwnerResources(ctx context.Context, ownerID int64) ([]*models.Resource, error)
	GetOwnerBookings(ctx context.Context, ownerID int64) ([]*models.Booking, error)
	GetOwnerStatistics(ctx context.Context, ownerID int64) (*OwnerStatistics, error)
}

// OwnerStatistics represents aggregated statistics for an owner
type OwnerStatistics struct {
	TotalResources   int     `json:"total_resources"`
	TotalBookings    int     `json:"total_bookings"`
	ActiveBookings   int     `json:"active_bookings"`
	CancelledBookings int    `json:"cancelled_bookings"`
	TotalRevenue     float64 `json:"total_revenue"`
	AverageRating    float64 `json:"average_rating"`
	TotalReviews     int     `json:"total_reviews"`
}

// ownerRepository implements OwnerRepository interface with PostgreSQL storage
type ownerRepository struct {
	db *sql.DB
}

// NewOwnerRepository creates a new instance of OwnerRepository
func NewOwnerRepository(db *sql.DB) OwnerRepository {
	return &ownerRepository{
		db: db,
	}
}

// GetOwnerResources retrieves all resources owned by a specific owner
func (r *ownerRepository) GetOwnerResources(ctx context.Context, ownerID int64) ([]*models.Resource, error) {
	query := `
		SELECT
			r.id, r.name, r.description, r.capacity, r.owner_id,
			r.category_id, r.address, r.city, r.latitude, r.longitude,
			r.amenities, r.rules, r.price_per_hour, r.is_active,
			r.created_at, r.updated_at,
			u.name as owner_name,
			c.name as category_name,
			COALESCE(AVG(rv.rating), 0) as rating,
			COUNT(DISTINCT rv.id) as reviews_count
		FROM resources r
		LEFT JOIN users u ON r.owner_id = u.id
		LEFT JOIN resource_categories c ON r.category_id = c.id
		LEFT JOIN reviews rv ON r.id = rv.resource_id
		WHERE r.owner_id = $1
		GROUP BY r.id, u.name, c.name
		ORDER BY r.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resources := make([]*models.Resource, 0)
	for rows.Next() {
		resource := &models.Resource{}
		var ownerID, categoryID sql.NullInt64
		var ownerName, categoryName, address, city, rules sql.NullString
		var latitude, longitude, pricePerHour sql.NullFloat64
		var amenities sql.NullString

		err := rows.Scan(
			&resource.ID,
			&resource.Name,
			&resource.Description,
			&resource.Capacity,
			&ownerID,
			&categoryID,
			&address,
			&city,
			&latitude,
			&longitude,
			&amenities,
			&rules,
			&pricePerHour,
			&resource.IsActive,
			&resource.CreatedAt,
			&resource.UpdatedAt,
			&ownerName,
			&categoryName,
			&resource.Rating,
			&resource.ReviewsCount,
		)
		if err != nil {
			return nil, err
		}

		resource.OwnerID = models.NullInt64ToPtr(ownerID)
		resource.CategoryID = models.NullInt64ToPtr(categoryID)
		resource.Latitude = models.NullFloat64ToPtr(latitude)
		resource.Longitude = models.NullFloat64ToPtr(longitude)
		resource.PricePerHour = models.NullFloat64ToPtr(pricePerHour)

		if ownerName.Valid {
			resource.OwnerName = ownerName.String
		}
		if categoryName.Valid {
			resource.CategoryName = categoryName.String
		}
		if address.Valid {
			resource.Address = address.String
		}
		if city.Valid {
			resource.City = city.String
		}
		if rules.Valid {
			resource.Rules = rules.String
		}
		if amenities.Valid {
			resource.Amenities, _ = models.ScanAmenities(amenities.String)
		}

		resources = append(resources, resource)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return resources, nil
}

// GetOwnerBookings retrieves all bookings for resources owned by a specific owner
func (r *ownerRepository) GetOwnerBookings(ctx context.Context, ownerID int64) ([]*models.Booking, error) {
	query := `
		SELECT
			b.id, b.user_id, b.resource_id, b.start_time, b.end_time,
			b.status, b.total_price, b.notes, b.created_at, b.updated_at,
			u.name as user_name, u.email as user_email,
			r.name as resource_name
		FROM bookings b
		INNER JOIN resources r ON b.resource_id = r.id
		INNER JOIN users u ON b.user_id = u.id
		WHERE r.owner_id = $1
		ORDER BY b.start_time DESC
	`

	rows, err := r.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := make([]*models.Booking, 0)
	for rows.Next() {
		booking := &models.Booking{}
		var totalPrice sql.NullFloat64
		var notes sql.NullString
		var userName, userEmail, resourceName string

		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.ResourceID,
			&booking.StartTime,
			&booking.EndTime,
			&booking.Status,
			&totalPrice,
			&notes,
			&booking.CreatedAt,
			&booking.UpdatedAt,
			&userName,
			&userEmail,
			&resourceName,
		)
		if err != nil {
			return nil, err
		}

		if totalPrice.Valid {
			booking.TotalPrice = totalPrice.Float64
		}
		if notes.Valid {
			booking.Notes = notes.String
		}

		// Add user and resource details (can extend models if needed)
		booking.UserName = userName
		booking.UserEmail = userEmail
		booking.ResourceName = resourceName

		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

// GetOwnerStatistics calculates and retrieves statistics for an owner
func (r *ownerRepository) GetOwnerStatistics(ctx context.Context, ownerID int64) (*OwnerStatistics, error) {
	stats := &OwnerStatistics{}

	// Get resource count
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM resources WHERE owner_id = $1
	`, ownerID).Scan(&stats.TotalResources)
	if err != nil {
		return nil, err
	}

	// Get booking statistics
	query := `
		SELECT
			COUNT(*) as total_bookings,
			SUM(CASE WHEN b.status IN ('pending', 'confirmed') THEN 1 ELSE 0 END) as active_bookings,
			SUM(CASE WHEN b.status = 'cancelled' THEN 1 ELSE 0 END) as cancelled_bookings,
			COALESCE(SUM(CASE WHEN b.status IN ('pending', 'confirmed') THEN b.total_price ELSE 0 END), 0) as total_revenue
		FROM bookings b
		INNER JOIN resources r ON b.resource_id = r.id
		WHERE r.owner_id = $1
	`

	err = r.db.QueryRowContext(ctx, query, ownerID).Scan(
		&stats.TotalBookings,
		&stats.ActiveBookings,
		&stats.CancelledBookings,
		&stats.TotalRevenue,
	)
	if err != nil {
		return nil, err
	}

	// Get rating statistics
	ratingQuery := `
		SELECT
			COALESCE(AVG(rv.rating), 0) as average_rating,
			COUNT(rv.id) as total_reviews
		FROM reviews rv
		INNER JOIN resources r ON rv.resource_id = r.id
		WHERE r.owner_id = $1
	`

	err = r.db.QueryRowContext(ctx, ratingQuery, ownerID).Scan(
		&stats.AverageRating,
		&stats.TotalReviews,
	)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
