package steganography

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// createTestImage generates a simple PNG image and saves it to the given filePath.
// The image is a 100x100 white square.
func createTestImage(t *testing.T, filePath string, width, height int) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Fill the image with white color
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.White)
		}
	}

	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Failed to create directory for test image: %v", err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create test image file: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}
}

// TestEncodeDecodeText_LSBMatching_NoEncryption tests encoding and decoding text without encryption.
func TestEncodeDecodeText_LSBMatching_NoEncryption(t *testing.T) {
	baseDir := t.TempDir() // Create a temporary directory for test files
	inputImagePath := filepath.Join(baseDir, "input_text_no_enc.png")
	outputImagePath := filepath.Join(baseDir, "output_text_no_enc.png")

	createTestImage(t, inputImagePath, 100, 100)

	originalText := "Hello, LSB Matching without Encryption!"
	originalData := []byte(originalText)

	// Encode
	encoder, err := NewEncoder(inputImagePath)
	if err != nil {
		t.Fatalf("NewEncoder failed: %v", err)
	}
	if err := encoder.Hide(originalData); err != nil {
		t.Fatalf("Encoder.Hide failed: %v", err)
	}
	if err := encoder.SaveOutput(outputImagePath); err != nil {
		t.Fatalf("Encoder.SaveOutput failed: %v", err)
	}

	// Decode
	decoder, err := NewDecoder(outputImagePath)
	if err != nil {
		t.Fatalf("NewDecoder failed: %v", err)
	}
	decodedData, isFile, _, stegoFlags, err := decoder.Extract()
	if err != nil {
		t.Fatalf("Decoder.Extract failed: %v", err)
	}

	// Assertions
	if !bytes.Equal(decodedData, originalData) {
		t.Errorf("Decoded data mismatch. Got %q, want %q", string(decodedData), originalText)
	}
	if isFile {
		t.Errorf("isFile flag should be false for text, got true")
	}
	expectedFlags := FlagLSBMatchingUsed
	if stegoFlags != expectedFlags {
		t.Errorf("StegoFlags mismatch. Got %02x, want %02x", stegoFlags, expectedFlags)
	}
}

// TestEncodeDecodeFile_LSBMatching_NoEncryption tests encoding and decoding a file without encryption.
func TestEncodeDecodeFile_LSBMatching_NoEncryption(t *testing.T) {
	baseDir := t.TempDir()
	inputImagePath := filepath.Join(baseDir, "input_file_no_enc.png")
	outputImagePath := filepath.Join(baseDir, "output_file_no_enc.png")

	createTestImage(t, inputImagePath, 200, 200) // Larger image for file metadata + content

	originalFileData := []byte("This is the content of a test file.")
	originalMetadata := &FileMetadata{
		Filename: "test.txt",
		Filesize: int64(len(originalFileData)),
		Filetype: "text/plain",
	}

	// Encode
	encoder, err := NewEncoder(inputImagePath)
	if err != nil {
		t.Fatalf("NewEncoder failed: %v", err)
	}
	if err := encoder.HideFile(originalFileData, originalMetadata); err != nil {
		t.Fatalf("Encoder.HideFile failed: %v", err)
	}
	if err := encoder.SaveOutput(outputImagePath); err != nil {
		t.Fatalf("Encoder.SaveOutput failed: %v", err)
	}

	// Decode
	decoder, err := NewDecoder(outputImagePath)
	if err != nil {
		t.Fatalf("NewDecoder failed: %v", err)
	}
	decodedFileData, isFile, decodedMetadata, stegoFlags, err := decoder.Extract()
	if err != nil {
		t.Fatalf("Decoder.Extract failed: %v", err)
	}

	// Assertions
	if !bytes.Equal(decodedFileData, originalFileData) {
		t.Errorf("Decoded file data mismatch. Got %q, want %q", string(decodedFileData), string(originalFileData))
	}
	if !isFile {
		t.Errorf("isFile flag should be true for file, got false")
	}
	if decodedMetadata == nil {
		t.Fatalf("Decoded metadata is nil")
	}
	if !reflect.DeepEqual(decodedMetadata, originalMetadata) {
		t.Errorf("Decoded metadata mismatch. Got %+v, want %+v", decodedMetadata, originalMetadata)
	}
	expectedFlags := FlagLSBMatchingUsed
	if stegoFlags != expectedFlags {
		t.Errorf("StegoFlags mismatch. Got %02x, want %02x", stegoFlags, expectedFlags)
	}
}

