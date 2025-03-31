package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"github.com/gin-gonic/gin"
	"github.com/pranaykumar2/steg-go/api/utils"
)

func ServeFile(c *gin.Context) {
	filename := c.Param("filename")

	if filepath.Ext(filename) == "" || filepath.Base(filename) != filename {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid filename")
		return
	}

	tempPath := filepath.Join(utils.TempDir, filename)
	if _, err := os.Stat(tempPath); err == nil {
		c.File(tempPath)
		return
	}

	uploadPath := filepath.Join(utils.UploadDir, filename)
	if _, err := os.Stat(uploadPath); err == nil {
		c.File(uploadPath)
		return
	}

	utils.NotFoundResponse(c, "File not found")
}
