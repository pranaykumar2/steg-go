package main

import (
  "encoding/hex"
  "fmt"
  "os"
  "os/user"
  "strings"
  _ "image/jpeg"
  _ "image/png"
  "github.com/fatih/color"
  "github.com/pranaykumar2/steg-go/internal/crypto"
  "github.com/pranaykumar2/steg-go/internal/steganography"
  "github.com/pranaykumar2/steg-go/internal/ui"
  "github.com/pranaykumar2/steg-go/pkg/exiftools"
)

const (
  appVersion = "1.1.0"
)

func main() {
  currentUser, err := user.Current()
  username := "user"
  if err == nil && currentUser.Username != "" {
    username = currentUser.Username
  }
  userInterface := ui.NewUI(username)
  userInterface.PrintHeader()

  if len(os.Args) < 2 {
    userInterface.ShowError("No command specified")
    printUsage(userInterface)
    os.Exit(1)
  }

  switch os.Args[1] {
  case "hide":
    if err := handleHideCommand(userInterface); err != nil {
      userInterface.ShowError(fmt.Sprintf("%v", err))
      os.Exit(1)
    }
  case "hideFile":
    if err := handleHideFileCommand(userInterface); err != nil {
      userInterface.ShowError(fmt.Sprintf("%v", err))
      os.Exit(1)
    }
  case "extract":
    if err := handleExtractCommand(userInterface); err != nil {
      userInterface.ShowError(fmt.Sprintf("%v", err))
      os.Exit(1)
    }
    case "metadata":
    if err := handleMetadataCommand(userInterface); err != nil {
      userInterface.ShowError(fmt.Sprintf("%v", err))
      os.Exit(1)
    }
  case "info":
    showInfo(userInterface)
  case "test":
    if len(os.Args) > 2 {
      testWithRealImage(os.Args[2])
    } else {
      userInterface.ShowError("Test command requires an image path")
    }
  default:
    userInterface.ShowError(fmt.Sprintf("Invalid command: %s", os.Args[1]))
    printUsage(userInterface)
    os.Exit(1)
  }
}

func printUsage(ui *ui.UI) {
  ui.PrintCommandHeader("USAGE INFORMATION")

  fmt.Printf("  %s <command> [options]\n\n", os.Args[0])

  ui.PrintFeatureList("Available Commands", []string{
    "hide        Hide a secret message in an image",
    "hideFile    Hide a file (PDF, document, audio, etc.) in an image",
    "extract     Extract hidden content from an image",
    "metadata    Display detailed metadata from an image",
    "info        Show information about this application",
  })

  ui.PrintFeatureList("Examples", []string{
    fmt.Sprintf("%s hide", os.Args[0]),
    fmt.Sprintf("%s hideFile", os.Args[0]),
    fmt.Sprintf("%s extract", os.Args[0]),
    fmt.Sprintf("%s metadata", os.Args[0]),
  })
}

