package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"image-resizing-service/internal/domain"
	"image-resizing-service/internal/ports"
	"time"
)

type ImageRepositoryImpl struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) ports.ImageRepository {
	return &ImageRepositoryImpl{db: db}
}

func (r *ImageRepositoryImpl) WithTx(tx *gorm.DB) ports.ImageRepository {
	return &ImageRepositoryImpl{db: tx}
}

func (r *ImageRepositoryImpl) Save(image *domain.Image) error {
	return r.db.Create(image).Error
}

func (r *ImageRepositoryImpl) Update(image *domain.Image) error {
	return r.db.Model(&domain.Image{}).
		Where("id = ?", image.ID).
		Updates(map[string]interface{}{
			"compressed_key": image.CompressedKey,
			"status":         image.Status,
			"error_message":  image.ErrorMessage,
		}).Error
}

func (r *ImageRepositoryImpl) FindByID(id string) (*domain.Image, error) {
	var img domain.Image
	err := r.db.Preload("Thumbnails").First(&img, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func (r *ImageRepositoryImpl) WaitForImage(ctx context.Context, id uuid.UUID, retries int, delay time.Duration) error {
	for attempt := 0; attempt < retries; attempt++ {
		var count int64
		err := r.db.WithContext(ctx).Model(&domain.Image{}).Where("id = ?", id).Count(&count).Error
		if err != nil {
			return fmt.Errorf("failed to check image existence: %w", err)
		}

		if count > 0 {
			return nil
		}

		time.Sleep(delay)
	}

	return fmt.Errorf("image %s not found after %d retries", id.String(), retries)
}
