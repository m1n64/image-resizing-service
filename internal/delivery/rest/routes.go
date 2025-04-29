package rest

import (
	"github.com/gin-gonic/gin"
	"image-resizing-service/internal/assembler"
	"image-resizing-service/internal/delivery/rest/handlers"
	"image-resizing-service/internal/ports"
	"net/http"
	"os"
)

func InitRoutes(router *gin.Engine, imageUseCase ports.ImageUseCase, restImageAssembler *assembler.RestImageAssembler) {
	imageHandler := handlers.NewImageHandler(imageUseCase, restImageAssembler)

	router.POST("/image/upload", UploadAuthMiddleware(), imageHandler.UploadImage)
	router.POST("/image/upload-binary", UploadAuthMiddleware(), imageHandler.UploadImageBinary)
	router.GET("/image/:id", imageHandler.GetImage)
}

func UploadAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uploadToken := os.Getenv("REST_UPLOAD_TOKEN")
		if uploadToken == "" {
			c.Next()
			return
		}

		if c.GetHeader("X-API-Key") != uploadToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		c.Next()
	}
}
