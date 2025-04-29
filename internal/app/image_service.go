package app

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"image-resizing-service/internal/domain"
	"image-resizing-service/internal/ports"
	"image-resizing-service/pkg/utils"
	"os"
)

type ImageService struct {
	db              *gorm.DB
	imageRepository ports.ImageRepository
	resizeService   ports.ResizeUseCase
	minio           *utils.MinioClient
	logger          *zap.Logger
}

func NewImageService(
	db *gorm.DB,
	imageRepo ports.ImageRepository,
	resizeService ports.ResizeUseCase,
	minio *utils.MinioClient,
	logger *zap.Logger,
) ports.ImageUseCase {
	return &ImageService{
		db:              db,
		imageRepository: imageRepo,
		resizeService:   resizeService,
		minio:           minio,
		logger:          logger,
	}
}

func (s *ImageService) UploadOriginal(ctx context.Context, filePath string, contentType string) (*ports.UploadResult, error) {
	id := uuid.New()

	originalKey := fmt.Sprintf("uploads/originals/%s", id.String())

	if err := s.minio.UploadFile(ctx, originalKey, filePath, contentType); err != nil {
		return nil, fmt.Errorf("failed to upload original file: %w", err)
	}

	image := &domain.Image{
		ID:          id.String(),
		OriginalKey: originalKey,
		Status:      domain.StatusPending,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		imgRepo := s.imageRepository.WithTx(tx)
		return imgRepo.Save(image)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save image record: %w", err)
	}

	go s.compressAndDispatch(context.Background(), id, originalKey)

	return &ports.UploadResult{
		ID:          id.String(),
		OriginalKey: originalKey,
		Status:      string(domain.StatusPending),
	}, nil
}

func (s *ImageService) FindByID(ctx context.Context, id string) (*domain.Image, error) {
	return s.imageRepository.FindByID(id)
}

func (s *ImageService) compressAndDispatch(ctx context.Context, id uuid.UUID, originalKey string) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("panic in compressAndDispatch", zap.Any("recover", r))
		}
	}()

	originalFile, err := s.minio.GetFileAsBytes(ctx, originalKey)
	if err != nil {
		s.markAsError(ctx, id, fmt.Errorf("failed to get file as bytes: %w", err))
		return
	}

	webpPath, err := utils.ConvertBytesToWebp(originalFile)
	if err != nil {
		s.markAsError(ctx, id, fmt.Errorf("failed to convert to webp: %w", err))
		return
	}
	defer os.Remove(webpPath)

	compressedKey := fmt.Sprintf("uploads/compressed/%s.webp", id.String())

	if err := s.minio.UploadFile(ctx, compressedKey, webpPath, "image/webp"); err != nil {
		s.markAsError(ctx, id, fmt.Errorf("failed to upload compressed webp: %w", err))
		return
	}

	image, err := s.imageRepository.FindByID(id.String())
	if err != nil {
		s.logger.Error("failed to find image for update", zap.Error(err))
		return
	}

	image.CompressedKey = compressedKey
	image.Status = domain.StatusProcessing

	if err := s.imageRepository.Update(image); err != nil {
		s.logger.Error("failed to update image after compression", zap.Error(err))
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.logger.Error("panic in resize thumbnails", zap.Any("recover", r))
			}
		}()

		if err := s.resizeService.ResizeThumbnails(context.Background(), id, compressedKey); err != nil {
			s.logger.Error("failed to resize thumbnails", zap.Error(err))
			s.markAsError(context.Background(), id, err)
		}
	}()
}

func (s *ImageService) markAsError(ctx context.Context, id uuid.UUID, originalErr error) {
	image, err := s.imageRepository.FindByID(id.String())
	if err != nil {
		s.logger.Error("failed to find image for error update", zap.Error(err))
		return
	}

	errorMessage := originalErr.Error()
	image.Status = domain.StatusError
	image.ErrorMessage = &errorMessage

	if err := s.imageRepository.Update(image); err != nil {
		s.logger.Error("failed to update image status to error", zap.Error(err))
	}
}
