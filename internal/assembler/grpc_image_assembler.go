package assembler

import (
	"context"
	images "image-resizing-service/internal/delivery/grpc/pb"
	"image-resizing-service/internal/domain"
	"image-resizing-service/internal/ports"
	"image-resizing-service/pkg/utils"
)

type GRPCImageAssembler struct {
	minio *utils.MinioClient
}

func NewGRPCImageAssembler(minio *utils.MinioClient) *GRPCImageAssembler {
	return &GRPCImageAssembler{minio: minio}
}

func (g *GRPCImageAssembler) BuildImage(ctx context.Context, image *ports.UploadResult) *images.ImageResponse {
	originalUrl, err := g.minio.GetFileURL(ctx, image.OriginalKey)
	if err != nil {
		return nil
	}

	return &images.ImageResponse{
		Id:          image.ID,
		OriginalUrl: originalUrl,
		Status:      image.Status,
	}
}

func (g *GRPCImageAssembler) BuildImageWithThumbnails(ctx context.Context, image *domain.Image) *images.ImageResponse {
	originalUrl, err := g.minio.GetFileURL(ctx, image.OriginalKey)
	if err != nil {
		return nil
	}

	compressedUrl, err := g.minio.GetFileURL(ctx, image.CompressedKey)
	if err != nil {
		return nil
	}

	var thumbnails []*images.ThumbnailShort
	for _, thumb := range image.Thumbnails {
		url, err := g.minio.GetFileURL(ctx, thumb.Key)
		if err != nil {
			return nil
		}

		thumbnails = append(thumbnails, &images.ThumbnailShort{
			Size: thumb.Size,
			Url:  url,
			Type: thumb.Type,
		})
	}

	return &images.ImageResponse{
		Id:            image.ID,
		OriginalUrl:   originalUrl,
		CompressedUrl: &compressedUrl,
		Status:        string(image.Status),
		Thumbnails:    thumbnails,
	}
}
