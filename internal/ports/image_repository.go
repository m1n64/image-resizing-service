package ports

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"image-resizing-service/internal/domain"
	"time"
)

type ImageRepository interface {
	WithTx(tx *gorm.DB) ImageRepository
	Save(image *domain.Image) error
	Update(image *domain.Image) error
	FindByID(id string) (*domain.Image, error)
	WaitForImage(ctx context.Context, id uuid.UUID, retries int, delay time.Duration) error
}
