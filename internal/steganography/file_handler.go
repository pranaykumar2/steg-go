package steganography

import (
  "encoding/binary"
  "errors"
  "os"
  "path/filepath"
  "strings"
)

const (
  FileModeEnabled byte = 0x01
  TextModeEnabled byte = 0x00
  MetadataSize         = 256
)

type FileMetadata struct {
  OriginalName string
  FileExt      string
  FileSize     uint64
}

type FileHandler struct{}

func NewFileHandler() *FileHandler {
  return &FileHandler{}
}

func (fh *FileHandler) ReadFileContent(filePath string) ([]byte, *FileMetadata, error) {
  data, err := os.ReadFile(filePath)
  if err != nil {
    return nil, nil, err
  }

  fileName := filepath.Base(filePath)
  fileExt := filepath.Ext(filePath)

  metadata := &FileMetadata{
    OriginalName: fileName,
    FileExt:      fileExt,
    FileSize:     uint64(len(data)),
  }

  return data, metadata, nil
}

func (fh *FileHandler) SaveFileContent(data []byte, metadata *FileMetadata, outputPath string) error {
  fileInfo, err := os.Stat(outputPath)
  if err == nil && fileInfo.IsDir() {
    outputPath = filepath.Join(outputPath, metadata.OriginalName)
  } else if filepath.Ext(outputPath) == "" {
    outputPath = outputPath + metadata.FileExt
  }

  return os.WriteFile(outputPath, data, 0644)
}

func (fh *FileHandler) SerializeMetadata(metadata *FileMetadata) []byte {
  result := make([]byte, MetadataSize)

  result[0] = FileModeEnabled

  binary.BigEndian.PutUint64(result[1:9], metadata.FileSize)

  nameBytes := []byte(metadata.OriginalName)
  if len(nameBytes) > 127 {
    nameBytes = nameBytes[:127]
  }
  copy(result[9:9+len(nameBytes)], nameBytes)

  result[9+len(nameBytes)] = 0

  return result
}

func (fh *FileHandler) DeserializeMetadata(data []byte) (*FileMetadata, error) {
  if len(data) < MetadataSize {
    return nil, errors.New("invalid metadata size")
  }

  if data[0] != FileModeEnabled {
    return nil, errors.New("not in file mode")
  }

  fileSize := binary.BigEndian.Uint64(data[1:9])

  filenameEnd := 9
  for i := 9; i < len(data); i++ {
    if data[i] == 0 {
      filenameEnd = i
      break
    }
  }

  originalName := string(data[9:filenameEnd])
  fileExt := filepath.Ext(originalName)

  return &FileMetadata{
    OriginalName: originalName,
    FileExt:      fileExt,
    FileSize:     fileSize,
  }, nil
}

func (fh *FileHandler) IsFileSupported(filePath string) (bool, string) {
  ext := strings.ToLower(filepath.Ext(filePath))

  // List of commonly supported file types
  supportedTypes := map[string]bool{
    ".pdf":  true,
    ".doc":  true,
    ".docx": true,
    ".txt":  true,
    ".mp3":  true,
    ".wav":  true,
    ".ogg":  true,
    ".jpg":  true,
    ".jpeg": true,
    ".png":  true,
    ".gif":  true,
    ".zip":  true,
    ".json": true,
    ".xml":  true,
    ".csv":  true,
  }

  return supportedTypes[ext], ext
}
