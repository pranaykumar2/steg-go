package steganography

import (
  "encoding/binary"
  "errors"
  "image"
  "os"
  _ "image/jpeg"
  _ "image/png"

  "github.com/pranaykumar2/steg-go/internal/crypto"
)

// Constants from encoder.go, or a shared constants package, would be ideal.
// For this exercise, we'll define them here as specified by the task.
const (
	// Steganography feature flags (bitmask)
	FlagEncryptionEnabled = byte(1 << 0) // 0x01
	FlagLSBMatchingUsed   = byte(1 << 1) // 0x02

	// Crypto constants
	SaltSize  = 16
	NonceSize = 12
	// formatVersion and headerPattern are implicitly available from the package context
)

type Decoder struct {
  image       image.Image
  fileHandler *FileHandler
  password    string
}

func NewDecoder(imagePath string, password ...string) (*Decoder, error) {
  file, err := os.Open(imagePath)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  img, _, err := image.Decode(file)
  if err != nil {
    return nil, err
  }

  dec := &Decoder{
    image:       img,
    fileHandler: NewFileHandler(),
  }

  if len(password) > 0 && password[0] != "" {
    dec.password = password[0]
  }

  return dec, nil
}

func (d *Decoder) Extract() ([]byte, bool, *FileMetadata, byte, error) {
  bounds := d.image.Bounds()
  width := bounds.Max.X - bounds.Min.X
  height := bounds.Max.Y - bounds.Min.Y

  headerBytes := make([]byte, len(headerPattern))
  bitIndex := 0
  for i := 0; i < len(headerPattern); i++ {
    headerBytes[i] = readByte(d.image, &bitIndex, width, height)
  }

  if string(headerBytes) != headerPattern {
    return nil, false, nil, 0, errors.New("no steganographic data found")
  }

  version := readByte(d.image, &bitIndex, width, height)
  if version != formatVersion {
    return nil, false, nil, 0, errors.New("unsupported steganography format version")
  }

  lengthBytes := make([]byte, 8)
  for i := 0; i < 8; i++ {
    lengthBytes[i] = readByte(d.image, &bitIndex, width, height)
  }

  dataLength := binary.BigEndian.Uint64(lengthBytes) // This is the total length of data following this length field
  if dataLength == 0 || dataLength > uint64((width*height*3)/8) { // Basic sanity check
    return nil, false, nil, 0, errors.New("invalid data length")
  }

  stegoFlags := readByte(d.image, &bitIndex, width, height)
  remainingDataLength := dataLength - 1 // Subtract 1 for the stegoFlags byte itself

  isEncrypted := (stegoFlags & FlagEncryptionEnabled) != 0
  _ = (stegoFlags & FlagLSBMatchingUsed) != 0 // Store this for potential future use or logging; assign to blank identifier if not used immediately

  // For debugging or logging, one might print lsbMatchingUsed.
  // For now, it doesn't change the decoding logic for LSBs.

  var salt, nonce []byte
  if isEncrypted {
    if d.password == "" {
      return nil, false, nil, stegoFlags, errors.New("password required for encrypted data, but not provided to decoder")
    }
    if remainingDataLength < SaltSize+NonceSize {
      return nil, false, nil, stegoFlags, errors.New("invalid data length for encrypted content (salt/nonce missing)")
    }
    salt = make([]byte, SaltSize)
    for i := 0; i < SaltSize; i++ {
      salt[i] = readByte(d.image, &bitIndex, width, height)
    }
    nonce = make([]byte, NonceSize)
    for i := 0; i < NonceSize; i++ {
      nonce[i] = readByte(d.image, &bitIndex, width, height)
    }
    remainingDataLength -= (SaltSize + NonceSize)
  }
  // No explicit error for unknown flags in stegoFlags, allowing for future additions.
  // We only care about the flags relevant to this version of the decoder.

  // Validate remainingDataLength
  // If encrypted, remainingDataLength can be 0 (empty encrypted message).
  // If not encrypted, remainingDataLength must be at least 1 for the mode indicator.
  if remainingDataLength < 0 || (!isEncrypted && remainingDataLength < 1) {
    return nil, false, nil, stegoFlags, errors.New("invalid remaining data length after processing headers")
  }


  // The 'data' read here is what's left after stegoFlags, salt, nonce.
  // This 'data' will contain the modeIndicator/FileMetadata and the actual (possibly encrypted) content.
  extractedDataPortion := make([]byte, remainingDataLength)
  for i := uint64(0); i < remainingDataLength; i++ {
    extractedDataPortion[i] = readByte(d.image, &bitIndex, width, height)
  }

  var finalPayload []byte
  var err error

  if isEncrypted {
    // If extractedDataPortion is empty and encryption is on, Decrypt might handle it or error.
    // crypto.Decrypt should ideally handle an empty ciphertext if that's a valid state.
    finalPayload, err = crypto.Decrypt(extractedDataPortion, d.password, salt, nonce)
    if err != nil {
      return nil, false, nil, stegoFlags, fmt.Errorf("failed to decrypt data: %w", err)
    }
  } else {
    finalPayload = extractedDataPortion
  }

  // Now, process the finalPayload which is either decrypted data or the original non-encrypted data.
  // This payload still contains the mode indicator or file metadata header.
  if len(finalPayload) == 0 {
      // This could happen if an empty message was hidden, or empty file.
      // For text mode, it might be an issue if TextModeEnabled byte is expected.
      // For file mode, it implies no metadata and no content.
      // Depending on requirements, this might be an error or a valid state.
      // For now, let's assume an empty payload means no further processing of mode/metadata.
      return []byte{}, false, nil, stegoFlags, nil // Return empty data, not a file, no metadata
  }


  modeIndicator := finalPayload[0] // First byte of the (decrypted or original) payload
  isFile := modeIndicator == FileModeEnabled

  var appMetadata *FileMetadata // Renamed from 'metadata' to avoid confusion with crypto metadata
  var contentData []byte

  if isFile {
    // The finalPayload for a file contains: [FileModeEnabled (1 byte)][SerializedFileMetadata (MetadataSize bytes)][ActualFileContent]
    // However, the `FileModeEnabled` byte itself is not part of the `SerializedFileMetadata`.
    // The `fileHandler.DeserializeMetadata` expects only the `MetadataSize` bytes.
    
    // We need to ensure finalPayload is long enough for modeIndicator + MetadataSize
    if len(finalPayload) < 1+MetadataSize {
        return nil, false, nil, stegoFlags, errors.New("invalid file data: too small for metadata")
    }
    
    appMetadata, err = d.fileHandler.DeserializeMetadata(finalPayload[1 : 1+MetadataSize])
    if err != nil {
      return nil, false, nil, stegoFlags, fmt.Errorf("failed to deserialize file metadata: %w", err)
    }
    contentData = finalPayload[1+MetadataSize:]
  } else {
    // For text mode, finalPayload contains: [TextModeEnabled (1 byte)][ActualTextContent]
    if len(finalPayload) < 1 { // Should always have at least the mode indicator
        return nil, false, nil, stegoFlags, errors.New("invalid text data: missing mode indicator")
    }
    contentData = finalPayload[1:]
  }

  return contentData, isFile, appMetadata, stegoFlags, nil
}

func readByte(img image.Image, bitIndex *int, width, height int) byte {
  var b byte
  for bit := 7; bit >= 0; bit-- {
    x := *bitIndex / (height * 3)
    y := (*bitIndex / 3) % height

    if x >= width {
      return 0
    }

    r, g, b_, _ := img.At(x, y).RGBA()
    var colorBit uint8

    switch *bitIndex % 3 {
    case 0:
      colorBit = uint8(r & 1)
    case 1:
      colorBit = uint8(g & 1)
    case 2:
      colorBit = uint8(b_ & 1)
    }

    b |= colorBit << uint(bit)
    *bitIndex++
  }
  return b
}
