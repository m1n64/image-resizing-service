package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Thumbnail struct {
	ID      string    `gorm:"type:uuid;primaryKey"`
	ImageID uuid.UUID `gorm:"type:uuid;not null;index:idx_thumbnails_image_id_size,unique"`
	Size    string    `gorm:"type:varchar(50);not null;index:idx_thumbnails_image_id_size,unique"`
	Key     string    `gorm:"not null"`
	Type    string    `gorm:"type:varchar(20);not null"`
	gorm.Model
}

func (thumb *Thumbnail) BeforeCreate(tx *gorm.DB) (err error) {
	thumb.ID = uuid.New().String()

	return
}
