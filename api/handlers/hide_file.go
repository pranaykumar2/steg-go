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
	// "github.com/pranaykumar2/steg-go/internal/crypto" // No longer directly used
)

type HideFileResponse struct {
	OutputFileURL string `json:"outputFileURL"`
	Encryption    string `json:"encryption"` // To indicate if encryption was used
	FileDetails   struct {
		OriginalName string `json:"originalName"`
		FileType     string `json:"fileType"`
		FileSize     int64  `json:"fileSize"`
	} `json:"fileDetails"`
}

func HideFile(c *gin.Context) {
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(utils.MaxFileSize); err != nil {
		utils.ValidationErrorResponse(c, "Invalid form data: "+err.Error())
		return
	}

	imageFile, err := c.FormFile("image")
	if err != nil {
		utils.ValidationErrorResponse(c, "No cover image uploaded")
		return
	}

	fileToHide, err := c.FormFile("file")
	if err != nil {
		utils.ValidationErrorResponse(c, "No file to hide uploaded")
		return
	}

	// Retrieve password from form data
	password := c.Request.FormValue("password")

	imagePath, err := utils.SaveUploadedFile(imageFile)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to save cover image: "+err.Error())
		return
	}

	// Save the file to hide to a temporary path to pass to ReadFileContent
	tempFileToHidePath, err := utils.SaveUploadedFile(fileToHide)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to save file to hide: "+err.Error())
		return
	}
	defer os.Remove(tempFileToHidePath) // Clean up the temp file

	fileHandler := steganography.NewFileHandler()
	// ReadFileContent expects a path, so we use the saved temp file path
	fileData, metadata, err := fileHandler.ReadFileContent(tempFileToHidePath)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to read file to hide: "+err.Error())
		return
	}
	// Ensure OriginalName and FileExt are from the uploaded file's metadata, not temp path
	metadata.OriginalName = fileToHide.Filename
	metadata.FileExt = filepath.Ext(fileToHide.Filename)


	var encoder *steganography.Encoder
	if password != "" {
		encoder, err = steganography.NewEncoder(imagePath, password)
	} else {
		encoder, err = steganography.NewEncoder(imagePath)
	}
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to initialize encoder: "+err.Error())
		return
	}

	// Encryption is handled by NewEncoder if password is provided.
	// Pass the original fileData to HideFile.
	if err := encoder.HideFile(fileData, metadata); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to hide file: "+err.Error())
		return
	}

	uniqueName := utils.GenerateUniqueFilename(imageFile.Filename)
	outputPath := filepath.Join(utils.TempDir, "stego_"+uniqueName)

	if err := encoder.SaveOutput(outputPath); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save output image: "+err.Error())
		return
	}

	outputURL := "/api/files/" + filepath.Base(outputPath)
	encryptionStatus := "disabled"
	if password != "" {
		encryptionStatus = "enabled"
	}

	response := HideFileResponse{
		OutputFileURL: outputURL,
		Encryption:    encryptionStatus,
	}
	response.FileDetails.OriginalName = metadata.OriginalName
	response.FileDetails.FileType = metadata.FileExt
	response.FileDetails.FileSize = metadata.FileSize // Assuming FileSize is already int64

	// Return success response
	utils.SuccessResponse(c, http.StatusOK, "File hidden successfully", response)
}