// TestEncodeDecodeText_LSBMatching_WithEncryption tests encoding and decoding text with encryption.
func TestEncodeDecodeText_LSBMatching_WithEncryption(t *testing.T) {
	baseDir := t.TempDir()
	inputImagePath := filepath.Join(baseDir, "input_text_with_enc.png")
	outputImagePath := filepath.Join(baseDir, "output_text_with_enc.png")
	password := "testPassword123"

	createTestImage(t, inputImagePath, 150, 150) // Slightly larger for encryption overhead

	originalText := "Secure message with LSB Matching & Encryption!"
	originalData := []byte(originalText)

	// Encode
	encoder, err := NewEncoder(inputImagePath, password)
	if err != nil {
		t.Fatalf("NewEncoder with password failed: %v", err)
	}
	if err := encoder.Hide(originalData); err != nil {
		t.Fatalf("Encoder.Hide failed: %v", err)
	}
	if err := encoder.SaveOutput(outputImagePath); err != nil {
		t.Fatalf("Encoder.SaveOutput failed: %v", err)
	}

	// Decode
	decoder, err := NewDecoder(outputImagePath, password)
	if err != nil {
		t.Fatalf("NewDecoder with password failed: %v", err)
	}
	decodedData, isFile, _, stegoFlags, err := decoder.Extract()
	if err != nil {
		t.Fatalf("Decoder.Extract failed: %v", err)
	}

	// Assertions
	if !bytes.Equal(decodedData, originalData) {
		t.Errorf("Decoded data mismatch. Got %q, want %q", string(decodedData), originalText)
	}
	if isFile {
		t.Errorf("isFile flag should be false for text, got true")
	}
	expectedFlags := FlagLSBMatchingUsed | FlagEncryptionEnabled
	if stegoFlags != expectedFlags {
		t.Errorf("StegoFlags mismatch. Got %02x, want %02x", stegoFlags, expectedFlags)
	}
}

// TestEncodeDecodeFile_LSBMatching_WithEncryption tests encoding and decoding a file with encryption.
func TestEncodeDecodeFile_LSBMatching_WithEncryption(t *testing.T) {
	baseDir := t.TempDir()
	inputImagePath := filepath.Join(baseDir, "input_file_with_enc.png")
	outputImagePath := filepath.Join(baseDir, "output_file_with_enc.png")
	password := "superSecretFilePassword"

	createTestImage(t, inputImagePath, 250, 250) // Larger for file + encryption

	originalFileData := []byte("This is a super secret file content that needs encryption.")
	originalMetadata := &FileMetadata{
		Filename: "secret_document.txt",
		Filesize: int64(len(originalFileData)),
		Filetype: "application/octet-stream",
	}

	// Encode
	encoder, err := NewEncoder(inputImagePath, password)
	if err != nil {
		t.Fatalf("NewEncoder with password failed: %v", err)
	}
	if err := encoder.HideFile(originalFileData, originalMetadata); err != nil {
		t.Fatalf("Encoder.HideFile failed: %v", err)
	}
	if err := encoder.SaveOutput(outputImagePath); err != nil {
		t.Fatalf("Encoder.SaveOutput failed: %v", err)
	}

	// Decode
	decoder, err := NewDecoder(outputImagePath, password)
	if err != nil {
		t.Fatalf("NewDecoder with password failed: %v", err)
	}
	decodedFileData, isFile, decodedMetadata, stegoFlags, err := decoder.Extract()
	if err != nil {
		t.Fatalf("Decoder.Extract failed: %v", err)
	}

	// Assertions
	if !bytes.Equal(decodedFileData, originalFileData) {
		t.Errorf("Decoded file data mismatch. Got %q, want %q", string(decodedFileData), string(originalFileData))
	}
	if !isFile {
		t.Errorf("isFile flag should be true for file, got false")
	}
	if decodedMetadata == nil {
		t.Fatalf("Decoded metadata is nil")
	}
	if !reflect.DeepEqual(decodedMetadata, originalMetadata) {
		t.Errorf("Decoded metadata mismatch. Got %+v, want %+v", decodedMetadata, originalMetadata)
	}
	expectedFlags := FlagLSBMatchingUsed | FlagEncryptionEnabled
	if stegoFlags != expectedFlags {
		t.Errorf("StegoFlags mismatch. Got %02x, want %02x", stegoFlags, expectedFlags)
	}
}

