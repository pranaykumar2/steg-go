package steganography

import (
  "encoding/binary"
  "errors"
  "os"
  "path/filepath"
  "strings"
)

// Constants FileModeEnabled, TextModeEnabled, and MetadataSize are now defined in common.go
// const (
//   FileModeEnabled byte = 0x01
//   TextModeEnabled byte = 0x00
//   MetadataSize         = 256 // common.go has 258. Code will use common.go's version.
// )

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
    if metadata.FileExt != "" {
      outputPath = outputPath + metadata.FileExt
    } else if ext := filepath.Ext(metadata.OriginalName); ext != "" {
      outputPath = outputPath + ext
    }
  }

  return os.WriteFile(outputPath, data, 0644)
}

func (fh *FileHandler) SerializeMetadata(metadata *FileMetadata) []byte {
  result := make([]byte, MetadataSize)

  result[0] = FileModeEnabled

  binary.BigEndian.PutUint64(result[1:9], metadata.FileSize)

  originalName := metadata.OriginalName

  nameBytes := []byte(originalName)
  if len(nameBytes) > 127 {
    nameBytes = nameBytes[:127]
  }
  copy(result[9:9+len(nameBytes)], nameBytes)

  result[9+len(nameBytes)] = 0

  extBytes := []byte(metadata.FileExt)
  if len(extBytes) > 10 {
    extBytes = extBytes[:10] // Limit extension length
  }

  result[137] = byte(len(extBytes))

  copy(result[138:138+len(extBytes)], extBytes)

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
  for i := 9; i < 137; i++ {
    if data[i] == 0 {
      filenameEnd = i
      break
    }
  }
  originalName := string(data[9:filenameEnd])

  var fileExt string

  if len(data) > 138 {
    extLen := int(data[137])
    if extLen > 0 && extLen <= 10 && 138+extLen <= len(data) {
      fileExt = string(data[138:138+extLen])
    }
  }

  if fileExt == "" {
    fileExt = filepath.Ext(originalName)
  }

  if fileExt != "" && !strings.HasSuffix(originalName, fileExt) {
    if filepath.Ext(originalName) == "" {
      originalName = originalName + fileExt
    }
  }

  return &FileMetadata{
    OriginalName: originalName,
    FileExt:      fileExt,
    FileSize:     fileSize,
  }, nil
}

func (fh *FileHandler) IsFileSupported(filePath string) (bool, string) {
  ext := strings.ToLower(filepath.Ext(filePath))

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
    ".xlsx": true,
    ".pptx": true,
    ".epub": true,
    ".mp4":  true,
    ".avi":  true,
    ".mov":  true,
    ".svg":  true,
    ".html": true,
    ".css":  true,
    ".js":   true,
    ".py":   true,
    ".java": true,
    ".go":   true,
    ".cpp":  true,
    ".h":    true,
    ".c":    true,
    ".rb":   true,
    ".php":  true,
    ".sql":  true,
    ".md":   true,
    ".rtf":  true,
    ".tar":  true,
    ".gz":   true,
    ".7z":   true,
    ".rar":  true,
    ".ico":  true,
    ".psd":  true,
    ".ai":   true,
    ".customization": true,
  }

  return supportedTypes[ext], ext
}
