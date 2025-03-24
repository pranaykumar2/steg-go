package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http" // Added the missing import
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Supported image MIME types
var supportedImageTypes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
	"image/jpg":  true,
}

// UploadDir is the directory where uploaded files will be stored
const UploadDir = "./uploads"

// TempDir is the directory where temporary files will be stored
const TempDir = "./temp"

// MaxFileSize is the maximum allowed size for uploaded files (10MB)
const MaxFileSize = 100 * 1024 * 1024

// EnsureDirectoryExists creates a directory if it doesn't exist
func EnsureDirectoryExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// GenerateUniqueFilename creates a unique filename based on timestamp and hash
func GenerateUniqueFilename(originalName string) string {
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	hash := sha256.Sum256([]byte(originalName + timestamp))
	hashString := hex.EncodeToString(hash[:])[:12]

	ext := filepath.Ext(originalName)
	return hashString + ext
}

// SaveUploadedFile saves a file from an HTTP request to disk with validation
func SaveUploadedFile(file *multipart.FileHeader) (string, error) {
	// Create upload directory if it doesn't exist
	if err := EnsureDirectoryExists(UploadDir); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Validate file size
	if file.Size > MaxFileSize {
		return "", fmt.Errorf("file too large (max %d bytes)", MaxFileSize)
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Detect file type by reading first 512 bytes
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Reset file pointer
	src.Seek(0, io.SeekStart)

	// Get content type and validate
	contentType := http.DetectContentType(buffer)
	if !supportedImageTypes[contentType] && strings.HasPrefix(file.Filename, "image/") {
		return "", fmt.Errorf("unsupported image format: %s", contentType)
	}

	// Generate unique filename to prevent overwrites
	uniqueFilename := GenerateUniqueFilename(file.Filename)
	filePath := filepath.Join(UploadDir, uniqueFilename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy file contents
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return filePath, nil
}

// CleanupTempFiles removes temporary files older than the specified duration
func CleanupTempFiles(maxAge time.Duration) error {
	// Ensure temp directory exists
	if err := EnsureDirectoryExists(TempDir); err != nil {
		return err
	}

	// Get current time
	now := time.Now()

	// Walk through temp directory
	return filepath.Walk(TempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directory itself
		if path == TempDir {
			return nil
		}

		// Check if file is older than maxAge
		if now.Sub(info.ModTime()) > maxAge {
			if err := os.Remove(path); err != nil {
				return err
			}
		}

		return nil
	})
}

// SaveOutputFile saves processed output to a file and returns the path
func SaveOutputFile(data []byte, extension string) (string, error) {
	// Create temp directory if it doesn't exist
	if err := EnsureDirectoryExists(TempDir); err != nil {
		return "", err
	}

	// Generate unique filename
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	hash := sha256.Sum256([]byte(timestamp))
	filename := hex.EncodeToString(hash[:])[:12] + extension
	filePath := filepath.Join(TempDir, filename)

	// Write data to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", err
	}

	return filePath, nil
}
