package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/pranaykumar2/steg-go/api/utils"
	"github.com/stretchr/testify/assert"
)

// TestExtractAPI_EncryptedContent_CorrectPassword tests extracting content from an image
// that was hidden with a password, using the correct password.
func TestExtractAPI_EncryptedContent_CorrectPassword(t *testing.T) {
	originalImagePath, cleanupOriginal := createTestImageFile(t, "original_for_extract_enc_correct.png")
	defer cleanupOriginal()

	secretMessage := "This is a super secret message for extraction with password."
	password := "testPassword123"
	stegoImageName := "stego_extract_enc_correct.png"

	// Create an image with an encrypted message
	stegoImagePath, cleanupStego := createTestStegoImage(t, originalImagePath, secretMessage, password, stegoImageName)
	defer cleanupStego()

	// Prepare request for extraction
	params := map[string]string{
		"password": password, // Correct password
	}
	fileParams := map[string]string{
		"image": stegoImagePath,
	}

	req, err := newFileUploadRequest("/api/extract", params, fileParams)
	assert.NoError(t, err)

	rr := serveHTTP(req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var response utils.SuccessResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success, "API success flag should be true")

	responseData, ok := response.Data.(map[string]interface{})
	assert.True(t, ok, "Response data is not map[string]interface{}")

	assert.False(t, responseData["isFile"].(bool), "IsFile should be false for text message")
	assert.Equal(t, secretMessage, responseData["message"].(string), "Extracted message mismatch")
}

// TestExtractAPI_EncryptedContent_IncorrectPassword tests extracting content from an image
// that was hidden with a password, using an incorrect password.
func TestExtractAPI_EncryptedContent_IncorrectPassword(t *testing.T) {
	originalImagePath, cleanupOriginal := createTestImageFile(t, "original_for_extract_enc_incorrect.png")
	defer cleanupOriginal()

	secretMessage := "This message won't be extracted."
	correctPassword := "correctPassword"
	incorrectPassword := "wrongPassword"
	stegoImageName := "stego_extract_enc_incorrect.png"

	stegoImagePath, cleanupStego := createTestStegoImage(t, originalImagePath, secretMessage, correctPassword, stegoImageName)
	defer cleanupStego()

	params := map[string]string{
		"password": incorrectPassword,
	}
	fileParams := map[string]string{
		"image": stegoImagePath,
	}

	req, err := newFileUploadRequest("/api/extract", params, fileParams)
	assert.NoError(t, err)

	rr := serveHTTP(req)
	// Expecting Unauthorized or similar error due to incorrect password leading to decryption failure
	checkAPIErrorResponse(t, rr, http.StatusUnauthorized, "Decryption failed. Incorrect password.")
}

// TestExtractAPI_EncryptedContent_NoPassword tests extracting content from an image
// that was hidden with a password, but providing no password during extraction.
func TestExtractAPI_EncryptedContent_NoPassword(t *testing.T) {
	originalImagePath, cleanupOriginal := createTestImageFile(t, "original_for_extract_enc_nopass.png")
	defer cleanupOriginal()

	secretMessage := "This message requires a password."
	password := "passwordIsNeeded"
	stegoImageName := "stego_extract_enc_nopass.png"

	stegoImagePath, cleanupStego := createTestStegoImage(t, originalImagePath, secretMessage, password, stegoImageName)
	defer cleanupStego()

	params := map[string]string{
		// No password provided
	}
	fileParams := map[string]string{
		"image": stegoImagePath,
	}

	req, err := newFileUploadRequest("/api/extract", params, fileParams)
	assert.NoError(t, err)

	rr := serveHTTP(req)
	// Expecting Bad Request or similar error because password is required
	checkAPIErrorResponse(t, rr, http.StatusBadRequest, "Password required to decrypt this image.")
}

