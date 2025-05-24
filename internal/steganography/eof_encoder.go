package steganography

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pranaykumar2/steg-go/internal/crypto"
	// Assuming common.go is in the same package or accessible.
	// If common.go is in 'steganography' package, direct reference to FlagEncryptionEnabled is fine.
	// For clarity, if it were a separate package, it'd be:
	// "github.com/pranaykumar2/steg-go/internal/steganography/common"
)

const (
	EOFStegMarker     = "::EOF_STEG_MAGIC::"
	FlagEOFMethodUsed = byte(1 << 3) // Indicates data is appended using EOF method
	FlagEOFMethodFile = byte(1 << 4) // Indicates the appended data is a file (vs text)
	// FlagEncryptionEnabled is defined in common.go (0x01)
)

// EOFEncoder struct can hold configuration or state if needed in the future.
type EOFEncoder struct {
}

// NewEOFEncoder creates a new EOFEncoder.
func NewEOFEncoder() *EOFEncoder {
	return &EOFEncoder{}
}

// writeAppendedBlockToEOF handles the core logic of preparing and appending data to the carrier file.
func (e *EOFEncoder) writeAppendedBlockToEOF(carrierPath string, stegoFlags byte, dataToEmbed []byte, password string) error {
	// Check if carrier file exists
	if _, err := os.Stat(carrierPath); os.IsNotExist(err) {
		return fmt.Errorf("carrier file does not exist: %s", carrierPath)
	}

	var finalPayloadToAppend []byte
	var err error

	if password != "" {
		stegoFlags |= FlagEncryptionEnabled // Use FlagEncryptionEnabled from common.go
		encryptedPayload, salt, nonce, errEnc := crypto.Encrypt(dataToEmbed, password)
		if errEnc != nil {
			return fmt.Errorf("failed to encrypt payload: %w", errEnc)
		}
		// Construct block: flags | salt | nonce | encrypted_payload
		finalPayloadToAppend = append(finalPayloadToAppend, stegoFlags)
		finalPayloadToAppend = append(finalPayloadToAppend, salt...)
		finalPayloadToAppend = append(finalPayloadToAppend, nonce...)
		finalPayloadToAppend = append(finalPayloadToAppend, encryptedPayload...)
	} else {
		// Construct block: flags | raw_payload
		finalPayloadToAppend = append(finalPayloadToAppend, stegoFlags)
		finalPayloadToAppend = append(finalPayloadToAppend, dataToEmbed...)
	}

	// Open file for appending
	file, err := os.OpenFile(carrierPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open carrier file for appending: %w", err)
	}
	defer file.Close()

	// 1. Write length of the (flags + [salt+nonce] + payload) block (8 bytes)
	blockLength := uint64(len(finalPayloadToAppend))
	lenBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(lenBytes, blockLength)
	if _, err := file.Write(lenBytes); err != nil {
		return fmt.Errorf("failed to write block length: %w", err)
	}

	// 2. Write the actual block (flags + [salt+nonce] + payload)
	if _, err := file.Write(finalPayloadToAppend); err != nil {
		return fmt.Errorf("failed to write appended block: %w", err)
	}

	// 3. Write the EOFStegMarker
	if _, err := file.Write([]byte(EOFStegMarker)); err != nil {
		return fmt.Errorf("failed to write EOF marker: %w", err)
	}

	return nil
}

// HideTextInEOF hides a text message by appending it to the carrier file.
func (e *EOFEncoder) HideTextInEOF(carrierPath string, message string, password string) error {
	stegoFlags := FlagEOFMethodUsed // Base flag for EOF text
	
	// For text, originalFileNameLength is 0, and originalFileName is empty.
	originalFileNameLength := byte(0)
	// Construct dataToEmbed: [fileNameLen (1 byte)] + [fileName (0 bytes)] + [actualMessage]
	dataToEmbed := append([]byte{originalFileNameLength}, []byte(message)...)

	return e.writeAppendedBlockToEOF(carrierPath, stegoFlags, dataToEmbed, password)
}

// HideFileInEOF hides a file's content by appending it to the carrier file.
func (e *EOFEncoder) HideFileInEOF(carrierPath string, fileToHidePath string, password string) error {
	// Check if fileToHidePath exists
	if _, err := os.Stat(fileToHidePath); os.IsNotExist(err) {
		return fmt.Errorf("file to hide does not exist: %s", fileToHidePath)
	}

	actualPayload, err := os.ReadFile(fileToHidePath)
	if err != nil {
		return fmt.Errorf("failed to read file to hide: %w", err)
	}

	stegoFlags := FlagEOFMethodUsed | FlagEOFMethodFile // Base flags for EOF file

	originalFileName := filepath.Base(fileToHidePath)
	if len(originalFileName) > 255 {
		// This is a simplistic truncation. Consider if this is the desired behavior.
		originalFileName = originalFileName[:255] 
		// Log or warn about truncation if a logging mechanism exists
	}
	originalFileNameLength := byte(len(originalFileName))
	
	// Construct dataToEmbed: [fileNameLen (1 byte)] + [fileNameBytes] + [actualFilePayload]
	dataToEmbed := append([]byte{originalFileNameLength}, []byte(originalFileName)...)
	dataToEmbed = append(dataToEmbed, actualPayload...)

	return e.writeAppendedBlockToEOF(carrierPath, stegoFlags, dataToEmbed, password)
}
