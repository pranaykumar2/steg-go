package main

import (
  "encoding/hex"
  "fmt"
  "os"
  "os/user"
  "strings"
  _ "image/jpeg"
  _ "image/png"
  "github.com/pranaykumar2/steg-go/internal/crypto"
  "github.com/pranaykumar2/steg-go/internal/steganography"
  "github.com/pranaykumar2/steg-go/internal/ui"
)

const (
  appVersion = "1.0.0"
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
    userInterface.ShowError("No command specified. Use 'hide', 'extract', or 'info'")
    printUsage()
    os.Exit(1)
  }

  switch os.Args[1] {
  case "hide":
    if err := handleHideCommand(userInterface); err != nil {
      userInterface.ShowError(fmt.Sprintf("Error: %v", err))
      os.Exit(1)
    }
  case "extract":
    if err := handleExtractCommand(userInterface); err != nil {
      userInterface.ShowError(fmt.Sprintf("Error: %v", err))
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
    userInterface.ShowError("Invalid command. Use 'hide', 'extract', or 'info'")
    printUsage()
    os.Exit(1)
  }
}

func printUsage() {
  fmt.Println("\nUsage:")
  fmt.Printf("  %s <command> [options]\n", os.Args[0])
  fmt.Println("\nCommands:")
  fmt.Println("  hide     Hide a secret message in an image")
  fmt.Println("  extract  Extract a secret message from an image")
  fmt.Println("  info     Show information about this application")
  fmt.Println("\nExamples:")
  fmt.Printf("  %s hide\n", os.Args[0])
  fmt.Printf("  %s extract\n", os.Args[0])
}

func handleHideCommand(ui *ui.UI) error {
  inputPath := ui.PromptInput("Enter input image path (PNG or JPG)")
  if !fileExists(inputPath) {
    return fmt.Errorf("input file does not exist: %s", inputPath)
  }

  encoder, err := steganography.NewEncoder(inputPath)
  if err != nil {
    return fmt.Errorf("failed to initialize encoder: %v", err)
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

  ui.ShowInfo("Processing image...")

  encryptor, err := crypto.NewEncryptor()
  if err != nil {
    return fmt.Errorf("failed to initialize encryption: %v", err)
  }

  encrypted, err := encryptor.Encrypt([]byte(message))
  if err != nil {
    return fmt.Errorf("failed to encrypt message: %v", err)
  }

  if err := encoder.Hide(encrypted); err != nil {
    return fmt.Errorf("failed to hide message: %v", err)
  }

  if err := encoder.SaveOutput(outputPath); err != nil {
    return fmt.Errorf("failed to save output image: %v", err)
  }

  ui.ShowSuccess("Message hidden successfully!")
  ui.ShowInfo(fmt.Sprintf("Encryption key (save this!): %x", encryptor.GetKey()))

  return nil
}

func handleExtractCommand(ui *ui.UI) error {

  inputPath := ui.PromptInput("Enter image path")
  if !fileExists(inputPath) {
    return fmt.Errorf("file does not exist: %s", inputPath)
  }

  keyStr := ui.PromptInput("Enter encryption key (hex)")
  keyStr = strings.TrimSpace(keyStr)
  if len(keyStr) != 64 {
    return fmt.Errorf("invalid key length. Expected 64 hexadecimal characters")
  }
  key, err := hex.DecodeString(keyStr)
  if err != nil {
    return fmt.Errorf("invalid encryption key format: must be hexadecimal")
  }

  ui.ShowInfo("Extracting message...")

  decoder, err := steganography.NewDecoder(inputPath)
  if err != nil {
    return fmt.Errorf("failed to initialize decoder: %v", err)
  }

  data, err := decoder.Extract()
  if err != nil {
    if err.Error() == "no steganographic data found" {
      return fmt.Errorf("no hidden message found in this image")
    }
    return fmt.Errorf("failed to extract message: %v", err)
  }

  encryptor, err := crypto.NewEncryptorWithKey(key)
  if err != nil {
    return fmt.Errorf("failed to initialize decryption: %v", err)
  }

  decrypted, err := encryptor.Decrypt(data)
  if err != nil {
    return fmt.Errorf("failed to decrypt message: %v", err)
  }

  ui.ShowSuccess("Message extracted successfully!")
  fmt.Printf("\nExtracted message: %s\n", string(decrypted))

  return nil
}

func showInfo(ui *ui.UI) {
  ui.ShowInfo(fmt.Sprintf("Steg-Go - Secure Image Steganography Tool v%s", appVersion))
  ui.ShowInfo("Created by pranaykumar2")
  ui.ShowInfo("This tool allows you to hide secret messages in images")
  ui.ShowInfo("Your messages are protected with AES-256 encryption")

  fmt.Println("\nSupported file formats:")
  fmt.Println("  - Input: PNG, JPG/JPEG")
  fmt.Println("  - Output: PNG (for maximum data integrity)")

  fmt.Println("\nUsage examples:")
  fmt.Printf("  %s hide\n", os.Args[0])
  fmt.Printf("  %s extract\n", os.Args[0])
}

func fileExists(path string) bool {
  _, err := os.Stat(path)
  return !os.IsNotExist(err)
}

func testWithRealImage(imagePath string) {
  // (for development purposes)
}