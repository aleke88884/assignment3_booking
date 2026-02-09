package models

import "time"

// Review представляет отзыв о ресурсе
type Review struct {
	ID         int64     `json:"id"`
	ResourceID int64     `json:"resource_id"`
	UserID     int64     `json:"user_id"`
	BookingID  *int64    `json:"booking_id,omitempty"`
	Rating     int       `json:"rating"`
	Comment    string    `json:"comment,omitempty"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Для JOIN запросов
	UserName string `json:"user_name,omitempty"`
}

// ReviewCreateRequest для создания отзыва
type ReviewCreateRequest struct {
	ResourceID int64  `json:"resource_id"`
	BookingID  *int64 `json:"booking_id,omitempty"`
	Rating     int    `json:"rating"`
	Comment    string `json:"comment"`
}
