package di

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"image-resizing-service/internal/app"
	"image-resizing-service/internal/assembler"
	"image-resizing-service/internal/infrastructure/db"
	"image-resizing-service/internal/ports"
	"image-resizing-service/pkg/utils"
	"os"
)

type Dependencies struct {
	Logger      *zap.Logger
	Redis       *redis.Client
	DB          *gorm.DB
	Validator   *validator.Validate
	MinioClient *utils.MinioClient
	// Repositories
	ImageRepo     ports.ImageRepository
	ThumbnailRepo ports.ThumbnailRepository
	// Usecases
	ImageUsecase  ports.ImageUseCase
	ResizeUsecase ports.ResizeUseCase

	// Builders
	RestImageAssembler *assembler.RestImageAssembler
	GRPCImageAssembler *assembler.GRPCImageAssembler
}

func InitDependencies() *Dependencies {
	// Infrastructure
	logger := utils.InitLogs()
	utils.LoadEnv()
	redisConn := utils.CreateRedisConn(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	dbConn := utils.InitDBConnection(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	utils.InitMigrations(dbConn)

	validate := utils.InitValidator()

	minioClient := utils.NewMinioClient(os.Getenv("MINIO_ENDPOINT"), os.Getenv("MINIO_ROOT_USER"), os.Getenv("MINIO_ROOT_PASSWORD"), utils.BucketName, os.Getenv("MINIO_SECURE") == "true")

	// Repositories
	imageRepo := db.NewImageRepository(dbConn)
	thumbnailRepo := db.NewThumbnailRepository(dbConn)

	// Usecases
	resizeUsecase := app.NewResizeService(dbConn, minioClient, thumbnailRepo, imageRepo)
	imageUsecase := app.NewImageService(dbConn, imageRepo, resizeUsecase, minioClient, logger)

	// Assemblers
	restImageAssembler := assembler.NewRestImageAssembler(minioClient)
	grpcImageAssembler := assembler.NewGRPCImageAssembler(minioClient)

	return &Dependencies{
		Logger:             logger,
		Redis:              redisConn,
		DB:                 dbConn,
		Validator:          validate,
		MinioClient:        minioClient,
		ImageRepo:          imageRepo,
		ThumbnailRepo:      thumbnailRepo,
		ImageUsecase:       imageUsecase,
		ResizeUsecase:      resizeUsecase,
		RestImageAssembler: restImageAssembler,
		GRPCImageAssembler: grpcImageAssembler,
	}
}
