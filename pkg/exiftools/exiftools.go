package exiftools

import (
  "fmt"
  "os"
  "path/filepath"
  "strings"
  "time"
  "image"
  _ "image/jpeg"
  _ "image/png"
)

type MetadataInfo struct {
  Filename       string
  FileSize       int64
  FileType       string
  MimeType       string
  ModTime        time.Time
  ImageWidth     int
  ImageHeight    int
  HasEXIF        bool
  PrivacyRisks   []string
  CameraMake     string
  CameraModel    string
  Software       string
  CreationTime   string
  GPSPresent     bool
  Properties     map[string]string
}

func GetImageMetadata(filePath string) (*MetadataInfo, error) {
  fileInfo, err := os.Stat(filePath)
  if err != nil {
    return nil, fmt.Errorf("failed to get file info: %v", err)
  }

  metadata := &MetadataInfo{
    Filename:   fileInfo.Name(),
    FileSize:   fileInfo.Size(),
    FileType:   strings.ToUpper(strings.TrimPrefix(filepath.Ext(filePath), ".")),
    ModTime:    fileInfo.ModTime(),
    Properties: make(map[string]string),
  }

  metadata.Properties["File Modified"] = fileInfo.ModTime().Format("2006-01-02 15:04:05")

  switch strings.ToLower(metadata.FileType) {
  case "jpg", "jpeg":
    metadata.MimeType = "image/jpeg"
    metadata.HasEXIF = true
  case "png":
    metadata.MimeType = "image/png"
  case "gif":
    metadata.MimeType = "image/gif"
  case "bmp":
    metadata.MimeType = "image/bmp"
  }

  file, err := os.Open(filePath)
  if err == nil {
    defer file.Close()

    img, _, err := image.DecodeConfig(file)
    if err == nil {
      metadata.ImageWidth = img.Width
      metadata.ImageHeight = img.Height
      metadata.Properties["Image Width"] = fmt.Sprintf("%d pixels", img.Width)
      metadata.Properties["Image Height"] = fmt.Sprintf("%d pixels", img.Height)
    }

    file.Seek(0, 0)
    data := make([]byte, 12)
    file.Read(data)

    if len(data) >= 2 {
      if data[0] == 0xFF && data[1] == 0xD8 {
        metadata.HasEXIF = true
        metadata.PrivacyRisks = append(metadata.PrivacyRisks,
          "JPEG file may contain EXIF metadata with personal information")
      }

      if len(data) >= 8 &&
        data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
        metadata.Properties["Format"] = "PNG"
      }
    }
  }

  if metadata.ImageWidth > 0 && metadata.ImageHeight > 0 {
    capacity := (metadata.ImageWidth * metadata.ImageHeight * 3) / 8
    if capacity > 1024 {
      metadata.Properties["Steganography Capacity"] = fmt.Sprintf("~%.2f KB", float64(capacity)/1024)
    } else {
      metadata.Properties["Steganography Capacity"] = fmt.Sprintf("~%d bytes", capacity)
    }
  }
  metadata.analyzePrivacyRisks()

  return metadata, nil
}

func (m *MetadataInfo) analyzePrivacyRisks() {
  if len(m.PrivacyRisks) > 0 {
    return
  }

  if m.FileType == "JPG" || m.FileType == "JPEG" {
    m.PrivacyRisks = append(m.PrivacyRisks,
      "JPEG files commonly contain EXIF data with camera info, location, and creation time")
  }

  if m.GPSPresent {
    m.PrivacyRisks = append(m.PrivacyRisks,
      "Location data detected - exact coordinates may be embedded")
  }

  if m.CameraMake != "" || m.CameraModel != "" {
    m.PrivacyRisks = append(m.PrivacyRisks,
      "Device information detected - can link image to specific camera")
  }

  m.PrivacyRisks = append(m.PrivacyRisks,
    "Steganography preserves most metadata - consider removing sensitive metadata")
}
