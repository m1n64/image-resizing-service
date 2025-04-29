package utils

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

func SaveTempFile(file io.Reader) (string, string, error) {
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", "", err
	}
	contentType := http.DetectContentType(buf[:n])

	reader := io.MultiReader(bytes.NewReader(buf[:n]), file)

	tempFile, err := os.CreateTemp("tmp", "upload-*.tmp")
	if err != nil {
		return "", "", err
	}
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, reader); err != nil {
		return "", "", err
	}

	return tempFile.Name(), contentType, nil
}

func RemoveFile(path string) {
	_ = os.Remove(path)
}
