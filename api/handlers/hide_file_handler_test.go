package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pranaykumar2/steg-go/api/utils"
	"github.com/stretchr/testify/assert"
)

func TestHideFileAPI_WithPassword(t *testing.T) {
	coverImagePath, cleanupCoverImage := createTestImageFile(t, "cover_hide_file_pass.png")
	defer cleanupCoverImage()

	fileToHidePath, cleanupFileToHide := createTestTextFile(t, "this is a secret file content", "secret.txt")
	defer cleanupFileToHide()

	params := map[string]string{
		"password": "securefilepassword",
	}
	fileParams := map[string]string{
		"image": coverImagePath,
		"file":  fileToHidePath,
	}

	req, err := newFileUploadRequest("/api/hideFile", params, fileParams)
	assert.NoError(t, err)

	rr := serveHTTP(req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var response utils.SuccessResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	responseData, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)

	assert.NotEmpty(t, responseData["outputFileURL"])
	assert.Equal(t, "enabled", responseData["encryption"])
	
	fileDetails, ok := responseData["fileDetails"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "secret.txt", fileDetails["originalName"])
}

func TestHideFileAPI_WithoutPassword(t *testing.T) {
	coverImagePath, cleanupCoverImage := createTestImageFile(t, "cover_hide_file_no_pass.png")
	defer cleanupCoverImage()

	fileToHidePath, cleanupFileToHide := createTestTextFile(t, "this is a public file content", "public.txt")
	defer cleanupFileToHide()

	params := map[string]string{} // No password
	fileParams := map[string]string{
		"image": coverImagePath,
		"file":  fileToHidePath,
	}

	req, err := newFileUploadRequest("/api/hideFile", params, fileParams)
	assert.NoError(t, err)

	rr := serveHTTP(req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var response utils.SuccessResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	responseData, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)

	assert.NotEmpty(t, responseData["outputFileURL"])
	assert.Equal(t, "disabled", responseData["encryption"])

	fileDetails, ok := responseData["fileDetails"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "public.txt", fileDetails["originalName"])
}


func TestHideFileAPI_NoCoverImage(t *testing.T) {
	fileToHidePath, cleanupFileToHide := createTestTextFile(t, "some content", "file.txt")
	defer cleanupFileToHide()

	params := map[string]string{}
	fileParams := map[string]string{
		"file": fileToHidePath,
	}

	req, err := newFileUploadRequest("/api/hideFile", params, fileParams)
	assert.NoError(t, err)

	rr := serveHTTP(req)
	checkAPIErrorResponse(t, rr, http.StatusBadRequest, "No cover image uploaded")
}

func TestHideFileAPI_NoFileToHide(t *testing.T) {
	coverImagePath, cleanupCoverImage := createTestImageFile(t, "cover_no_file.png")
	defer cleanupCoverImage()

	params := map[string]string{}
	fileParams := map[string]string{
		"image": coverImagePath,
	}

	req, err := newFileUploadRequest("/api/hideFile", params, fileParams)
	assert.NoError(t, err)

	rr := serveHTTP(req)
	checkAPIErrorResponse(t, rr, http.StatusBadRequest, "No file to hide uploaded")
}
