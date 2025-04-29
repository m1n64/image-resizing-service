package assembler

import (
	"context"
	"image-resizing-service/internal/domain"
	"image-resizing-service/internal/dto"
	"image-resizing-service/internal/ports"
	"image-resizing-service/pkg/utils"
)

type RestImageAssembler struct {
	minio *utils.MinioClient
}

func NewRestImageAssembler(minio *utils.MinioClient) *RestImageAssembler {
	return &RestImageAssembler{minio: minio}
}

func (a *RestImageAssembler) BuildImage(image *ports.UploadResult) *dto.ImageWithThumbnails {
	ctx := context.Background()

	originalUrl, err := a.minio.GetFileURL(ctx, image.OriginalKey)
	if err != nil {
		return nil
	}

	return &dto.ImageWithThumbnails{
		ID:          image.ID,
		OriginalUrl: originalUrl,
		Status:      image.Status,
	}
}

func (a *RestImageAssembler) BuildImageWithThumbnails(image *domain.Image) *dto.ImageWithThumbnails {
	ctx := context.Background()

	thumbnails := make([]dto.ThumbnailShort, 0, len(image.Thumbnails))
	for _, thumb := range image.Thumbnails {
		url, err := a.minio.GetFileURL(ctx, thumb.Key)
		if err != nil {
			return nil
		}

		thumbnails = append(thumbnails, dto.ThumbnailShort{
			Size: thumb.Size,
			Url:  url,
			Type: thumb.Type,
		})
	}

	originalUrl, err := a.minio.GetFileURL(ctx, image.OriginalKey)
	if err != nil {
		return nil
	}

	compressedUrl, err := a.minio.GetFileURL(ctx, image.CompressedKey)
	if err != nil {
		return nil
	}

	return &dto.ImageWithThumbnails{
		ID:            image.ID,
		OriginalUrl:   originalUrl,
		CompressedUrl: compressedUrl,
		Status:        string(image.Status),
		ErrorMessage:  image.ErrorMessage,
		Thumbnails:    thumbnails,
	}
}
