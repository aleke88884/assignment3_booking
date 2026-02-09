package models

import "time"

// BookingStatus represents the status of a booking
type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
)

// Booking represents a reservation of a resource by a user
type Booking struct {
	ID         int64         `json:"id"`
	UserID     int64         `json:"user_id"`
	ResourceID int64         `json:"resource_id"`
	StartTime  time.Time     `json:"start_time"`
	EndTime    time.Time     `json:"end_time"`
	Status     BookingStatus `json:"status"`
	TotalPrice float64       `json:"total_price,omitempty"`
	Notes      string        `json:"notes,omitempty"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`

	// For JOIN queries
	UserName     string `json:"user_name,omitempty"`
	UserEmail    string `json:"user_email,omitempty"`
	ResourceName string `json:"resource_name,omitempty"`
}
