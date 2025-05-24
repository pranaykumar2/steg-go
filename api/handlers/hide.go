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
	Message  string `form:"message" binding:"required"` // Changed from json to form
	Password string `form:"password,omitempty"`         // Added password, using form tag
}

type HideTextResponse struct {
	OutputFileURL string `json:"outputFileURL"`
	Encryption    string `json:"encryption"` // To indicate if encryption was used
}

func HideText(c *gin.Context) {
	// Parse multipart form, which is necessary for file uploads and form fields
	if err := c.Request.ParseMultipartForm(utils.MaxFileSize); err != nil {
		utils.ValidationErrorResponse(c, "Invalid form data: "+err.Error())
		return
	}

	var req HideTextRequest
	// Bind form data (including message and password) to the struct
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

	var encoder *steganography.Encoder
	if req.Password != "" {
		encoder, err = steganography.NewEncoder(inputPath, req.Password)
	} else {
		encoder, err = steganography.NewEncoder(inputPath)
	}
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to initialize encoder: "+err.Error())
		return
	}

	// Encryption is handled by NewEncoder if password is provided.
	// The original message is passed directly to Hide.
	if err := encoder.Hide([]byte(req.Message)); err != nil {
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
	encryptionStatus := "disabled"
	if req.Password != "" {
		encryptionStatus = "enabled"
	}

	utils.SuccessResponse(c, http.StatusOK, "Message hidden successfully", HideTextResponse{
		OutputFileURL: outputURL,
		Encryption:    encryptionStatus,
	})
}
