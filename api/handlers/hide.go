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

type HideTextRequest struct {
	Message string `json:"message" binding:"required"`
}

type HideTextResponse struct {
	Key           string `json:"key"`
	OutputFileURL string `json:"outputFileURL"`
}

func HideText(c *gin.Context) {
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(utils.MaxFileSize); err != nil {
		utils.ValidationErrorResponse(c, "Invalid form data: "+err.Error())
		return
	}

	var req HideTextRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request: "+err.Error())
		return
	}

	if req.Message == "" {
		utils.ValidationErrorResponse(c, "Message cannot be empty")
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		utils.ValidationErrorResponse(c, "No image file uploaded")
		return
	}

	inputPath, err := utils.SaveUploadedFile(file)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to save uploaded file: "+err.Error())
		return
	}

	encoder, err := steganography.NewEncoder(inputPath)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to initialize encoder: "+err.Error())
		return
	}

	encryptor, err := crypto.NewEncryptor()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to initialize encryption: "+err.Error())
		return
	}

	encrypted, err := encryptor.Encrypt([]byte(req.Message))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to encrypt message: "+err.Error())
		return
	}

	if err := encoder.Hide(encrypted); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to hide message: "+err.Error())
		return
	}

	uniqueName := utils.GenerateUniqueFilename(file.Filename)
	outputPath := filepath.Join(utils.TempDir, "stego_"+uniqueName)

	if err := encoder.SaveOutput(outputPath); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save output image: "+err.Error())
		return
	}

	outputURL := "/api/files/" + filepath.Base(outputPath)

	keyHex := hex.EncodeToString(encryptor.GetKey())

	utils.SuccessResponse(c, http.StatusOK, "Message hidden successfully", HideTextResponse{
		Key:           keyHex,
		OutputFileURL: outputURL,
	})
}
