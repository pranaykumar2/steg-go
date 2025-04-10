package steganography

import (
  "encoding/binary"
  "errors"
  "image"
  "os"
  _ "image/jpeg"
  _ "image/png"
)

type Decoder struct {
  image       image.Image
  fileHandler *FileHandler
}

func NewDecoder(imagePath string) (*Decoder, error) {
  file, err := os.Open(imagePath)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  img, _, err := image.Decode(file)
  if err != nil {
    return nil, err
  }

  return &Decoder{
    image:       img,
    fileHandler: NewFileHandler(),
  }, nil
}

func (d *Decoder) Extract() ([]byte, bool, *FileMetadata, error) {
  bounds := d.image.Bounds()
  width := bounds.Max.X - bounds.Min.X
  height := bounds.Max.Y - bounds.Min.Y

  headerBytes := make([]byte, len(headerPattern))
  bitIndex := 0
  for i := 0; i < len(headerPattern); i++ {
    headerBytes[i] = readByte(d.image, &bitIndex, width, height)
  }

  if string(headerBytes) != headerPattern {
    return nil, false, nil, errors.New("no steganographic data found")
  }

  version := readByte(d.image, &bitIndex, width, height)
  if version != formatVersion {
    return nil, false, nil, errors.New("unsupported steganography format version")
  }

  lengthBytes := make([]byte, 8)
  for i := 0; i < 8; i++ {
    lengthBytes[i] = readByte(d.image, &bitIndex, width, height)
  }

  dataLength := binary.BigEndian.Uint64(lengthBytes)
  if dataLength == 0 || dataLength > uint64((width*height*3)/8) {
    return nil, false, nil, errors.New("invalid data length")
  }

  modeIndicator := readByte(d.image, &bitIndex, width, height)

  bitIndex -= 8

  data := make([]byte, dataLength)
  for i := uint64(0); i < dataLength; i++ {
    data[i] = readByte(d.image, &bitIndex, width, height)
  }

  isFile := modeIndicator == FileModeEnabled

  var metadata *FileMetadata
  var contentData []byte

  if isFile {
    if len(data) <= MetadataSize {
      return nil, false, nil, errors.New("invalid file data: too small")
    }

    var err error
    metadata, err = d.fileHandler.DeserializeMetadata(data[:MetadataSize])
    if err != nil {
      return nil, false, nil, err
    }

    contentData = data[MetadataSize:]
  } else {
    contentData = data[1:]
  }

  return contentData, isFile, metadata, nil
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
