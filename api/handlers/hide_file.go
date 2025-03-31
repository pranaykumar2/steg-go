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

type HideFileResponse struct {
	Key           string `json:"key"`
	OutputFileURL string `json:"outputFileURL"`
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

	imagePath, err := utils.SaveUploadedFile(imageFile)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to save cover image: "+err.Error())
		return
	}

	tempFileToHide, err := os.CreateTemp(utils.TempDir, "file_to_hide_*")
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create temporary file: "+err.Error())
		return
	}
	defer os.Remove(tempFileToHide.Name())
	defer tempFileToHide.Close()
	src, err := fileToHide.Open()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to open uploaded file: "+err.Error())
		return
	}
	defer src.Close()
	if _, err = io.Copy(tempFileToHide, src); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to copy file: "+err.Error())
		return
	}

	fileHandler := steganography.NewFileHandler()

	fileData, metadata, err := fileHandler.ReadFileContent(tempFileToHide.Name())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to read file: "+err.Error())
		return
	}

	metadata.OriginalName = fileToHide.Filename
	metadata.FileExt = filepath.Ext(fileToHide.Filename)

	encoder, err := steganography.NewEncoder(imagePath)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to initialize encoder: "+err.Error())
		return
	}

	encryptor, err := crypto.NewEncryptor()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to initialize encryption: "+err.Error())
		return
	}

	encrypted, err := encryptor.Encrypt(fileData)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to encrypt file: "+err.Error())
		return
	}

	if err := encoder.HideFile(encrypted, metadata); err != nil {
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

	keyHex := hex.EncodeToString(encryptor.GetKey())

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
