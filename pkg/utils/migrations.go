package utils

import (
	"gorm.io/gorm"
	"image-resizing-service/internal/domain"
)

func InitMigrations(db *gorm.DB) {
	db.AutoMigrate(&domain.Image{}, &domain.Thumbnail{})
}
