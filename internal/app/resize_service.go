package app

import (
	"bytes"
	"context"
	"fmt"
	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"image"
	"image-resizing-service/internal/domain"
	"image-resizing-service/internal/ports"
	"image-resizing-service/pkg/utils"
)

type ResizeService struct {
	db                  *gorm.DB
	minio               *utils.MinioClient
	thumbnailRepository ports.ThumbnailRepository
	imageRepository     ports.ImageRepository
}

func NewResizeService(
	db *gorm.DB,
	minio *utils.MinioClient,
	thumbnailRepo ports.ThumbnailRepository,
	imageRepo ports.ImageRepository,
) ports.ResizeUseCase {
	return &ResizeService{
		db:                  db,
		minio:               minio,
		thumbnailRepository: thumbnailRepo,
		imageRepository:     imageRepo,
	}
}

func (s *ResizeService) ResizeThumbnails(ctx context.Context, imageID uuid.UUID, originalKey string) error {
	originalBytes, err := s.minio.GetFileAsBytes(ctx, originalKey)
	if err != nil {
		return fmt.Errorf("failed to download original: %w", err)
	}

	img, _, err := image.Decode(bytes.NewReader(originalBytes))
	if err != nil {
		return fmt.Errorf("failed to decode original image: %w", err)
	}

	for _, size := range domain.ThumbnailSizes {
		if err := s.generateAndSaveThumbnail(ctx, imageID, img, size); err != nil {
			return fmt.Errorf("failed to process thumbnail %s: %w", size.Label, err)
		}
	}

	imageEntity, err := s.imageRepository.FindByID(imageID.String())
	if err != nil {
		return fmt.Errorf("failed to find image entity: %w", err)
	}

	imageEntity.Status = domain.StatusReady

	if err := s.imageRepository.Update(imageEntity); err != nil {
		return fmt.Errorf("failed to update image status: %w", err)
	}

	return nil
}

func (s *ResizeService) generateAndSaveThumbnail(ctx context.Context, imageID uuid.UUID, img image.Image, size domain.ThumbnailSize) error {
	thumb := imaging.Resize(img, size.Width, size.Height, imaging.Lanczos)
	
	var buf bytes.Buffer
	options := &webp.Options{Quality: 80}
	if err := webp.Encode(&buf, thumb, options); err != nil {
		return fmt.Errorf("failed to encode thumbnail to webp: %w", err)
	}

	thumbKey := fmt.Sprintf("uploads/thumbnails/%s_%s.webp", imageID, size.Label)

	if err := s.minio.UploadBytes(ctx, thumbKey, buf.Bytes(), "image/webp"); err != nil {
		return fmt.Errorf("failed to upload thumbnail to storage: %w", err)
	}

	thumbnail := &domain.Thumbnail{
		ID:      uuid.New().String(),
		ImageID: imageID,
		Size:    size.Label,
		Key:     thumbKey,
		Type:    string(size.Type),
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		thumbRepo := s.thumbnailRepository.WithTx(tx)
		return thumbRepo.Save(thumbnail)
	})
	if err != nil {
		return fmt.Errorf("failed to save thumbnail record: %w", err)
	}

	return nil
}
