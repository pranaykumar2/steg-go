package handlers

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pranaykumar2/steg-go/internal/steganography"
	"github.com/pranaykumar2/steg-go/api/utils"
)

var testRouter *gin.Engine

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	testRouter = gin.Default()
	
	// Setup routes (simplified for testing - ideally load from actual router setup)
	// This is a basic setup. In a real app, you'd import your router setup function.
	if testRouter != nil {
		testRouter.POST("/api/hide", HideText)
		testRouter.POST("/api/hideFile", HideFile)
		testRouter.POST("/api/extract", Extract)
		// Add other routes as needed for testing
	}

	// Create temp directory for test files if it doesn't exist
	// This should align with where utils.SaveUploadedFile saves files.
	// For simplicity, assuming utils.TempDir is "temp" relative to where tests are run.
	// In a real scenario, utils.TempDir should be configurable or use os.TempDir().
	if _, err := os.Stat(utils.TempDir); os.IsNotExist(err) {
		_ = os.MkdirAll(utils.TempDir, 0755)
	}


	code := m.Run()

	// Clean up (optional, as os.TempDir() would handle system temp)
	// If utils.TempDir is custom, you might want to clean it.
	// For this example, we'll leave it to manual cleanup or OS.
	
	os.Exit(code)
}

// createTestImageFile creates a dummy PNG image file for testing uploads.
// It returns the path to the created file and a cleanup function.
func createTestImageFile(t *testing.T, filename string) (string, func()) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 10, 10)) // Small 10x10 image
	filePath := filepath.Join(os.TempDir(), filename) // Use system temp dir for safety

	f, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create test image file: %v", err)
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		t.Fatalf("Failed to encode test image: %v", err)
	}
	f.Close()

	return filePath, func() { os.Remove(filePath) }
}

// newFileUploadRequest creates a new file upload HTTP request with extra fields.
func newFileUploadRequest(uri string, params map[string]string, fileParams map[string]string) (*http.Request, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	for key, filePath := range fileParams {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		part, err := writer.CreateFormFile(key, filepath.Base(filePath))
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return nil, err
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}


// Helper to create a steganographic image for extract tests
func createTestStegoImage(t *testing.T, originalImagePath, message, password, outputStegoFilename string) (string, func()) {
	t.Helper()
	
	encoder, err := steganography.NewEncoder(originalImagePath, password) // Password can be empty for no encryption
	if err != nil {
		t.Fatalf("Failed to create encoder for test stego image: %v", err)
	}

	err = encoder.Hide([]byte(message))
	if err != nil {
		t.Fatalf("Failed to hide message for test stego image: %v", err)
	}
	
	stegoImagePath := filepath.Join(os.TempDir(), outputStegoFilename)
	err = encoder.SaveOutput(stegoImagePath)
	if err != nil {
		t.Fatalf("Failed to save test stego image: %v", err)
	}
	
	return stegoImagePath, func() { os.Remove(stegoImagePath) }
}

// Helper to create a steganographic image with a file for extract tests
func createTestStegoFileImage(t *testing.T, originalCoverPath, fileToHidePath, fileToHideOriginalName, password, outputStegoFilename string) (string, func()) {
	t.Helper()

	encoder, err := steganography.NewEncoder(originalCoverPath, password)
	if err != nil {
		t.Fatalf("Failed to create encoder for test stego file image: %v", err)
	}

	fileHandler := steganography.NewFileHandler()
	fileData, metadata, err := fileHandler.ReadFileContent(fileToHidePath)
	if err != nil {
		t.Fatalf("Failed to read file content for test stego file: %v", err)
	}
	metadata.OriginalName = fileToHideOriginalName // Ensure original name is set

	err = encoder.HideFile(fileData, metadata)
	if err != nil {
		t.Fatalf("Failed to hide file for test stego file image: %v", err)
	}

	stegoImagePath := filepath.Join(os.TempDir(), outputStegoFilename)
	err = encoder.SaveOutput(stegoImagePath)
	if err != nil {
		t.Fatalf("Failed to save test stego file image: %v", err)
	}

	return stegoImagePath, func() { os.Remove(stegoImagePath) }
}


// Helper to create a simple text file for embedding
func createTestTextFile(t *testing.T, content, filename string) (string, func()) {
    t.Helper()
    filePath := filepath.Join(os.TempDir(), filename)
    err := os.WriteFile(filePath, []byte(content), 0644)
    if err != nil {
        t.Fatalf("Failed to create test text file: %v", err)
    }
    return filePath, func() { os.Remove(filePath) }
}

func serveHTTP(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	testRouter.ServeHTTP(rr, req)
	return rr
}

// checkAPIErrorResponse is a helper to check for standard error responses.
func checkAPIErrorResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedStatusCode int, expectedErrorMessageSubstring string) {
	t.Helper()
	if rr.Code != expectedStatusCode {
		t.Errorf("handler returned wrong status code: got %v want %v, body: %s", rr.Code, expectedStatusCode, rr.Body.String())
	}
	if !bytes.Contains(rr.Body.Bytes(), []byte(expectedErrorMessageSubstring)) {
		t.Errorf("handler returned unexpected body: got %s want substring %s", rr.Body.String(), expectedErrorMessageSubstring)
	}
}
