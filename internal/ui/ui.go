package ui

import (
  "bufio"
  "fmt"
  "os"
  "strings"
  "time"
  "github.com/pranaykumar2/steg-go/pkg/exiftools"
  "github.com/fatih/color"
  "github.com/briandowns/spinner"
)

type UI struct {
  reader      *bufio.Reader
  username    string
  spinnerInst *spinner.Spinner
}

func NewUI(username string) *UI {
  return &UI{
    reader:   bufio.NewReader(os.Stdin),
    username: username,
  }
}

func (u *UI) PrintMetadata(metadata *exiftools.MetadataInfo) {
  fmt.Println()
  color.New(color.FgHiCyan).Println("  ┌─ Image Metadata ─────────────────────────────────┐")

  color.New(color.FgCyan).Println("  │                                                   │")
  color.New(color.FgCyan).Println("  │  FILE INFORMATION                                 │")
  color.New(color.FgCyan).Printf("  │  • Filename:    %-32s │\n", truncateString(metadata.Filename, 32))
  color.New(color.FgCyan).Printf("  │  • File Type:   %-32s │\n", metadata.FileType)
  color.New(color.FgCyan).Printf("  │  • File Size:   %-32s │\n", formatFileSize(metadata.FileSize))
  color.New(color.FgCyan).Printf("  │  • Modified:    %-32s │\n", metadata.ModTime.Format("2006-01-02 15:04:05"))

  if metadata.MimeType != "" {
    color.New(color.FgCyan).Printf("  │  • MIME Type:   %-32s │\n", metadata.MimeType)
  }

  if metadata.ImageWidth > 0 && metadata.ImageHeight > 0 {
    color.New(color.FgCyan).Println("  │                                                   │")
    color.New(color.FgCyan).Println("  │  IMAGE DETAILS                                    │")
    color.New(color.FgCyan).Printf("  │  • Dimensions:  %-32s │\n",
      fmt.Sprintf("%d × %d pixels", metadata.ImageWidth, metadata.ImageHeight))
  }

  color.New(color.FgCyan).Println("  │                                                   │")
  color.New(color.FgCyan).Println("  │  STEGANOGRAPHY INFORMATION                        │")

  capacityBytes := (metadata.ImageWidth * metadata.ImageHeight * 3) / 8
  capacityText := ""
  if capacityBytes > 1048576 {
    capacityText = fmt.Sprintf("%.2f MB", float64(capacityBytes)/1048576)
  } else if capacityBytes > 1024 {
    capacityText = fmt.Sprintf("%.2f KB", float64(capacityBytes)/1024)
  } else {
    capacityText = fmt.Sprintf("%d bytes", capacityBytes)
  }

  color.New(color.FgCyan).Printf("  │  • Max Capacity: %-32s │\n", capacityText)

  exifStatus := "Not detected"
  if metadata.HasEXIF {
    exifStatus = "Present (may contain personal data)"
  }
  color.New(color.FgCyan).Printf("  │  • Metadata:    %-32s │\n", exifStatus)

  if len(metadata.Properties) > 0 {
    color.New(color.FgCyan).Println("  │                                                   │")
    color.New(color.FgCyan).Println("  │  ADDITIONAL PROPERTIES                            │")

    count := 0
    for k, v := range metadata.Properties {
      if k == "Image Width" || k == "Image Height" ||
         k == "MIME Type" || k == "File Modified" ||
         k == "Steganography Capacity" {
        continue
      }

      if count < 5 { // Limit to 5 additional properties
        color.New(color.FgCyan).Printf("  │  • %-12s %-32s │\n",
          truncateString(k+":", 12), truncateString(v, 32))
        count++
      }
    }
  }

  color.New(color.FgHiCyan).Println("  └───────────────────────────────────────────────┘")

  // Privacy Risks Warning
  if len(metadata.PrivacyRisks) > 0 {
    fmt.Println()
    color.New(color.FgHiYellow).Println("  ┌─ PRIVACY CONSIDERATIONS ──────────────────────┐")
    color.New(color.FgHiYellow).Println("  │                                               │")

    for i, risk := range metadata.PrivacyRisks {
      if i < 3 { // Limit to top 3 risks to avoid too much output
        wrappedText := wrapText(risk, 45)
        for j, line := range wrappedText {
          if j == 0 {
            color.New(color.FgYellow).Printf("  │  • %-44s │\n", line)
          } else {
            color.New(color.FgYellow).Printf("  │    %-44s │\n", line)
          }
        }
      }
    }

    color.New(color.FgHiYellow).Println("  │                                               │")
    color.New(color.FgHiYellow).Println("  └───────────────────────────────────────────────┘")
  }
}

func truncateString(s string, maxLen int) string {
  if len(s) <= maxLen {
    return s
  }
  return s[:maxLen-3] + "..."
}

func formatFileSize(size int64) string {
  const (
    _          = iota
    KB float64 = 1 << (10 * iota)
    MB
    GB
  )

  switch {
  case size >= int64(GB):
    return fmt.Sprintf("%.2f GB", float64(size)/GB)
  case size >= int64(MB):
    return fmt.Sprintf("%.2f MB", float64(size)/MB)
  case size >= int64(KB):
    return fmt.Sprintf("%.2f KB", float64(size)/KB)
  default:
    return fmt.Sprintf("%d bytes", size)
  }
}

func wrapText(text string, maxLen int) []string {
  var result []string
  words := strings.Fields(text)
  if len(words) == 0 {
    return []string{""}
  }

  currentLine := words[0]
  for i := 1; i < len(words); i++ {
    if len(currentLine)+1+len(words[i]) <= maxLen {
      currentLine += " " + words[i]
    } else {
      result = append(result, currentLine)
      currentLine = words[i]
    }
  }

  result = append(result, currentLine)
  return result
}


