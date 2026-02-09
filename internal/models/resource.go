package models

import (
	"database/sql"
	"time"
)

// Resource представляет бронируемый ресурс
type Resource struct {
	ID           int64            `json:"id"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Capacity     int              `json:"capacity"`
	CategoryID   *int64           `json:"category_id,omitempty"`
	Address      string           `json:"address,omitempty"`
	City         string           `json:"city,omitempty"`
	Latitude     *float64         `json:"latitude,omitempty"`
	Longitude    *float64         `json:"longitude,omitempty"`
	Amenities    []string         `json:"amenities,omitempty"`
	Rules        string           `json:"rules,omitempty"`
	PricePerHour *float64         `json:"price_per_hour,omitempty"`
	IsActive     bool             `json:"is_active"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`

	// Для JOIN запросов
	CategoryName string            `json:"category_name,omitempty"`
	Photos       []ResourcePhoto   `json:"photos,omitempty"`
	Rating       float64           `json:"rating,omitempty"`
	ReviewsCount int               `json:"reviews_count,omitempty"`
}

// ResourceCreateRequest для создания ресурса
type ResourceCreateRequest struct {
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Capacity     int       `json:"capacity"`
	CategoryID   *int64    `json:"category_id"`
	Address      string    `json:"address"`
	City         string    `json:"city"`
	Latitude     *float64  `json:"latitude"`
	Longitude    *float64  `json:"longitude"`
	Amenities    []string  `json:"amenities"`
	Rules        string    `json:"rules"`
	PricePerHour *float64  `json:"price_per_hour"`
}

// ResourceUpdateRequest для обновления ресурса
type ResourceUpdateRequest struct {
	Name         *string   `json:"name,omitempty"`
	Description  *string   `json:"description,omitempty"`
	Capacity     *int      `json:"capacity,omitempty"`
	CategoryID   *int64    `json:"category_id,omitempty"`
	Address      *string   `json:"address,omitempty"`
	City         *string   `json:"city,omitempty"`
	Latitude     *float64  `json:"latitude,omitempty"`
	Longitude    *float64  `json:"longitude,omitempty"`
	Amenities    []string  `json:"amenities,omitempty"`
	Rules        *string   `json:"rules,omitempty"`
	PricePerHour *float64  `json:"price_per_hour,omitempty"`
	IsActive     *bool     `json:"is_active,omitempty"`
}

// ResourceFilterParams для фильтрации ресурсов
type ResourceFilterParams struct {
	CategoryID *int64   `json:"category_id"`
	City       string   `json:"city"`
	MinPrice   *float64 `json:"min_price"`
	MaxPrice   *float64 `json:"max_price"`
	IsActive   *bool    `json:"is_active"`
	Limit      int      `json:"limit"`
	Offset     int      `json:"offset"`
}

// ScanAmenities помощник для сканирования amenities из БД
func ScanAmenities(src interface{}) ([]string, error) {
	if src == nil {
		return []string{}, nil
	}

	switch v := src.(type) {
	case []uint8:
		return parsePostgresArray(string(v)), nil
	case string:
		return parsePostgresArray(v), nil
	default:
		return []string{}, nil
	}
}

func parsePostgresArray(s string) []string {
	if s == "" || s == "{}" {
		return []string{}
	}
	// Убираем { и }
	s = s[1 : len(s)-1]

	result := []string{}
	current := ""
	inQuotes := false

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '"' {
			inQuotes = !inQuotes
		} else if c == ',' && !inQuotes {
			result = append(result, current)
			current = ""
		} else {
			current += string(c)
		}
	}

	if current != "" {
		result = append(result, current)
	}

	return result
}

// NullInt64ToPtr конвертирует sql.NullInt64 в *int64
func NullInt64ToPtr(n sql.NullInt64) *int64 {
	if n.Valid {
		return &n.Int64
	}
	return nil
}

// NullFloat64ToPtr конвертирует sql.NullFloat64 в *float64
func NullFloat64ToPtr(n sql.NullFloat64) *float64 {
	if n.Valid {
		return &n.Float64
	}
	return nil
}
