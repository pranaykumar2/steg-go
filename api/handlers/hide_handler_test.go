package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pranaykumar2/steg-go/api/utils"
	"github.com/stretchr/testify/assert"
)

func TestHideTextAPI_WithPassword(t *testing.T) {
	imagePath, cleanupImage := createTestImageFile(t, "test_hide_with_pass.png")
	defer cleanupImage()

	params := map[string]string{
		"message":  "secret message with password",
		"password": "supersecretpassword",
	}
	fileParams := map[string]string{
		"image": imagePath,
	}

	req, err := newFileUploadRequest("/api/hide", params, fileParams)
	assert.NoError(t, err)

	rr := serveHTTP(req) // Using the global testRouter from main_test.go

	assert.Equal(t, http.StatusOK, rr.Code)

	var response utils.SuccessResponse // Assuming HideTextResponse is wrapped in SuccessResponse
	// The actual data is in response.Data, which needs to be asserted against HideTextResponse struct
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	// Assert specific fields within response.Data (which should be HideTextResponse)
	// Need to type assert response.Data or unmarshal it into the specific HideTextResponse type
	responseData, ok := response.Data.(map[string]interface{})
	assert.True(t, ok, "response.Data is not a map[string]interface{}")

	assert.NotEmpty(t, responseData["outputFileURL"])
	assert.Equal(t, "enabled", responseData["encryption"])
}

func TestHideTextAPI_WithoutPassword(t *testing.T) {
	imagePath, cleanupImage := createTestImageFile(t, "test_hide_no_pass.png")
	defer cleanupImage()

	params := map[string]string{
		"message": "secret message no password",
	}
	fileParams := map[string]string{
		"image": imagePath,
	}

	req, err := newFileUploadRequest("/api/hide", params, fileParams)
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
}

func TestHideTextAPI_NoMessage(t *testing.T) {
	imagePath, cleanupImage := createTestImageFile(t, "test_hide_no_message.png")
	defer cleanupImage()

	params := map[string]string{
		"message": "", // Empty message
	}
	fileParams := map[string]string{
		"image": imagePath,
	}

	req, err := newFileUploadRequest("/api/hide", params, fileParams)
	assert.NoError(t, err)
	
	rr := serveHTTP(req)
	checkAPIErrorResponse(t, rr, http.StatusBadRequest, "Message cannot be empty")
}

func TestHideTextAPI_NoImage(t *testing.T) {
	params := map[string]string{
		"message": "test message",
	}
	// No image file
	fileParams := map[string]string{}

	req, err := newFileUploadRequest("/api/hide", params, fileParams)
	assert.NoError(t, err)

	rr := serveHTTP(req)
	checkAPIErrorResponse(t, rr, http.StatusBadRequest, "No image file uploaded")
}