// TestDecodeText_WithEncryption_WrongPassword tests decoding encrypted text with a wrong password.
func TestDecodeText_WithEncryption_WrongPassword(t *testing.T) {
	baseDir := t.TempDir()
	inputImagePath := filepath.Join(baseDir, "input_text_wrong_pass.png")
	outputImagePath := filepath.Join(baseDir, "output_text_wrong_pass.png")
	correctPassword := "correctPassword"
	wrongPassword := "wrongPassword"

	createTestImage(t, inputImagePath, 150, 150)
	originalText := "This message should not be decoded."
	originalData := []byte(originalText)

	// Encode with correct password
	encoder, err := NewEncoder(inputImagePath, correctPassword)
	if err != nil {
		t.Fatalf("NewEncoder failed: %v", err)
	}
	if err := encoder.Hide(originalData); err != nil {
		t.Fatalf("Encoder.Hide failed: %v", err)
	}
	if err := encoder.SaveOutput(outputImagePath); err != nil {
		t.Fatalf("Encoder.SaveOutput failed: %v", err)
	}

	// Attempt to decode with wrong password
	decoder, err := NewDecoder(outputImagePath, wrongPassword)
	if err != nil {
		t.Fatalf("NewDecoder with wrong password failed: %v", err) // This itself shouldn't fail
	}
	_, _, _, _, err = decoder.Extract()
	if err == nil {
		t.Errorf("Decoder.Extract should have failed with wrong password, but it succeeded.")
	}
	// Optionally, check for a specific error message if crypto.Decrypt returns a distinct error
	// For example: if !strings.Contains(err.Error(), "cipher: message authentication failed") { ... }
}

// TestDecodeText_WithEncryption_NoPassword tests decoding encrypted text without providing a password.
func TestDecodeText_WithEncryption_NoPassword(t *testing.T) {
	baseDir := t.TempDir()
	inputImagePath := filepath.Join(baseDir, "input_text_no_pass_provided.png")
	outputImagePath := filepath.Join(baseDir, "output_text_no_pass_provided.png")
	password := "aPassword"

	createTestImage(t, inputImagePath, 150, 150)
	originalText := "This message is encrypted."
	originalData := []byte(originalText)

	// Encode with password
	encoder, err := NewEncoder(inputImagePath, password)
	if err != nil {
		t.Fatalf("NewEncoder failed: %v", err)
	}
	if err := encoder.Hide(originalData); err != nil {
		t.Fatalf("Encoder.Hide failed: %v", err)
	}
	if err := encoder.SaveOutput(outputImagePath); err != nil {
		t.Fatalf("Encoder.SaveOutput failed: %v", err)
	}

	// Attempt to decode without password
	decoder, err := NewDecoder(outputImagePath) // No password provided
	if err != nil {
		t.Fatalf("NewDecoder without password failed: %v", err) // This itself shouldn't fail
	}
	_, _, _, _, err = decoder.Extract()
	if err == nil {
		t.Errorf("Decoder.Extract should have failed when no password provided for encrypted data, but it succeeded.")
	} else {
		// Check if the error message is as expected
		expectedErrorMsg := "password required for encrypted data, but not provided to decoder"
		if err.Error() != expectedErrorMsg {
			t.Errorf("Expected error %q, got %q", expectedErrorMsg, err.Error())
		}
	}
}

// TestMain can be used for package-level setup/teardown if needed
// For now, t.TempDir() handles cleanup of test files effectively.
/*
func TestMain(m *testing.M) {
	// setup
	code := m.Run()
	// teardown
	os.Exit(code)
}
*/

