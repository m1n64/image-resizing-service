package ports

import (
	"context"
	"image-resizing-service/internal/domain"
)

type UploadResult struct {
	ID          string
	OriginalKey string
	Status      string
}

type ImageUseCase interface {
	UploadOriginal(ctx context.Context, filePath string, contentType string) (*UploadResult, error)
	FindByID(ctx context.Context, id string) (*domain.Image, error)
}
