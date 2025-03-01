package steganography

import (
  "encoding/binary"
  "fmt"
  "image"
  "image/color"
  "image/png"
  "os"
  "path/filepath"
  _ "image/jpeg"

  "github.com/pranaykumar2/steg-go/pkg/imageprocessing"
)

const (
  headerPattern = "STEG"
  headerSize    = 13
  bitsPerByte   = 8
  formatVersion = byte(1)
)

type Encoder struct {
  processor *imageprocessing.ImageProcessor
  image     image.Image
}

func NewEncoder(imagePath string) (*Encoder, error) {
  processor, err := imageprocessing.NewImageProcessor(imagePath)
  if err != nil {
    return nil, err
  }

  return &Encoder{
    processor: processor,
    image:     processor.GetImage(),
  }, nil
}

func (e *Encoder) Hide(data []byte) error {
  bounds := e.image.Bounds()
  width := bounds.Max.X - bounds.Min.X
  height := bounds.Max.Y - bounds.Min.Y

  totalDataSize := len(headerPattern) + 1 + 8 + len(data)
  requiredBits := totalDataSize * 8
  availableBits := width * height * 3

  if requiredBits > availableBits {
    return fmt.Errorf("image too small, need %d bits but have %d", requiredBits, availableBits)
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

  lengthBytes := make([]byte, 8)
  binary.BigEndian.PutUint64(lengthBytes, uint64(len(data)))
  for i := 0; i < len(lengthBytes); i++ {
    writeByte(output, lengthBytes[i], &bitIndex, width, height)
  }

  for i := 0; i < len(data); i++ {
    writeByte(output, data[i], &bitIndex, width, height)
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
    colorBit := (b >> uint(bit)) & 1
    c := img.RGBAAt(x, y)

    switch *bitIndex % 3 {
    case 0:
      c.R = (c.R & 0xFE) | uint8(colorBit)
    case 1:
      c.G = (c.G & 0xFE) | uint8(colorBit)
    case 2:
      c.B = (c.B & 0xFE) | uint8(colorBit)
    }

    img.SetRGBA(x, y, c)
    *bitIndex++
  }
}