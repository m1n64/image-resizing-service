package domain

import (
	"gorm.io/gorm"
)

type ImageStatus string

const (
	StatusPending    ImageStatus = "pending"
	StatusProcessing ImageStatus = "processing"
	StatusReady      ImageStatus = "ready"
	StatusError      ImageStatus = "error"
)

type Image struct {
	ID            string      `gorm:"type:uuid;primaryKey"`
	OriginalKey   string      `gorm:"not null"`
	CompressedKey string      `gorm:""`
	Status        ImageStatus `gorm:"type:varchar(20);not null;default:'pending'"`
	ErrorMessage  *string     `gorm:""`
	Thumbnails    []Thumbnail `gorm:"foreignKey:ImageID;constraint:OnDelete:CASCADE;"`
	gorm.Model
}

/*func (img *Image) BeforeCreate(tx *gorm.DB) (err error) {
	img.ID = uuid.New().String()
	img.Status = StatusPending
	return
}*/