func (u *UI) PrintHeader() {
  // Clear the screen
  fmt.Print("\033[H\033[2J")

  color.New(color.FgHiCyan).Println("\n╔═══════════════════════════════════════════════════════════════╗")
  color.New(color.FgHiCyan).Println("║                                                               ║")
  color.New(color.FgHiCyan).Println("║                          STEG-GO                              ║")
  color.New(color.FgHiCyan).Println("║             Secret Text Embedded Generates Output             ║")
  color.New(color.FgHiCyan).Println("║                                                               ║")
  color.New(color.FgHiCyan).Println("╚═══════════════════════════════════════════════════════════════╝")

  timeStr := time.Now().UTC().Format("2006-01-02 15:04:05")
  color.New(color.FgBlue).Println("\n  SESSION INFORMATION")
  color.New(color.FgBlue).Println("  ─────────────────────────────────────")
  color.New(color.FgHiBlue).Printf("  • Date & Time (UTC): %s\n", timeStr)
  color.New(color.FgHiBlue).Printf("  • User: %s\n", u.username)
  color.New(color.FgHiBlue).Printf("  • Encryption: AES-256-GCM\n")
  color.New(color.FgHiBlue).Printf("  • Version: 1.1.0 (with File Embedding Support)\n\n")
}

func (u *UI) PrintCommandHeader(title string) {
  fmt.Println()
  headerText := fmt.Sprintf(" %s ", strings.ToUpper(title))
  padding := 60 - len(headerText)
  leftPadding := padding / 2
  rightPadding := padding - leftPadding

  line := strings.Repeat("═", 60)
  color.New(color.FgYellow).Println(line)

  left := strings.Repeat(" ", leftPadding)
  right := strings.Repeat(" ", rightPadding)
  color.New(color.FgYellow).Printf("║%s%s%s║\n", left, headerText, right)

  color.New(color.FgYellow).Println(line)
  fmt.Println()
}

func (u *UI) PromptInput(prompt string) string {
  color.New(color.FgCyan, color.Bold).Printf("➜ %s: ", prompt)
  input, _ := u.reader.ReadString('\n')
  return strings.TrimSpace(input)
}

func (u *UI) PromptConfirmation(prompt string) bool {
  color.New(color.FgMagenta, color.Bold).Printf("❓ %s (y/n): ", prompt)
  input, _ := u.reader.ReadString('\n')
  input = strings.TrimSpace(input)
  return strings.ToLower(input) == "y" || strings.ToLower(input) == "yes"
}

func (u *UI) ShowSuccess(message string) {
  color.New(color.FgGreen, color.Bold).Printf("✓ SUCCESS: %s\n", message)
}

func (u *UI) ShowError(message string) {
  color.New(color.FgRed, color.Bold).Printf("✘ ERROR: %s\n", message)
}

func (u *UI) ShowInfo(message string) {
  color.New(color.FgBlue).Printf("ℹ %s\n", message)
}

func (u *UI) ShowWarning(message string) {
  color.New(color.FgYellow).Printf("⚠ WARNING: %s\n", message)
}

func (u *UI) StartProgress(message string) {
  u.spinnerInst = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
  u.spinnerInst.Prefix = "  "
  u.spinnerInst.Suffix = " " + message
  u.spinnerInst.Color("cyan")
  u.spinnerInst.Start()
}

func (u *UI) StopProgress() {
  if u.spinnerInst != nil {
    u.spinnerInst.Stop()
  }
}

func (u *UI) UpdateProgress(message string) {
  if u.spinnerInst != nil {
    u.spinnerInst.Suffix = " " + message
  }
}

func (u *UI) PrintFeatureList(title string, features []string) {
  fmt.Println()
  color.New(color.FgHiMagenta).Printf("  %s:\n", title)
  for _, feature := range features {
    color.New(color.FgHiWhite).Printf("  • %s\n", feature)
  }
  fmt.Println()
}

func (u *UI) PrintDataDetails(details map[string]string) {
  fmt.Println()
  color.New(color.FgHiCyan).Println("  ┌─ Operation Details ───────────────────────────┐")

  for key, value := range details {
    color.New(color.FgCyan).Printf("  │ %-20s", key+":")
    color.New(color.FgHiWhite).Printf(" %-32s │\n", value)
  }

  color.New(color.FgHiCyan).Println("  └─────────────────────────────────────────────┘")
  fmt.Println()
}

func (u *UI) PrintKeyBox(key string) {
  fmt.Println()
  keyLines := splitStringByLength(key, 48)

  color.New(color.FgHiYellow).Println("  ┌─ ENCRYPTION KEY ─────────────────────────────┐")

  for _, line := range keyLines {
    color.New(color.FgHiYellow).Print("  │ ")
    color.New(color.FgHiWhite, color.BgBlack).Printf(" %s ", line)

    padding := 47 - len(line)
    if padding > 0 {
      fmt.Print(strings.Repeat(" ", padding))
    }

    color.New(color.FgHiYellow).Println(" │")
  }

  color.New(color.FgHiYellow).Println("  └─────────────────────────────────────────────┘")
  color.New(color.FgHiRed).Println("    IMPORTANT: Save this key to extract your data later!")
  fmt.Println()
}

func splitStringByLength(input string, length int) []string {
  var result []string
  for i := 0; i < len(input); i += length {
    end := i + length
    if end > len(input) {
      end = len(input)
    }
    result = append(result, input[i:end])
  }
  return result
}
