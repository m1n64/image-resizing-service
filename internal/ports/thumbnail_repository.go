package ports

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"image-resizing-service/internal/domain"
)

type ThumbnailRepository interface {
	WithTx(tx *gorm.DB) ThumbnailRepository
	Save(thumbnail *domain.Thumbnail) error
	FindByImageID(imageID uuid.UUID) ([]domain.Thumbnail, error)
	Exists(imageID uuid.UUID, sizeLabel string) (bool, error)
}
