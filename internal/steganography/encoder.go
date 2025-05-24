package steganography

import (
  "encoding/binary"
  "fmt"
  "image"
  "image/color"
  "image/png"
  "math/rand"
  "os"
  "path/filepath"
  "time"
  _ "image/jpeg"

  "github.com/pranaykumar2/steg-go/internal/crypto"
  "github.com/pranaykumar2/steg-go/pkg/imageprocessing"
)

const (
  headerPattern = "STEG"
  headerSize    = 13
  bitsPerByte   = 8
  formatVersion = byte(1)

  // Steganography feature flags (bitmask) are now in common.go
  // FlagEncryptionEnabled = byte(1 << 0) // 0x01
  // FlagLSBMatchingUsed   = byte(1 << 1) // 0x02

  // Crypto constants (SaltSize and NonceSize) are now in common.go
  // SaltSize  = 16
  // NonceSize = 12
)

// LSB specific constants (headerPattern, formatVersion) remain here.
// TextModeEnabled and FileModeEnabled are in common.go

type Encoder struct {
  processor   *imageprocessing.ImageProcessor
  image       image.Image
  fileHandler *FileHandler
  password    string
}

func NewEncoder(imagePath string, password ...string) (*Encoder, error) {
  processor, err := imageprocessing.NewImageProcessor(imagePath)
  if err != nil {
    return nil, err
  }

  rand.Seed(time.Now().UnixNano())

  enc := &Encoder{
    processor:   processor,
    image:       processor.GetImage(),
    fileHandler: NewFileHandler(),
  }

  if len(password) > 0 && password[0] != "" {
    enc.password = password[0]
  }

  return enc, nil
}

func (e *Encoder) Hide(data []byte) error {
  bounds := e.image.Bounds()
  width := bounds.Max.X - bounds.Min.X
  height := bounds.Max.Y - bounds.Min.Y

  payload := data
  stegoFlags := FlagLSBMatchingUsed // LSB Matching is always used by this encoder
  var salt, nonce []byte
  var err error

  if e.password != "" {
    payload, salt, nonce, err = crypto.Encrypt(data, e.password)
    if err != nil {
      return fmt.Errorf("failed to encrypt data: %w", err)
    }
    stegoFlags |= FlagEncryptionEnabled
  }

  metadata := []byte{TextModeEnabled} // This is the app-level metadata (mode indicator)

  // Calculate totalDataSize (total bits required for the image)
  // This includes header, version, length field, stegoFlags, optional salt/nonce, app metadata, and payload.
  totalHeaderSize := len(headerPattern) + 1 /*version*/ + 8 /*length field*/ + 1 /*stegoFlags*/
  dataPortionSize := len(metadata) + len(payload)
  if (stegoFlags & FlagEncryptionEnabled) != 0 {
    dataPortionSize += SaltSize + NonceSize
  }
  totalDataSize := totalHeaderSize + dataPortionSize

  requiredBits := totalDataSize * 8
  availableBits := width * height * 3

  if requiredBits > availableBits {
    return fmt.Errorf("image too small, need %d bits but have %d (data len: %d, payload len: %d, salt: %d, nonce: %d, meta: %d)",
			requiredBits, availableBits, len(data), len(payload), len(salt), len(nonce), len(metadata))
  }

  output := image.NewRGBA(bounds)

  for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for x := bounds.Min.X; x < bounds.Max.X; x++ {
      r, g, b, a := e.image.At(x, y).RGBA()
      output.Set(x, y, color.RGBA{
        R: uint8(r >> 8),
        G: uint8(g >> 8),
        B: uint8(b >> 8),
        A: uint8(a >> 8),
      })
    }
  }

  bitIndex := 0

  for i := 0; i < len(headerPattern); i++ {
    b := headerPattern[i]
    writeByte(output, b, &bitIndex, width, height)
  }

  writeByte(output, formatVersion, &bitIndex, width, height)

  // Write total embedded data length.
  // This length is for everything *after* the length field itself:
  // stegoFlags (1) + optional salt & nonce + app metadata + payload.
  actualEmbeddedLength := 1 /*stegoFlags*/
  if (stegoFlags & FlagEncryptionEnabled) != 0 {
    actualEmbeddedLength += SaltSize + NonceSize
  }
  actualEmbeddedLength += len(metadata) + len(payload)

  lengthBytes := make([]byte, 8)
  binary.BigEndian.PutUint64(lengthBytes, uint64(actualEmbeddedLength))
  for i := 0; i < len(lengthBytes); i++ {
    writeByte(output, lengthBytes[i], &bitIndex, width, height)
  }

  writeByte(output, stegoFlags, &bitIndex, width, height)

  if (stegoFlags & FlagEncryptionEnabled) != 0 {
    for i := 0; i < SaltSize; i++ {
      writeByte(output, salt[i], &bitIndex, width, height)
    }
    for i := 0; i < NonceSize; i++ {
      writeByte(output, nonce[i], &bitIndex, width, height)
    }
  }

  for i := 0; i < len(metadata); i++ { // This writes the TextModeEnabled byte
    writeByte(output, metadata[i], &bitIndex, width, height)
  }

  for i := 0; i < len(payload); i++ { // This writes the (potentially encrypted) data
    writeByte(output, payload[i], &bitIndex, width, height)
  }

  e.processor = &imageprocessing.ImageProcessor{}
  e.image = output
  return nil
}

