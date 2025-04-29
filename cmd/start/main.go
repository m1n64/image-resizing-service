package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"image-resizing-service/internal/delivery/grpc/handlers"
	images "image-resizing-service/internal/delivery/grpc/pb"
	"image-resizing-service/internal/delivery/rest"
	"image-resizing-service/pkg/di"
	"image-resizing-service/pkg/utils"
	"log"
	"net"
	"os"
)

var dependencies *di.Dependencies

func init() {
	dependencies = di.InitDependencies()
}

func main() {
	utils.InitMigrations(dependencies.DB)

	go func() {
		fmt.Println("Server started on port 8000")

		r := gin.Default()
		rest.InitRoutes(r, dependencies.ImageUsecase, dependencies.RestImageAssembler)

		r.Run(":8000")
	}()

	go func() {
		fmt.Println("GRPC server started on port 50051")

		listener, err := net.Listen("tcp", "0.0.0.0:50051")
		if err != nil {
			log.Fatal("failed to listen", zap.Error(err))
		}

		grpcToken := os.Getenv("GRPC_TOKEN")

		var grpcServer *grpc.Server
		if grpcToken != "" {
			grpcServer = grpc.NewServer(
				grpc.UnaryInterceptor(tokenAuthInterceptor(grpcToken)),
			)
		} else {
			grpcServer = grpc.NewServer()
		}

		imagesHandler := handlers.NewImageGRPCHandler(dependencies.ImageUsecase, dependencies.GRPCImageAssembler)
		images.RegisterImageServiceServer(grpcServer, imagesHandler)

		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("failed to serve", zap.Error(err))
		}
	}()

	select {}
}

func tokenAuthInterceptor(token string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeader, exists := md["authorization"]
		if !exists || len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		if authHeader[0] != token {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		return handler(ctx, req)
	}
}
