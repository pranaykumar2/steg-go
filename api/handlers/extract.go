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

type ExtractRequest struct {
	Key string `json:"key" binding:"required"`
}

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

	var req ExtractRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request: "+err.Error())
		return
	}

	if len(req.Key) != 64 {
		utils.ValidationErrorResponse(c, "Invalid key length. Expected 64 hexadecimal characters")
		return
	}

	key, err := hex.DecodeString(req.Key)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid key format. Must be hexadecimal")
		return
	}

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

	decoder, err := steganography.NewDecoder(imagePath)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to initialize decoder: "+err.Error())
		return
	}

	data, isFile, metadata, err := decoder.Extract()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to extract data: "+err.Error())
		return
	}

	encryptor, err := crypto.NewEncryptorWithKey(key)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to initialize decryption: "+err.Error())
		return
	}

	decrypted, err := encryptor.Decrypt(data)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to decrypt data: "+err.Error())
		return
	}

	response := ExtractResponse{
		IsFile: isFile,
	}

	if isFile && metadata != nil {
		outputPath := filepath.Join(utils.TempDir, metadata.OriginalName)
		fileHandler := steganography.NewFileHandler()
		if err := fileHandler.SaveFileContent(decrypted, metadata, outputPath); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save extracted file: "+err.Error())
			return
		}

		fileURL := "/api/files/" + filepath.Base(outputPath)

		response.FileURL = fileURL
		response.FileName = metadata.OriginalName
		response.FileType = metadata.FileExt
		response.FileSize = int64(metadata.FileSize)

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
		response.Message = string(decrypted)
	}

	utils.SuccessResponse(c, http.StatusOK, "Content extracted successfully", response)
}
