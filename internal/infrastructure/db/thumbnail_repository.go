package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"image-resizing-service/internal/domain"
	"image-resizing-service/internal/ports"
)

type ThumbnailRepositoryImpl struct {
	db *gorm.DB
}

func NewThumbnailRepository(db *gorm.DB) ports.ThumbnailRepository {
	return &ThumbnailRepositoryImpl{db: db}
}

func (r *ThumbnailRepositoryImpl) WithTx(tx *gorm.DB) ports.ThumbnailRepository {
	return &ThumbnailRepositoryImpl{db: tx}
}

func (r *ThumbnailRepositoryImpl) Save(thumbnail *domain.Thumbnail) error {
	return r.db.Create(&thumbnail).Error
}

func (r *ThumbnailRepositoryImpl) FindByImageID(imageID uuid.UUID) ([]domain.Thumbnail, error) {
	var thumbnails []domain.Thumbnail
	err := r.db.Where("image_id = ?", imageID).Find(&thumbnails).Error
	if err != nil {
		return nil, err
	}

	return thumbnails, nil
}

func (r *ThumbnailRepositoryImpl) Exists(imageID uuid.UUID, sizeLabel string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Thumbnail{}).
		Where("image_id = ? AND size = ?", imageID, sizeLabel).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