func handleMetadataCommand(ui *ui.UI) error {
  ui.PrintCommandHeader("IMAGE METADATA ANALYSIS")

  imagePath := ui.PromptInput("Enter image path to analyze")
  if !fileExists(imagePath) {
    return fmt.Errorf("file does not exist: %s", imagePath)
  }

  ui.StartProgress("Analyzing image metadata")

  metadata, err := exiftools.GetImageMetadata(imagePath)
  if err != nil {
    ui.StopProgress()
    return fmt.Errorf("failed to extract metadata: %v", err)
  }

  ui.StopProgress()

  ui.PrintMetadata(metadata)

  if metadata.ImageWidth > 0 && metadata.ImageHeight > 0 {
    analyzeCapacity := ui.PromptConfirmation("Would you like a detailed steganography capacity analysis?")
    if analyzeCapacity {
      fmt.Println()
      color.New(color.FgHiCyan).Println("  ┌─ Steganography Capacity Analysis ─────────────┐")

      totalPixels := metadata.ImageWidth * metadata.ImageHeight

      color.New(color.FgCyan).Printf("  │ • Image dimensions: %d × %d pixels            │\n",
        metadata.ImageWidth, metadata.ImageHeight)
      color.New(color.FgCyan).Printf("  │ • Total pixels: %-33d │\n", totalPixels)

      lsbBytes := (totalPixels * 3) / 8
      color.New(color.FgCyan).Printf("  │ • LSB capacity (1-bit): %-23s │\n", formatBytes(lsbBytes))
      textChars := int(float64(lsbBytes) * 8 / 5.1)
      color.New(color.FgCyan).Printf("  │ • Estimated text capacity: ~%-20d │\n", textChars)
      color.New(color.FgCyan).Printf("  │   (characters)                               │\n")
      color.New(color.FgCyan).Println("  │                                                 │")
      color.New(color.FgCyan).Println("  │ This image could potentially store:             │")

      pdfPages := lsbBytes / 100000
      if pdfPages < 1 {
        color.New(color.FgCyan).Println("  │ • A small portion of a PDF document             │")
      } else {
        color.New(color.FgCyan).Printf("  │ • A PDF document of ~%d pages                  │\n", pdfPages)
      }

      mp3Seconds := lsbBytes / 16000
      if mp3Seconds < 1 {
        color.New(color.FgCyan).Println("  │ • Less than a second of MP3 audio               │")
      } else if mp3Seconds < 60 {
        color.New(color.FgCyan).Printf("  │ • ~%d seconds of MP3 audio                     │\n", mp3Seconds)
      } else {
        color.New(color.FgCyan).Printf("  │ • ~%.1f minutes of MP3 audio                   │\n", float64(mp3Seconds)/60)
      }

      if lsbBytes > 50000 {
        color.New(color.FgCyan).Println("  │ • A small JPG image                             │")
      }

      color.New(color.FgHiCyan).Println("  └─────────────────────────────────────────────┘")

      fmt.Println()
      color.New(color.FgYellow).Println("  Note: Actual capacity may be lower due to overhead")
      color.New(color.FgYellow).Println("  and limitations of the steganography algorithm.")
    }
  }

  return nil
}

func formatBytes(bytes int) string {
  if bytes >= 1048576 {
    return fmt.Sprintf("%.2f MB", float64(bytes)/1048576)
  } else if bytes >= 1024 {
    return fmt.Sprintf("%.2f KB", float64(bytes)/1024)
  } else {
    return fmt.Sprintf("%d bytes", bytes)
  }
}


func handleHideCommand(ui *ui.UI) error {
  ui.PrintCommandHeader("HIDE TEXT MESSAGE")

  inputPath := ui.PromptInput("Enter input image path (PNG or JPG)")
  if !fileExists(inputPath) {
    return fmt.Errorf("input file does not exist: %s", inputPath)
  }

  outputPath := ui.PromptInput("Enter output image path (will be saved as PNG)")
  if !strings.HasSuffix(strings.ToLower(outputPath), ".png") {
    outputPath += ".png"
  }

  message := ui.PromptInput("Enter the secret message")
  message = strings.TrimSpace(message)
  if message == "" {
    return fmt.Errorf("message cannot be empty")
  }

  ui.StartProgress("Initializing encoder")
  encoder, err := steganography.NewEncoder(inputPath)
  if err != nil {
    ui.StopProgress()
    return fmt.Errorf("failed to initialize encoder: %v", err)
  }

  ui.UpdateProgress("Generating encryption key")
  encryptor, err := crypto.NewEncryptor()
  if err != nil {
    ui.StopProgress()
    return fmt.Errorf("failed to initialize encryption: %v", err)
  }

  ui.UpdateProgress("Encrypting message")
  encrypted, err := encryptor.Encrypt([]byte(message))
  if err != nil {
    ui.StopProgress()
    return fmt.Errorf("failed to encrypt message: %v", err)
  }

  ui.UpdateProgress("Hiding message in image")
  if err := encoder.Hide(encrypted); err != nil {
    ui.StopProgress()
    return fmt.Errorf("failed to hide message: %v", err)
  }

  ui.UpdateProgress("Saving output image")
  if err := encoder.SaveOutput(outputPath); err != nil {
    ui.StopProgress()
    return fmt.Errorf("failed to save output image: %v", err)
  }
  ui.StopProgress()

  details := map[string]string{
    "Input Image": inputPath,
    "Output Image": outputPath,
    "Message Length": fmt.Sprintf("%d characters", len(message)),
    "Encrypted Size": fmt.Sprintf("%.2f KB", float64(len(encrypted))/1024),
  }
  ui.PrintDataDetails(details)

  ui.ShowSuccess("Message hidden successfully in the image")
  ui.PrintKeyBox(hex.EncodeToString(encryptor.GetKey()))

  return nil
}