// TestEncodeDecode_EmptyText tests encoding and decoding an empty text string.
func TestEncodeDecode_EmptyText(t *testing.T) {
	baseDir := t.TempDir()
	inputImagePath := filepath.Join(baseDir, "input_empty_text.png")
	outputImagePath := filepath.Join(baseDir, "output_empty_text.png")

	createTestImage(t, inputImagePath, 100, 100)
	originalText := ""
	originalData := []byte(originalText)

	// Encode
	encoder, err := NewEncoder(inputImagePath)
	if err != nil {
		t.Fatalf("NewEncoder failed: %v", err)
	}
	if err := encoder.Hide(originalData); err != nil {
		t.Fatalf("Encoder.Hide failed: %v", err)
	}
	if err := encoder.SaveOutput(outputImagePath); err != nil {
		t.Fatalf("Encoder.SaveOutput failed: %v", err)
	}

	// Decode
	decoder, err := NewDecoder(outputImagePath)
	if err != nil {
		t.Fatalf("NewDecoder failed: %v", err)
	}
	decodedData, isFile, _, stegoFlags, err := decoder.Extract()
	if err != nil {
		// The current decoder logic for empty payload:
		// return []byte{}, false, nil, stegoFlags, nil
		// So, no error is expected here for an empty message.
		t.Fatalf("Decoder.Extract failed for empty text: %v", err)
	}

	if !bytes.Equal(decodedData, originalData) {
		t.Errorf("Decoded data mismatch for empty text. Got %q, want %q", string(decodedData), originalText)
	}
	if isFile {
		t.Errorf("isFile flag should be false for empty text, got true")
	}
	expectedFlags := FlagLSBMatchingUsed
	if stegoFlags != expectedFlags {
		t.Errorf("StegoFlags mismatch for empty text. Got %02x, want %02x", stegoFlags, expectedFlags)
	}
}

// TestEncodeDecode_EmptyFileContent tests encoding and decoding an empty file.
func TestEncodeDecode_EmptyFileContent(t *testing.T) {
    baseDir := t.TempDir()
    inputImagePath := filepath.Join(baseDir, "input_empty_file.png")
    outputImagePath := filepath.Join(baseDir, "output_empty_file.png")

    createTestImage(t, inputImagePath, 200, 200) // Image needs to be large enough for metadata

    originalFileData := []byte{} // Empty file content
    originalMetadata := &FileMetadata{
        Filename: "empty.txt",
        Filesize: 0, // Size is 0
        Filetype: "text/plain",
    }

    // Encode
    encoder, err := NewEncoder(inputImagePath)
    if err != nil {
        t.Fatalf("NewEncoder failed: %v", err)
    }
    if err := encoder.HideFile(originalFileData, originalMetadata); err != nil {
        t.Fatalf("Encoder.HideFile for empty file failed: %v", err)
    }
    if err := encoder.SaveOutput(outputImagePath); err != nil {
        t.Fatalf("Encoder.SaveOutput for empty file failed: %v", err)
    }

    // Decode
    decoder, err := NewDecoder(outputImagePath)
    if err != nil {
        t.Fatalf("NewDecoder for empty file failed: %v", err)
    }
    decodedFileData, isFile, decodedMetadata, stegoFlags, err := decoder.Extract()
    if err != nil {
        // Similar to empty text, current decoder should handle this.
        t.Fatalf("Decoder.Extract for empty file failed: %v", err)
    }

    if !bytes.Equal(decodedFileData, originalFileData) {
        t.Errorf("Decoded file data mismatch for empty file. Got %q, want %q", string(decodedFileData), string(originalFileData))
    }
    if !isFile {
        t.Errorf("isFile flag should be true for empty file, got false")
    }
    if decodedMetadata == nil {
        t.Fatalf("Decoded metadata is nil for empty file")
    }
    if !reflect.DeepEqual(decodedMetadata, originalMetadata) {
        t.Errorf("Decoded metadata mismatch for empty file. Got %+v, want %+v", decodedMetadata, originalMetadata)
    }
    expectedFlags := FlagLSBMatchingUsed
    if stegoFlags != expectedFlags {
        t.Errorf("StegoFlags mismatch for empty file. Got %02x, want %02x", stegoFlags, expectedFlags)
    }
}