// TestExtractAPI_NonEncryptedContent_NoPassword tests extracting content from an image
// that was hidden without a password, and providing no password during extraction.
func TestExtractAPI_NonEncryptedContent_NoPassword(t *testing.T) {
	originalImagePath, cleanupOriginal := createTestImageFile(t, "original_for_extract_noenc_nopass.png")
	defer cleanupOriginal()

	secretMessage := "This is a public message, no password needed."
	stegoImageName := "stego_extract_noenc_nopass.png"

	// Create an image with a non-encrypted message (empty password)
	stegoImagePath, cleanupStego := createTestStegoImage(t, originalImagePath, secretMessage, "", stegoImageName)
	defer cleanupStego()

	params := map[string]string{
		// No password provided
	}
	fileParams := map[string]string{
		"image": stegoImagePath,
	}

	req, err := newFileUploadRequest("/api/extract", params, fileParams)
	assert.NoError(t, err)

	rr := serveHTTP(req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var response utils.SuccessResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	
	responseData, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)

	assert.False(t, responseData["isFile"].(bool))
	assert.Equal(t, secretMessage, responseData["message"].(string))
}

// TestExtractAPI_NonEncryptedContent_WithPassword tests extracting content from an image
// that was hidden without a password, but providing a password during extraction (should be ignored).
func TestExtractAPI_NonEncryptedContent_WithPassword(t *testing.T) {
	originalImagePath, cleanupOriginal := createTestImageFile(t, "original_for_extract_noenc_withpass.png")
	defer cleanupOriginal()

	secretMessage := "Public message, password should be ignored."
	stegoImageName := "stego_extract_noenc_withpass.png"

	stegoImagePath, cleanupStego := createTestStegoImage(t, originalImagePath, secretMessage, "", stegoImageName)
	defer cleanupStego()

	params := map[string]string{
		"password": "someRandomPassword", // Password provided but not needed
	}
	fileParams := map[string]string{
		"image": stegoImagePath,
	}

	req, err := newFileUploadRequest("/api/extract", params, fileParams)
	assert.NoError(t, err)

	rr := serveHTTP(req)
	assert.Equal(t, http.StatusOK, rr.Code) // Should still succeed

	var response utils.SuccessResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	responseData, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	
	assert.False(t, responseData["isFile"].(bool))
	assert.Equal(t, secretMessage, responseData["message"].(string))
}


// TestExtractAPI_EncryptedFileContent_CorrectPassword tests extracting an encrypted file.
func TestExtractAPI_EncryptedFileContent_CorrectPassword(t *testing.T) {
    originalCoverPath, cleanupCover := createTestImageFile(t, "cover_extract_enc_file.png")
    defer cleanupCover()

    fileContent := "This is secret file data for extraction."
    originalFileName := "secret_data.txt"
    password := "filePassword123"
    stegoImageName := "stego_extract_enc_file.png"

    // Create a dummy file to hide
    fileToHidePath, cleanupFileToHide := createTestTextFile(t, fileContent, originalFileName)
    defer cleanupFileToHide()

    // Create an image with an encrypted file
    stegoImagePath, cleanupStego := createTestStegoFileImage(t, originalCoverPath, fileToHidePath, originalFileName, password, stegoImageName)
    defer cleanupStego()

    // Prepare request for extraction
    params := map[string]string{"password": password}
    fileParams := map[string]string{"image": stegoImagePath}

    req, err := newFileUploadRequest("/api/extract", params, fileParams)
    assert.NoError(t, err)

    rr := serveHTTP(req)
    assert.Equal(t, http.StatusOK, rr.Code)

    var response utils.SuccessResponse
    err = json.Unmarshal(rr.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.True(t, response.Success, "API success flag should be true")

    responseData, ok := response.Data.(map[string]interface{})
    assert.True(t, ok, "Response data is not map[string]interface{}")

    assert.True(t, responseData["isFile"].(bool), "IsFile should be true for file content")
    assert.NotEmpty(t, responseData["fileURL"].(string), "FileURL should not be empty")
    assert.Equal(t, originalFileName, responseData["fileName"].(string), "Extracted filename mismatch")
    
    // Further check: Download the file from fileURL and verify content (more involved, optional for this test scope)
}

// TestExtractAPI_NoStegoContent tests extracting from an image with no hidden content.
func TestExtractAPI_NoStegoContent(t *testing.T) {
    imagePath, cleanupImage := createTestImageFile(t, "plain_image_for_extract.png")
    defer cleanupImage()

    params := map[string]string{} // No password
    fileParams := map[string]string{"image": imagePath}

    req, err := newFileUploadRequest("/api/extract", params, fileParams)
    assert.NoError(t, err)

    rr := serveHTTP(req)
    checkAPIErrorResponse(t, rr, http.StatusNotFound, "No hidden content found in this image.")
}
