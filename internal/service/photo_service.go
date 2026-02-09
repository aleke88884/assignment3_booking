package service

import (
	"context"
	"fmt"
	"io"
	"strings"

	"smartbooking/internal/models"
	"smartbooking/internal/repository"
	"smartbooking/internal/storage"
)

// PhotoService интерфейс для работы с фотографиями
type PhotoService interface {
	UploadPhoto(ctx context.Context, resourceID int64, file io.Reader, fileName string, isPrimary bool) (*models.ResourcePhoto, error)
	GetResourcePhotos(ctx context.Context, resourceID int64) ([]*models.ResourcePhoto, error)
	DeletePhoto(ctx context.Context, id int64) error
	SetPrimaryPhoto(ctx context.Context, id int64, resourceID int64) error
}

type photoService struct {
	photoRepo repository.PhotoRepository
	storage   storage.StorageService
}

func NewPhotoService(photoRepo repository.PhotoRepository, storage storage.StorageService) PhotoService {
	return &photoService{
		photoRepo: photoRepo,
		storage:   storage,
	}
}

func (s *photoService) UploadPhoto(ctx context.Context, resourceID int64, file io.Reader, fileName string, isPrimary bool) (*models.ResourcePhoto, error) {
	// Определяем MIME type по расширению
	contentType := getContentType(fileName)
	if !isImageMimeType(contentType) {
		return nil, fmt.Errorf("неподдерживаемый тип файла: %s", contentType)
	}

	// Загружаем в storage
	result, err := s.storage.UploadFile(ctx, file, fileName, contentType)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки файла: %w", err)
	}

	// Создаем запись в БД
	photo := &models.ResourcePhoto{
		ResourceID: resourceID,
		URL:        result.URL,
		StorageKey: result.StorageKey,
		FileName:   result.FileName,
		FileSize:   result.FileSize,
		MimeType:   contentType,
		IsPrimary:  isPrimary,
	}

	err = s.photoRepo.Create(ctx, photo)
	if err != nil {
		// Если не удалось создать запись в БД, удаляем файл из storage
		_ = s.storage.DeleteFile(ctx, result.StorageKey)
		return nil, fmt.Errorf("ошибка сохранения в БД: %w", err)
	}

	// Если это главное фото, обновляем остальные
	if isPrimary {
		_ = s.photoRepo.SetPrimary(ctx, photo.ID, resourceID)
	}

	return photo, nil
}

func (s *photoService) GetResourcePhotos(ctx context.Context, resourceID int64) ([]*models.ResourcePhoto, error) {
	return s.photoRepo.GetByResourceID(ctx, resourceID)
}

func (s *photoService) DeletePhoto(ctx context.Context, id int64) error {
	// Получаем фото чтобы удалить из storage
	photo, err := s.photoRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Удаляем из storage
	err = s.storage.DeleteFile(ctx, photo.StorageKey)
	if err != nil {
		return fmt.Errorf("ошибка удаления из storage: %w", err)
	}

	// Удаляем из БД
	err = s.photoRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления из БД: %w", err)
	}

	return nil
}

func (s *photoService) SetPrimaryPhoto(ctx context.Context, id int64, resourceID int64) error {
	return s.photoRepo.SetPrimary(ctx, id, resourceID)
}

// Вспомогательные функции

func getContentType(fileName string) string {
	ext := strings.ToLower(fileName[strings.LastIndex(fileName, ".")+1:])
	switch ext {
	case "jpg", "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "gif":
		return "image/gif"
	case "webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

func isImageMimeType(mimeType string) bool {
	return strings.HasPrefix(mimeType, "image/")
}
