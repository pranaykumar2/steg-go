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

// HideTextRequest represents the request to hide text in an image
type HideTextRequest struct {
	Message string `json:"message" binding:"required"`
}

// HideTextResponse represents the response after hiding text
type HideTextResponse struct {
	Key           string `json:"key"`
	OutputFileURL string `json:"outputFileURL"`
}

// HideText handles requests to hide text in an image
func HideText(c *gin.Context) {
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(utils.MaxFileSize); err != nil {
		utils.ValidationErrorResponse(c, "Invalid form data: "+err.Error())
		return
	}

	// Get text message from form
	var req HideTextRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request: "+err.Error())
		return
	}

	if req.Message == "" {
		utils.ValidationErrorResponse(c, "Message cannot be empty")
		return
	}

	// Get uploaded image file
	file, err := c.FormFile("image")
	if err != nil {
		utils.ValidationErrorResponse(c, "No image file uploaded")
		return
	}

	// Save uploaded file
	inputPath, err := utils.SaveUploadedFile(file)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to save uploaded file: "+err.Error())
		return
	}

	// Initialize steganography encoder
	encoder, err := steganography.NewEncoder(inputPath)
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

	// Encrypt message
	encrypted, err := encryptor.Encrypt([]byte(req.Message))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to encrypt message: "+err.Error())
		return
	}

	// Hide message in image
	if err := encoder.Hide(encrypted); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to hide message: "+err.Error())
		return
	}

	// Generate output filename
	uniqueName := utils.GenerateUniqueFilename(file.Filename)
	outputPath := filepath.Join(utils.TempDir, "stego_"+uniqueName)

	// Save output image
	if err := encoder.SaveOutput(outputPath); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save output image: "+err.Error())
		return
	}

	// Generate public URL for the output file
	// In a real-world scenario, this would be a publicly accessible URL
	outputURL := "/api/files/" + filepath.Base(outputPath)

	// Get encryption key
	keyHex := hex.EncodeToString(encryptor.GetKey())

	// Return success response with key and file URL
	utils.SuccessResponse(c, http.StatusOK, "Message hidden successfully", HideTextResponse{
		Key:           keyHex,
		OutputFileURL: outputURL,
	})
}
