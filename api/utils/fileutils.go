package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var supportedImageTypes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
	"image/jpg":  true,
}

const UploadDir = "./uploads"

const TempDir = "./temp"

const MaxFileSize = 100 * 1024 * 1024

func EnsureDirectoryExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func GenerateUniqueFilename(originalName string) string {
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	hash := sha256.Sum256([]byte(originalName + timestamp))
	hashString := hex.EncodeToString(hash[:])[:12]

	ext := filepath.Ext(originalName)
	return hashString + ext
}

func SaveUploadedFile(file *multipart.FileHeader) (string, error) {
	if err := EnsureDirectoryExists(UploadDir); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	if file.Size > MaxFileSize {
		return "", fmt.Errorf("file too large (max %d bytes)", MaxFileSize)
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	src.Seek(0, io.SeekStart)

	contentType := http.DetectContentType(buffer)
	if !supportedImageTypes[contentType] && strings.HasPrefix(file.Filename, "image/") {
		return "", fmt.Errorf("unsupported image format: %s", contentType)
	}

	uniqueFilename := GenerateUniqueFilename(file.Filename)
	filePath := filepath.Join(UploadDir, uniqueFilename)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return filePath, nil
}

func CleanupTempFiles(maxAge time.Duration) error {
	if err := EnsureDirectoryExists(TempDir); err != nil {
		return err
	}

	now := time.Now()

	return filepath.Walk(TempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == TempDir {
			return nil
		}

		if now.Sub(info.ModTime()) > maxAge {
			if err := os.Remove(path); err != nil {
				return err
			}
		}

		return nil
	})
}

func SaveOutputFile(data []byte, extension string) (string, error) {
	if err := EnsureDirectoryExists(TempDir); err != nil {
		return "", err
	}

	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	hash := sha256.Sum256([]byte(timestamp))
	filename := hex.EncodeToString(hash[:])[:12] + extension
	filePath := filepath.Join(TempDir, filename)

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", err
	}

	return filePath, nil
}
