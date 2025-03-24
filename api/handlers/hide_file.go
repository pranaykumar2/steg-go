package handlers

import (
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/pranaykumar2/steg-go/api/utils"
	"github.com/pranaykumar2/steg-go/internal/crypto"
	"github.com/pranaykumar2/steg-go/internal/steganography"
)

// HideFileResponse represents the response after hiding a file
type HideFileResponse struct {
	Key           string `json:"key"`
	OutputFileURL string `json:"outputFileURL"`
	FileDetails   struct {
		OriginalName string `json:"originalName"`
		FileType     string `json:"fileType"`
		FileSize     int64  `json:"fileSize"`
	} `json:"fileDetails"`
}

// HideFile handles requests to hide a file in an image
func HideFile(c *gin.Context) {
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(utils.MaxFileSize); err != nil {
		utils.ValidationErrorResponse(c, "Invalid form data: "+err.Error())
		return
	}

	// Get uploaded image file
	imageFile, err := c.FormFile("image")
	if err != nil {
		utils.ValidationErrorResponse(c, "No cover image uploaded")
		return
	}

	// Get file to hide
	fileToHide, err := c.FormFile("file")
	if err != nil {
		utils.ValidationErrorResponse(c, "No file to hide uploaded")
		return
	}

	// Save uploaded cover image
	imagePath, err := utils.SaveUploadedFile(imageFile)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to save cover image: "+err.Error())
		return
	}

	// Create temporary file for the file to hide
	tempFileToHide, err := os.CreateTemp(utils.TempDir, "file_to_hide_*")
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create temporary file: "+err.Error())
		return
	}
	defer os.Remove(tempFileToHide.Name())
	defer tempFileToHide.Close()

	// Open the uploaded file
	src, err := fileToHide.Open()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to open uploaded file: "+err.Error())
		return
	}
	defer src.Close()

	// Copy the uploaded file to the temporary file
	if _, err = io.Copy(tempFileToHide, src); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to copy file: "+err.Error())
		return
	}

	// Initialize file handler
	fileHandler := steganography.NewFileHandler()

	// Read file content
	fileData, metadata, err := fileHandler.ReadFileContent(tempFileToHide.Name())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to read file: "+err.Error())
		return
	}

	// Initialize steganography encoder
	encoder, err := steganography.NewEncoder(imagePath)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to initialize encoder: "+err.Error())
		return
	}

	// Initialize encryptor
	encryptor, err := crypto.NewEncryptor()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to initialize encryption: "+err.Error())
		return
	}

	// Encrypt file data
	encrypted, err := encryptor.Encrypt(fileData)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to encrypt file: "+err.Error())
		return
	}

	// Hide file in image
	if err := encoder.HideFile(encrypted, metadata); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to hide file: "+err.Error())
		return
	}

	// Generate output filename
	uniqueName := utils.GenerateUniqueFilename(imageFile.Filename)
	outputPath := filepath.Join(utils.TempDir, "stego_"+uniqueName)

	// Save output image
	if err := encoder.SaveOutput(outputPath); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save output image: "+err.Error())
		return
	}

	// Generate public URL for the output file
	outputURL := "/api/files/" + filepath.Base(outputPath)

	// Get encryption key
	keyHex := hex.EncodeToString(encryptor.GetKey())

	// Prepare response
	response := HideFileResponse{
		Key:           keyHex,
		OutputFileURL: outputURL,
	}
	response.FileDetails.OriginalName = metadata.OriginalName
	response.FileDetails.FileType = metadata.FileExt
	response.FileDetails.FileSize = int64(metadata.FileSize)

	// Return success response
	utils.SuccessResponse(c, http.StatusOK, "File hidden successfully", response)
}