func (e *Encoder) HideFile(fileData []byte, metadata *FileMetadata) error {
  bounds := e.image.Bounds()
  width := bounds.Max.X - bounds.Min.X
  height := bounds.Max.Y - bounds.Min.Y

  payload := fileData
  stegoFlags := FlagLSBMatchingUsed // LSB Matching is always used
  var salt, nonce []byte
  var err error

  if e.password != "" {
    payload, salt, nonce, err = crypto.Encrypt(fileData, e.password)
    if err != nil {
      return fmt.Errorf("failed to encrypt file data: %w", err)
    }
    stegoFlags |= FlagEncryptionEnabled
  }

  metadataBytes := e.fileHandler.SerializeMetadata(metadata) // This is app-level metadata

  // Calculate totalDataSize (total bits required for the image)
  // This includes header, version, length field, stegoFlags, optional salt/nonce, app metadata, and payload.
  totalHeaderSize := len(headerPattern) + 1 /*version*/ + 8 /*length field*/ + 1 /*stegoFlags*/
  dataPortionSize := len(metadataBytes) + len(payload)
  if (stegoFlags & FlagEncryptionEnabled) != 0 {
    dataPortionSize += SaltSize + NonceSize
  }
  totalDataSize := totalHeaderSize + dataPortionSize

  requiredBits := totalDataSize * 8
  availableBits := width * height * 3

  if requiredBits > availableBits {
    return fmt.Errorf("image too small, need %d bits but have %d (file data len: %d, payload len: %d, salt: %d, nonce: %d, meta: %d)",
			requiredBits, availableBits, len(fileData), len(payload), len(salt), len(nonce), len(metadataBytes))
  }

  output := image.NewRGBA(bounds)

  for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
    for x := bounds.Min.X; x < bounds.Max.X; x++ {
      r, g, b, a := e.image.At(x, y).RGBA()
      output.Set(x, y, color.RGBA{
        R: uint8(r >> 8),
        G: uint8(g >> 8),
        B: uint8(b >> 8),
        A: uint8(a >> 8),
      })
    }
  }

  bitIndex := 0

  for i := 0; i < len(headerPattern); i++ {
    b := headerPattern[i]
    writeByte(output, b, &bitIndex, width, height)
  }

  writeByte(output, formatVersion, &bitIndex, width, height)

  // Write total embedded data length.
  // This length is for everything *after* the length field itself:
  // stegoFlags (1) + optional salt & nonce + app metadata + payload.
  actualEmbeddedLength := 1 /*stegoFlags*/
  if (stegoFlags & FlagEncryptionEnabled) != 0 {
    actualEmbeddedLength += SaltSize + NonceSize
  }
  actualEmbeddedLength += len(metadataBytes) + len(payload)

  lengthBytes := make([]byte, 8)
  binary.BigEndian.PutUint64(lengthBytes, uint64(actualEmbeddedLength))
  for i := 0; i < len(lengthBytes); i++ {
    writeByte(output, lengthBytes[i], &bitIndex, width, height)
  }

  writeByte(output, stegoFlags, &bitIndex, width, height)

  if (stegoFlags & FlagEncryptionEnabled) != 0 {
    for i := 0; i < SaltSize; i++ {
      writeByte(output, salt[i], &bitIndex, width, height)
    }
    for i := 0; i < NonceSize; i++ {
      writeByte(output, nonce[i], &bitIndex, width, height)
    }
  }

  for i := 0; i < len(metadataBytes); i++ {
    writeByte(output, metadataBytes[i], &bitIndex, width, height)
  }

  for i := 0; i < len(payload); i++ { // This writes the (potentially encrypted) file data
    writeByte(output, payload[i], &bitIndex, width, height)
  }

  e.processor = &imageprocessing.ImageProcessor{}
  e.image = output
  return nil
}

func (e *Encoder) SaveOutput(outputPath string) error {
  if filepath.Ext(outputPath) == "" {
    outputPath += ".png"
  }

  output, err := os.Create(outputPath)
  if err != nil {
    return err
  }
  defer output.Close()

  return png.Encode(output, e.image)
}

func writeByte(img *image.RGBA, b byte, bitIndex *int, width, height int) {
  for bit := 7; bit >= 0; bit-- {
    x := *bitIndex / (height * 3)
    y := (*bitIndex / 3) % height
    bitToEmbed := (b >> uint(bit)) & 1
    c := img.RGBAAt(x, y)

    processChannel := func(channelValue uint8) uint8 {
      if channelValue%2 == bitToEmbed {
        return channelValue // LSB already matches, do nothing
      }

      // LSB needs to be flipped
      if channelValue == 0 {
        return 1 // Can only change to 1
      } else if channelValue == 255 {
        return 254 // Can only change to 254
      } else {
        // Randomly add or subtract 1
        if rand.Intn(2) == 0 { // Subtract 1
          return channelValue - 1
        } else { // Add 1
          return channelValue + 1
        }
      }
    }

    switch *bitIndex % 3 {
    case 0:
      c.R = processChannel(c.R)
    case 1:
      c.G = processChannel(c.G)
    case 2:
      c.B = processChannel(c.B)
    }

    img.SetRGBA(x, y, c)
    *bitIndex++
  }
}