func handleHideFileCommand(ui *ui.UI) error {
  ui.PrintCommandHeader("HIDE FILE IN IMAGE")

  // Collect input information
  inputPath := ui.PromptInput("Enter input image path (PNG or JPG)")
  if !fileExists(inputPath) {
    return fmt.Errorf("input file does not exist: %s", inputPath)
  }

  outputPath := ui.PromptInput("Enter output image path (will be saved as PNG)")
  if !strings.HasSuffix(strings.ToLower(outputPath), ".png") {
    outputPath += ".png"
  }

  filePath := ui.PromptInput("Enter path to the file you want to hide")
  if !fileExists(filePath) {
    return fmt.Errorf("file does not exist: %s", filePath)
  }

  ui.StartProgress("Checking file compatibility")
  fileHandler := steganography.NewFileHandler()
  supported, ext := fileHandler.IsFileSupported(filePath)
  if !supported {
    ui.ShowWarning(fmt.Sprintf("File type %s is not in the standard supported list, but we'll try anyway", ext))
  }

  ui.UpdateProgress("Reading file data")
  fileData, metadata, err := fileHandler.ReadFileContent(filePath)
  if err != nil {
    ui.StopProgress()
    return fmt.Errorf("failed to read file: %v", err)
  }

  ui.UpdateProgress("Initializing encoder")
  encoder, err := steganography.NewEncoder(inputPath)
  if err != nil {
    ui.StopProgress()
    return fmt.Errorf("failed to initialize encoder: %v", err)
  }

  ui.UpdateProgress("Generating encryption key")
  encryptor, err := crypto.NewEncryptor()
  if err != nil {
    ui.StopProgress()
    return fmt.Errorf("failed to initialize encryption: %v", err)
  }

  ui.UpdateProgress("Encrypting file data")
  encrypted, err := encryptor.Encrypt(fileData)
  if err != nil {
    ui.StopProgress()
    return fmt.Errorf("failed to encrypt file data: %v", err)
  }

  ui.UpdateProgress("Hiding file in image")
  if err := encoder.HideFile(encrypted, metadata); err != nil {
    ui.StopProgress()
    return fmt.Errorf("failed to hide file: %v", err)
  }

  ui.UpdateProgress("Saving output image")
  if err := encoder.SaveOutput(outputPath); err != nil {
    ui.StopProgress()
    return fmt.Errorf("failed to save output image: %v", err)
  }
  ui.StopProgress()

  details := map[string]string{
    "Input Image": inputPath,
    "Output Image": outputPath,
    "File Name": metadata.OriginalName,
    "File Type": metadata.FileExt,
    "File Size": fmt.Sprintf("%.2f KB", float64(metadata.FileSize)/1024),
    "Encrypted Size": fmt.Sprintf("%.2f KB", float64(len(encrypted))/1024),
  }
  ui.PrintDataDetails(details)

  ui.ShowSuccess("File hidden successfully in the image")
  ui.PrintKeyBox(hex.EncodeToString(encryptor.GetKey()))

  return nil
}

