package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pranaykumar2/steg-go/api/utils"
	"github.com/pranaykumar2/steg-go/pkg/exiftools"
)

// MetadataResponse represents the response for image metadata analysis
type MetadataResponse struct {
	Filename     string            `json:"filename"`
	FileSize     int64             `json:"fileSize"`
	FileType     string            `json:"fileType"`
	MimeType     string            `json:"mimeType"`
	ModTime      string            `json:"modTime"`
	ImageWidth   int               `json:"imageWidth"`
	ImageHeight  int               `json:"imageHeight"`
	HasEXIF      bool              `json:"hasEXIF"`
	PrivacyRisks []string          `json:"privacyRisks"`
	Properties   map[string]string `json:"properties"`

	// Steganography specific information
	SteganoCapacity struct {
		Bytes     int     `json:"bytes"`
		Kilobytes float64 `json:"kilobytes"`
		Megabytes float64 `json:"megabytes"`
		Text      struct {
			Characters int `json:"characters"`
			Words      int `json:"words"`
		} `json:"text"`
	} `json:"steganoCapacity"`
}

// AnalyzeMetadata handles requests to analyze image metadata
func AnalyzeMetadata(c *gin.Context) {
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

	// Get metadata
	metadata, err := exiftools.GetImageMetadata(imagePath)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to extract metadata: "+err.Error())
		return
	}

	// Calculate steganography capacity
	capacityBytes := (metadata.ImageWidth * metadata.ImageHeight * 3) / 8

	// Prepare response
	response := MetadataResponse{
		Filename:     metadata.Filename,
		FileSize:     metadata.FileSize,
		FileType:     metadata.FileType,
		MimeType:     metadata.MimeType,
		ModTime:      metadata.ModTime.Format("2006-01-02T15:04:05Z"),
		ImageWidth:   metadata.ImageWidth,
		ImageHeight:  metadata.ImageHeight,
		HasEXIF:      metadata.HasEXIF,
		PrivacyRisks: metadata.PrivacyRisks,
		Properties:   metadata.Properties,
	}

	// Set steganography capacity
	response.SteganoCapacity.Bytes = capacityBytes
	response.SteganoCapacity.Kilobytes = float64(capacityBytes) / 1024
	response.SteganoCapacity.Megabytes = float64(capacityBytes) / (1024 * 1024)
	response.SteganoCapacity.Text.Characters = capacityBytes
	response.SteganoCapacity.Text.Words = capacityBytes / 5 // Rough estimate assuming average word length of 5 chars

	// Return success response
	utils.SuccessResponse(c, http.StatusOK, "Metadata analysis completed", response)
}
