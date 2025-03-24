package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/pranaykumar2/steg-go/api/utils"
)

// ServeFile handles requests to serve generated files
func ServeFile(c *gin.Context) {
	filename := c.Param("filename")

	// Prevent directory traversal attacks
	if filepath.Ext(filename) == "" || filepath.Base(filename) != filename {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid filename")
		return
	}

	// Check in temp directory first
	tempPath := filepath.Join(utils.TempDir, filename)
	if _, err := os.Stat(tempPath); err == nil {
		c.File(tempPath)
		return
	}

	// Then check in uploads directory
	uploadPath := filepath.Join(utils.UploadDir, filename)
	if _, err := os.Stat(uploadPath); err == nil {
		c.File(uploadPath)
		return
	}

	utils.NotFoundResponse(c, "File not found")
}
