package handlers

import (
	"encoding/hex"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/pranaykumar2/steg-go/api/utils"
	"github.com/pranaykumar2/steg-go/internal/crypto"
	"github.com/pranaykumar2/steg-go/internal/steganography"
)

// ExtractRequest represents the request to extract content from an image
type ExtractRequest struct {
	Key string `json:"key" binding:"required"`
}

// ExtractResponse represents the response after extracting content
type ExtractResponse struct {
	IsFile      bool   `json:"isFile"`
	Message     string `json:"message,omitempty"`
	FileURL     string `json:"fileURL,omitempty"`
	FileName    string `json:"fileName,omitempty"`
	FileType    string `json:"fileType,omitempty"`
	FileSize    int64  `json:"fileSize,omitempty"`
	ContentType string `json:"contentType,omitempty"`
}

// Extract handles requests to extract hidden content from an image
func Extract(c *gin.Context) {
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(utils.MaxFileSize); err != nil {
		utils.ValidationErrorResponse(c, "Invalid form data: "+err.Error())
		return
	}

	// Bind the request
	var req ExtractRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request: "+err.Error())
		return
	}

	// Validate key
	if len(req.Key) != 64 {
		utils.ValidationErrorResponse(c, "Invalid key length. Expected 64 hexadecimal characters")
		return
	}

	// Decode key
	key, err := hex.DecodeString(req.Key)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid key format. Must be hexadecimal")
		return
	}

	// Get uploaded image file
	file, err := c.FormFile("image")
	if err != nil {
		utils.ValidationErrorResponse(c, "No image file uploaded")
		return
	}

	// Save uploaded file
	imagePath, err := utils.SaveUploadedFile(file)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to save uploaded file: "+err.Error())
		return
	}

	// Initialize decoder
	decoder, err := steganography.NewDecoder(imagePath)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to initialize decoder: "+err.Error())
		return
	}

	// Extract data
	data, isFile, metadata, err := decoder.Extract()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to extract data: "+err.Error())
		return
	}

	// Initialize encryptor with provided key
	encryptor, err := crypto.NewEncryptorWithKey(key)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to initialize decryption: "+err.Error())
		return
	}

	// Decrypt data
	decrypted, err := encryptor.Decrypt(data)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to decrypt data: "+err.Error())
		return
	}

	// Prepare response
	response := ExtractResponse{
		IsFile: isFile,
	}

	if isFile && metadata != nil {
		// Create temporary file for the extracted content
		outputPath := filepath.Join(utils.TempDir, metadata.OriginalName)

		// Save the extracted file
		fileHandler := steganography.NewFileHandler()
		if err := fileHandler.SaveFileContent(decrypted, metadata, outputPath); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save extracted file: "+err.Error())
			return
		}

		// Generate public URL for the extracted file
		fileURL := "/api/files/" + filepath.Base(outputPath)

		// Set file details in response
		response.FileURL = fileURL
		response.FileName = metadata.OriginalName
		response.FileType = metadata.FileExt
		response.FileSize = int64(metadata.FileSize)

		// Determine content type
		switch filepath.Ext(metadata.OriginalName) {
		case ".pdf":
			response.ContentType = "application/pdf"
		case ".txt":
			response.ContentType = "text/plain"
		case ".jpg", ".jpeg":
			response.ContentType = "image/jpeg"
		case ".png":
			response.ContentType = "image/png"
		case ".mp3":
			response.ContentType = "audio/mpeg"
		default:
			response.ContentType = "application/octet-stream"
		}
	} else {
		// It's a text message
		response.Message = string(decrypted)
	}

	// Return success response
	utils.SuccessResponse(c, http.StatusOK, "Content extracted successfully", response)
}