func handleExtractCommand(ui *ui.UI) error {
  ui.PrintCommandHeader("EXTRACT HIDDEN CONTENT")

  inputPath := ui.PromptInput("Enter image path")
  if !fileExists(inputPath) {
    return fmt.Errorf("file does not exist: %s", inputPath)
  }

  keyStr := ui.PromptInput("Enter encryption key (hex)")
  keyStr = strings.TrimSpace(keyStr)
  if len(keyStr) != 64 {
    return fmt.Errorf("invalid key length. Expected 64 hexadecimal characters")
  }

  ui.StartProgress("Validating encryption key")
  key, err := hex.DecodeString(keyStr)
  if err != nil {
    ui.StopProgress()
    return fmt.Errorf("invalid encryption key format: must be hexadecimal")
  }

  ui.UpdateProgress("Initializing decoder")
  decoder, err := steganography.NewDecoder(inputPath)
  if err != nil {
    ui.StopProgress()
    return fmt.Errorf("failed to initialize decoder: %v", err)
  }

  ui.UpdateProgress("Extracting hidden content")
  data, isFile, metadata, err := decoder.Extract()
  if err != nil {
    ui.StopProgress()
    if err.Error() == "no steganographic data found" {
      return fmt.Errorf("no hidden content found in this image")
    }
    return fmt.Errorf("failed to extract content: %v", err)
  }

  ui.UpdateProgress("Initializing decryption")
  encryptor, err := crypto.NewEncryptorWithKey(key)
  if err != nil {
    ui.StopProgress()
    return fmt.Errorf("failed to initialize decryption: %v", err)
  }

  ui.UpdateProgress("Decrypting content")
  decrypted, err := encryptor.Decrypt(data)
  if err != nil {
    ui.StopProgress()
    return fmt.Errorf("failed to decrypt content: %v", err)
  }
  ui.StopProgress()

  if isFile && metadata != nil {
    details := map[string]string{
      "Content Type": "File",
      "File Name": metadata.OriginalName,
      "File Type": metadata.FileExt,
      "File Size": fmt.Sprintf("%.2f KB", float64(len(decrypted))/1024),
    }
    ui.PrintDataDetails(details)

    outputPath := ui.PromptInput("Enter path to save the extracted file (or press Enter to use original filename)")
    if outputPath == "" {
      outputPath = metadata.OriginalName
    }

    ui.StartProgress("Saving extracted file")
    fileHandler := steganography.NewFileHandler()
    if err := fileHandler.SaveFileContent(decrypted, metadata, outputPath); err != nil {
      ui.StopProgress()
      return fmt.Errorf("failed to save extracted file: %v", err)
    }
    ui.StopProgress()

    ui.ShowSuccess(fmt.Sprintf("File extracted and saved to: %s", outputPath))
  } else {
    details := map[string]string{
      "Content Type": "Text Message",
      "Length": fmt.Sprintf("%d characters", len(decrypted)),
      "Input Image": inputPath,
    }
    ui.PrintDataDetails(details)

    ui.ShowSuccess("Message extracted successfully!")

    fmt.Println()
    color.New(color.FgHiCyan).Println("  ┌─ Extracted Message ───────────────────────────┐")
    messageLines := splitMessage(string(decrypted), 45)
    for _, line := range messageLines {
      color.New(color.FgHiCyan).Print("  │ ")
      color.New(color.FgHiWhite).Printf("%s", line)

      padding := 45 - len(line)
      if padding > 0 {
        fmt.Print(strings.Repeat(" ", padding))
      }

      color.New(color.FgHiCyan).Println(" │")
    }
    color.New(color.FgHiCyan).Println("  └─────────────────────────────────────────────┘")
    fmt.Println()
  }

  return nil
}

func splitMessage(message string, maxLength int) []string {
  var lines []string

  words := strings.Fields(message)
  if len(words) == 0 {
    return []string{""}
  }

  currentLine := words[0]
  for i := 1; i < len(words); i++ {
    if len(currentLine)+1+len(words[i]) <= maxLength {
      currentLine += " " + words[i]
    } else {
      lines = append(lines, currentLine)
      currentLine = words[i]
    }
  }

  lines = append(lines, currentLine)
  return lines
}

func showInfo(ui *ui.UI) {
  ui.PrintCommandHeader("ABOUT STEG-GO")

  ui.ShowInfo(fmt.Sprintf("Steg-Go - Secret Text Embedded Generates Output v%s", appVersion))
  ui.ShowInfo("Created by pranaykumar2")

  ui.PrintFeatureList("Security Features", []string{
    "AES-256-GCM encryption for all hidden content",
    "Secure key generation using crypto/rand",
    "LSB (Least Significant Bit) steganography algorithm",
    "Encrypted metadata for file operations",
  })

  ui.PrintFeatureList("Supported Formats", []string{
    "Input images: PNG, JPG/JPEG",
    "Output images: PNG (for maximum data integrity)",
    "Embeddable files: PDF, DOC/DOCX, TXT, MP3, WAV, and many more",
  })

  ui.PrintFeatureList("Capabilities", []string{
    "Hide text messages in images",
    "Hide entire files in images (documents, audio, etc.)",
    "Automatic file type detection and handling",
    "Secure encryption of all embedded content",
    "Advanced terminal UI with progress indicators",
  })

  ui.PrintFeatureList("Usage Examples", []string{
    fmt.Sprintf("%s hide     - Hide a text message in an image", os.Args[0]),
    fmt.Sprintf("%s hideFile - Hide a file in an image", os.Args[0]),
    fmt.Sprintf("%s extract  - Extract hidden content from an image", os.Args[0]),
    fmt.Sprintf("%s metadata - Show metadata of an image", os.Args[0]),
  })
}

func fileExists(path string) bool {
  _, err := os.Stat(path)
  return !os.IsNotExist(err)
}

func testWithRealImage(imagePath string) {
  // (for development purposes)
}
