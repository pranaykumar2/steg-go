package handlers

import (
	"encoding/hex"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/pranaykumar2/steg-go/api/utils"
	"github.com/pranaykumar2/steg-go/internal/crypto"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pranaykumar2/steg-go/api/utils"
	"github.com/pranaykumar2/steg-go/internal/steganography"
	// "github.com/pranaykumar2/steg-go/internal/crypto" // No longer directly used here
)

// ExtractRequest might not be needed if all inputs are from form-data
// type ExtractRequest struct {
// No fields needed from JSON body if password and image are from form
// }

type ExtractResponse struct {
	IsFile      bool   `json:"isFile"`
	Message     string `json:"message,omitempty"`
	FileURL     string `json:"fileURL,omitempty"`
	FileName    string `json:"fileName,omitempty"`
	FileType    string `json:"fileType,omitempty"`
	FileSize    int64  `json:"fileSize,omitempty"`
	ContentType string `json:"contentType,omitempty"`
}

func Extract(c *gin.Context) {
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(utils.MaxFileSize); err != nil {
		utils.ValidationErrorResponse(c, "Invalid form data: "+err.Error())
		return
	}

	// Password will be read from form value
	password := c.Request.FormValue("password")

	file, err := c.FormFile("image")
	if err != nil {
		utils.ValidationErrorResponse(c, "No image file uploaded")
		return
	}

	imagePath, err := utils.SaveUploadedFile(file)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to save uploaded file: "+err.Error())
		return
	}

	var decoder *steganography.Decoder
	if password != "" {
		decoder, err = steganography.NewDecoder(imagePath, password)
	} else {
		decoder, err = steganography.NewDecoder(imagePath)
	}
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to initialize decoder: "+err.Error())
		return
	}

	// stegoFlags (4th return) is ignored for now in API response
	extractedData, isFile, fileMetadata, _, err := decoder.Extract()
	if err != nil {
		if strings.Contains(err.Error(), "password required for encrypted data") {
			utils.ErrorResponse(c, http.StatusBadRequest, "Password required to decrypt this image.")
			return
		}
		if strings.Contains(err.Error(), "failed to decrypt data") { // This implies wrong password
			utils.ErrorResponse(c, http.StatusUnauthorized, "Decryption failed. Incorrect password.")
			return
		}
		if strings.Contains(err.Error(), "no steganographic data found") {
			utils.ErrorResponse(c, http.StatusNotFound, "No hidden content found in this image.")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to extract data: "+err.Error())
		return
	}

	// Decryption is now handled within decoder.Extract if password was provided.
	// The 'extractedData' is the final decrypted data.

	response := ExtractResponse{
		IsFile: isFile,
	}

	if isFile && fileMetadata != nil {
		// Save the extracted file to a temporary location to make it available via URL
		// Ensure the filename is safe to use in a path
		safeFileName := filepath.Base(fileMetadata.OriginalName) // Use only the filename part
		outputPath := filepath.Join(utils.TempDir, "extracted_"+utils.GenerateUniqueFilename(safeFileName))
		
		fileHandler := steganography.NewFileHandler()
		// Use extractedData (which is already decrypted)
		if err := fileHandler.SaveFileContent(extractedData, fileMetadata, outputPath); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save extracted file: "+err.Error())
			return
		}

		fileURL := "/api/files/" + filepath.Base(outputPath)

		response.FileURL = fileURL
		response.FileName = fileMetadata.OriginalName
		response.FileType = fileMetadata.FileExt
		response.FileSize = fileMetadata.FileSize // Assuming FileSize in FileMetadata is already int64

		// Determine ContentType (optional, can be set by client or derived)
		// This is a simplified version. A more robust solution would use mime.TypeByExtension.
		switch filepath.Ext(safeFileName) {
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
			response.ContentType = "application/octet-stream" // Generic binary
		}
	} else {
		response.Message = string(extractedData)
	}

	utils.SuccessResponse(c, http.StatusOK, "Content extracted successfully", response)
}
