package ports

import (
	"context"
	"github.com/google/uuid"
)

type ResizeUseCase interface {
	ResizeThumbnails(ctx context.Context, imageID uuid.UUID, originalKey string) error
}
