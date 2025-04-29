package utils

import (
	"bytes"
	"fmt"
	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func ConvertToWebp(filePath string) (string, error) {
	originalFile, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open original file: %w", err)
	}
	defer originalFile.Close()

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read original file: %w", err)
	}

	return ConvertBytesToWebp(fileBytes)
}

func ConvertBytesToWebp(file []byte) (string, error) {
	img, format, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	if format == "jpeg" {
		exifData, err := exif.Decode(bytes.NewReader(file))
		if err == nil {
			orientationTag, err := exifData.Get(exif.Orientation)
			if err == nil {
				orientation, _ := orientationTag.Int(0)
				img = autoRotate(img, orientation)
			}
		}
	}

	webpTempFile, err := os.CreateTemp("", "converted-*.webp")
	if err != nil {
		return "", fmt.Errorf("failed to create temp webp file: %w", err)
	}
	defer webpTempFile.Close()

	options := &webp.Options{
		Lossless: false,
		Quality:  80,
	}

	if err := webp.Encode(webpTempFile, img, options); err != nil {
		return "", fmt.Errorf("failed to encode image to webp: %w", err)
	}

	return webpTempFile.Name(), nil
}

func autoRotate(img image.Image, orientation int) image.Image {
	switch orientation {
	case 3:
		return imaging.Rotate180(img)
	case 6:
		return imaging.Rotate270(img)
	case 8:
		return imaging.Rotate90(img)
	default:
		return img
	}
}

func IsValidImage(data []byte) (bool, string) {
	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return false, ""
	}
	return true, format
}
