package handlers

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"image-resizing-service/internal/assembler"
	images "image-resizing-service/internal/delivery/grpc/pb"
	"image-resizing-service/internal/ports"
	"image-resizing-service/pkg/utils"
)

type ImageGRPCHandler struct {
	images.UnimplementedImageServiceServer
	useCase       ports.ImageUseCase
	grpcAssembler *assembler.GRPCImageAssembler
}

func NewImageGRPCHandler(useCase ports.ImageUseCase, grpcAssembler *assembler.GRPCImageAssembler) *ImageGRPCHandler {
	return &ImageGRPCHandler{
		useCase:       useCase,
		grpcAssembler: grpcAssembler,
	}
}

func (h *ImageGRPCHandler) UploadImage(ctx context.Context, req *images.UploadImageRequest) (*images.ImageResponse, error) {
	if req.Data == nil || len(req.Data) == 0 {
		return nil, status.Error(codes.InvalidArgument, "image data is required")
	}

	tempFilePath, contentType, err := utils.SaveTempFile(bytes.NewReader(req.Data))
	if err != nil {
		return nil, err
	}
	defer utils.RemoveFile(tempFilePath)

	result, err := h.useCase.UploadOriginal(ctx, tempFilePath, contentType)
	if err != nil {
		return nil, err
	}

	return h.grpcAssembler.BuildImage(ctx, result), nil
}

func (h *ImageGRPCHandler) GetImage(ctx context.Context, req *images.GetImageRequest) (*images.ImageResponse, error) {
	if uuid.Validate(req.Id) != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	imageData, err := h.useCase.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return h.grpcAssembler.BuildImageWithThumbnails(ctx, imageData), nil
}
