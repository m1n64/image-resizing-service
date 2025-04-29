package utils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const BucketName = "images"

type MinioClient struct {
	client *minio.Client
	bucket string
}

type customHostTransport struct {
	Transport  http.RoundTripper
	PublicHost string
}

func NewMinioClient(endpoint, accessKey, secretKey, bucketName string, useSSL bool) *MinioClient {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	ctx := context.Background()
	err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists", bucketName)
		} else {
			log.Fatalf("Failed to create bucket: %v", err)
		}
	}

	return &MinioClient{client: client, bucket: bucketName}
}

func (m *MinioClient) UploadFile(ctx context.Context, objectName, filePath, contentType string) error {
	_, err := m.client.FPutObject(ctx, m.bucket, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	return err
}

func (m *MinioClient) GetFile(ctx context.Context, objectName, destinationPath string) error {
	object, err := m.client.GetObject(ctx, m.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer object.Close()

	file, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, object)
	if err != nil {
		return err
	}

	return nil
}

func (m *MinioClient) GetFileAsStream(ctx context.Context, objectName string) (io.Reader, error) {
	object, err := m.client.GetObject(ctx, m.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (m *MinioClient) GetFileAsBytes(ctx context.Context, objectName string) ([]byte, error) {
	object, err := m.client.GetObject(ctx, m.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %v", err)
	}

	if object == nil {
		return nil, fmt.Errorf("received nil object for %s", objectName)
	}

	defer object.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, object)
	if err != nil {
		return nil, fmt.Errorf("failed to copy object data: %v", err)
	}

	return buf.Bytes(), nil
}

func (m *MinioClient) UploadBytes(ctx context.Context, key string, data []byte, contentType string) error {
	reader := bytes.NewReader(data)

	_, err := m.client.PutObject(ctx, m.bucket, key, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (m *MinioClient) GetFileURL(ctx context.Context, objectName string) (string, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", objectName))

	presignedURL, err := m.client.PresignedGetObject(
		ctx,
		m.bucket,
		objectName,
		time.Hour*24,
		reqParams,
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	return presignedURL.String(), nil
}
