package ui

import (
  "bufio"
  "os"
  "strings"
  "time"

  "github.com/fatih/color"
)

type UI struct {
  reader   *bufio.Reader
  username string
}

func NewUI(username string) *UI {
  return &UI{
    reader:   bufio.NewReader(os.Stdin),
    username: username,
  }
}

func (u *UI) PrintHeader() {
  color.Cyan("\n╔══════════════════════════════════════════╗")
  color.Cyan("║      Secure Image Steganography Tool     ║")
  color.Cyan("╚══════════════════════════════════════════╝\n")

  color.Blue("Current Time (UTC): %s", time.Now().UTC().Format("2006-01-02 15:04:05"))
  color.Blue("User: %s\n", u.username)
}

func (u *UI) PromptInput(prompt string) string {
  color.Cyan("➜ %s: ", prompt)
  input, _ := u.reader.ReadString('\n')
  return strings.TrimSpace(input)
}

func (u *UI) ShowSuccess(message string) {
  color.Green("✓ %s", message)
}

func (u *UI) ShowError(message string) {
  color.Red("✘ %s", message)
}

func (u *UI) ShowInfo(message string) {
  color.Blue("ℹ %s", message)
}
