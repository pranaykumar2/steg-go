package imageprocessing

import (
  "errors"
  "image"
  "image/png"
  "os"
  "path/filepath"
  "strings"
  _ "image/jpeg"
  _ "image/png"
)

type ImageProcessor struct {
  image  image.Image
  format string
}

func NewImageProcessor(imagePath string) (*ImageProcessor, error) {
  file, err := os.Open(imagePath)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  format := strings.ToLower(filepath.Ext(imagePath))
  if format != ".png" && format != ".jpg" && format != ".jpeg" {
    return nil, errors.New("unsupported image format. Use PNG or JPEG")
  }

  img, _, err := image.Decode(file)
  if err != nil {
    return nil, err
  }

  if format == ".jpeg" {
    format = ".jpg"
  }

  return &ImageProcessor{
    image:  img,
    format: strings.TrimPrefix(format, "."),
  }, nil
}

func (p *ImageProcessor) SaveImage(outputPath string) error {
  output, err := os.Create(outputPath)
  if err != nil {
    return err
  }
  defer output.Close()

  return png.Encode(output, p.image)
}

func (p *ImageProcessor) GetImage() image.Image {
  return p.image
}

func Test() string {
  return "Image processing package initialized"
}