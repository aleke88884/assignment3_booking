package models

import "time"

// ResourcePhoto представляет фотографию ресурса
type ResourcePhoto struct {
	ID           int64     `json:"id"`
	ResourceID   int64     `json:"resource_id"`
	URL          string    `json:"url"`
	StorageKey   string    `json:"storage_key,omitempty"`
	FileName     string    `json:"file_name"`
	FileSize     int64     `json:"file_size"`
	MimeType     string    `json:"mime_type"`
	Width        int       `json:"width,omitempty"`
	Height       int       `json:"height,omitempty"`
	IsPrimary    bool      `json:"is_primary"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// PhotoUploadRequest для загрузки фото
type PhotoUploadRequest struct {
	ResourceID   int64  `json:"resource_id"`
	IsPrimary    bool   `json:"is_primary"`
	DisplayOrder int    `json:"display_order"`
}
