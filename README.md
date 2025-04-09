<div align="center">

<picture>
  <source media="(prefers-color-scheme: light)" srcset="web/static/img/head.png">
  <img alt="Steg-Go Logo" src="web/static/img/head.png" width="155" height="150">
</picture>

### *Secret Text Embedded Generates Output*

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge)](LICENSE)
[![Stars](https://img.shields.io/github/stars/pranaykumar2/steg-go?style=for-the-badge&logo=github)](https://github.com/pranaykumar2/steg-go/stargazers)
[![Run on Replit](https://img.shields.io/badge/Run%20on-Replit-orange?style=for-the-badge&logo=replit)](https://replit.com/github/pranaykumar2/steg-go)

</div>

---

## 🌟 What is Steg-Go?

**Steg-Go** is a powerful command-line tool that lets you **hide encrypted messages inside ordinary images**. Using advanced **Least Significant Bit (LSB) Steganography** combined with **AES-256 encryption**, it provides a secure way to conceal sensitive information in plain sight.

<div align="center">

### "If you want to keep a secret, you must also hide it from yourself." — George Orwell

</div>

---

## ✨ Key Features

<table align="center">
  <tr>
    <td align="center"><img src="https://raw.githubusercontent.com/PKief/vscode-material-icon-theme/main/icons/lock.svg" width="60"/></td>
    <td align="center"><img src="https://raw.githubusercontent.com/PKief/vscode-material-icon-theme/main/icons/console.svg" width="60"/></td>
    <td align="center"><img src="https://raw.githubusercontent.com/PKief/vscode-material-icon-theme/main/icons/image.svg" width="60"/></td>
  </tr>
  <tr>
    <td align="center"><b>Military-Grade Encryption</b></td>
    <td align="center"><b>Intuitive CLI</b></td>
    <td align="center"><b>Transparent Visual Quality</b></td>
  </tr>
  <tr>
    <td>AES-256 encryption ensures your data remains secure even if steganography is detected</td>
    <td>Simple, guided interface for both hiding and extracting data</td>
    <td>No visible changes to images — your secrets remain truly hidden</td>
  </tr>
</table>

### Security & Stealth

- ✅ **Undetectable to the human eye** - Modifies only the least significant bits
- ✅ **Double-layer protection** - Steganography + encryption
- ✅ **Format preservation** - Maintains image quality
- ✅ **Cross-platform** - Works on Linux, macOS, and Windows

---

## 🚀 Installation

<details>
<summary><b>📦 Unix/Linux/MacOS</b></summary>

```bash
# Clone the repository
git clone https://github.com/pranaykumar2/steg-go.git
cd steg-go

# Build and run
chmod +x build.sh
./build.sh
./stego info
```
</details>

<details>
<summary><b>🪟 Windows</b></summary>

1. **Install Go** from [golang.org/dl](https://golang.org/dl/)

2. **Clone the repository**
   ```cmd
   git clone https://github.com/pranaykumar2/steg-go.git
   cd steg-go
   ```

3. **Build the application**
   ```cmd
   go mod tidy
   go build -o stego.exe ./cmd/stego
   ```

4. **Run Steg-Go**
   ```cmd
   .\stego.exe info
   ```

<details>
<summary>💡 Windows Batch File (Optional)</summary>

Create `build.bat` with the following content:

```batch
@echo off
echo Building Steganography Tool - Initial Setup...

echo Tidying Go modules...
go mod tidy

echo Building application...
go build -v -o stego.exe ./cmd/stego

if %ERRORLEVEL% EQU 0 (
    echo Build successful! Run: .\stego.exe
) else (
    echo Build failed! Check for errors.
    exit /b 1
)
```

Run: `build.bat`
</details>
</details>

<details>
<summary><b>☁️ Try Online</b></summary>
  
No installation required! Try Steg-Go instantly:

[![Run on Replit](https://replit.com/badge/github/pranaykumar2/steg-go)](https://replit.com/github/pranaykumar2/steg-go)
</details>

---

## 🎮 How to Use

<div align="center">
  <img src="terminal-demo.gif" alt="Terminal Demo" width="700">
</div>

### Hide a Secret Message

```bash
./stego hide
```

<details>
<summary>📝 Example Session</summary>

```
╔══════════════════════════════════════════╗
║      Secure Image Steganography Tool     ║
╚══════════════════════════════════════════╝
Current Time (UTC): 2025-03-01 09:41:34
User: runner

➜ Enter input image path (PNG or JPG): 
sample.jpg

➜ Enter output image path (will be saved as PNG): 
sample-hidden.png

➜ Enter the secret message: 
This is a top secret message!

ℹ Processing image...
✓ Message hidden successfully!
ℹ Encryption key (save this!): 5e365d1e972297e6f6b028a6720385a1ccf126463a111537687aa1713024c4c6
```
</details>

### Extract a Hidden Message

```bash
./stego extract
```

<details>
<summary>📝 Example Session</summary>

```
➜ Enter image path: 
sample-hidden.png

➜ Enter encryption key (hex): 
5e365d1e972297e6f6b028a6720385a1ccf126463a111537687aa1713024c4c6

ℹ Extracting message...
✓ Message extracted successfully!

Extracted message: This is a top secret message!
```
</details>

---

## 🔍 The Magic Behind Steg-Go

### Complete Data Flow Process

The diagram below shows how Steg-Go transforms your secret message and embeds it invisibly into an image:
```mermaid
---
config:
  theme: neutral
---
flowchart TD
 subgraph Input["Input Phase"]
        C["Image Processing"]
        A["Original Image"]
        D["Message Processing"]
        B["Secret Message"]
  end
 subgraph Encryption["Encryption Phase"]
        E["Generate AES-256 Key"]
        F["Encrypt Message"]
        G["Encrypted Payload"]
  end
 subgraph Steganography["Steganography Phase"]
        H["Extract Pixel Data"]
        I["Convert to Bit Stream"]
        J["LSB Replacement Algorithm"]
        K["Modified Pixel Data"]
        L["Assemble New Image"]
  end
 subgraph Output["Output Phase"]
        M["Save as PNG"]
        N["Display Encryption Key"]
  end
 subgraph LSB["LSB Modification Detail"]
        LSB2["RGB Components"]
        LSB1["Original Pixel Value"]
        LSB3["LSB Replacement"]
        LSB4["Secret Bit Stream"]
        LSB5["New Pixel Value"]
  end
    A L_A_C_0@-- Load --> C
    B L_B_D_0@-- Prepare --> D
    D L_D_E_0@--> E
    E L_E_F_0@--> F & N
    F L_F_G_0@--> G
    C L_C_H_0@--> H
    G L_G_I_0@--> I
    H --> J
    I L_I_J_0@--> J
    J L_J_K_0@--> K
    K L_K_L_0@--> L
    L L_L_M_0@--> M
    LSB1 L_LSB1_LSB2_0@-- Extract --> LSB2
    LSB2 L_LSB2_LSB3_0@-- Modify Last Bit --> LSB3
    LSB4 L_LSB4_LSB3_0@--> LSB3
    LSB3 L_LSB3_LSB5_0@--> LSB5
    J L_J_LSB_0@-.-> LSB
     C:::process
     A:::data
     D:::process
     B:::data
     E:::process
     F:::process
     G:::data
     H:::process
     I:::process
     J:::process
     K:::data
     L:::process
     M:::process
     N:::data
     LSB2:::detail
     LSB1:::detail
     LSB3:::detail
     LSB4:::detail
     LSB5:::detail
    classDef phase fill:#ffe6cc,stroke:#d79b00,stroke-width:2px,color:#000000
    classDef process fill:#e1f5fe,stroke:#01579b,stroke-width:2px,color:#000000
    classDef data fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#000000
    classDef detail fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#000000
    L_A_C_0@{ animation: fast } 
    L_B_D_0@{ animation: fast } 
    L_D_E_0@{ animation: fast } 
    L_E_F_0@{ animation: fast } 
    L_E_N_0@{ animation: slow } 
    L_F_G_0@{ animation: fast } 
    L_C_H_0@{ animation: fast } 
    L_G_I_0@{ animation: fast } 
    L_I_J_0@{ animation: fast } 
    L_J_K_0@{ animation: fast } 
    L_K_L_0@{ animation: fast } 
    L_L_M_0@{ animation: fast } 
    L_LSB1_LSB2_0@{ animation: fast } 
    L_LSB2_LSB3_0@{ animation: fast } 
    L_LSB4_LSB3_0@{ animation: fast } 
    L_LSB3_LSB5_0@{ animation: fast } 
    L_J_LSB_0@{ animation: fast }
```

### How It Works

1. **Input Phase**: The original image and secret message are loaded and prepared
2. **Encryption Phase**: Your message is secured with AES-256 encryption
3. **Steganography Phase**: The encrypted data is embedded bit by bit into the image
4. **Output Phase**: The modified image is saved, looking identical to the original

The LSB (Least Significant Bit) modification detail shows exactly how each pixel is subtly altered to store your secret data without visible changes.

### LSB Steganography Explained

Steg-Go hides your data by modifying the least significant bit of each color channel in image pixels:

<table align="center">
  <tr>
    <th>Original Pixel</th>
    <th>Secret Bits</th>
    <th>Modified Pixel</th>
  </tr>
  <tr>
    <td>
      R: 100 (01100100)<br>
      G: 150 (10010110)<br>
      B: 200 (11001000)
    </td>
    <td align="center">
      1<br>1<br>1
    </td>
    <td>
      R: 101 (01100101)<br>
      G: 151 (10010111)<br>
      B: 201 (11001001)
    </td>
  </tr>
</table>

### AES-256 Encryption Flow

```mermaid
---
config:
  look: handDrawn
  theme: default
---
graph LR
    A[Original Message] -->|Random Key Generation| B[AES-256 Encryption]
    B --> C[Encrypted Data]
    C -->|Embedding| D[Modified Image]
    E[Original Image] -->|Pixel Modification| D

```

---

## 👀 See the Difference (or Not!)

<div align="center">
  <table>
    <tr>
      <td align="center"><b>Original Image</b></td>
      <td align="center"><b>Image with Secret</b></td>
    </tr>
    <tr>
      <td><img src="sample-image.jpg" width="400"/></td>
      <td><img src="sample-hidden-image.png" width="400"/></td>
    </tr>
  </table>

  <i>Can you spot the difference? Nobody can—that's the point!</i>
</div>

---

## 🧠 Technical Architecture

<div align="center">

```mermaid
classDiagram
    class Main {
        +main()
    }
    class Steganography {
        +HideData()
        +ExtractData()
    }
    class Crypto {
        +Encrypt()
        +Decrypt()
        -GenerateKey()
    }
    class ImageProcessor {
        +LoadImage()
        +SaveImage()
        +ModifyPixels()
    }
    class UI {
        +PrintBanner()
        +GetUserInput()
        +DisplayResult()
    }
    
    Main --> UI
    Main --> Steganography
    Steganography --> Crypto
    Steganography --> ImageProcessor
```

</div>

| Component | Purpose |
|-----------|---------|
| `cmd/stego` | Entry point and command handling |
| `internal/steganography` | Core steganography algorithms |
| `internal/crypto` | Encryption and decryption logic |
| `pkg/imageprocessing` | Image manipulation utilities |
| `internal/ui` | User interface and interaction |

---

## 🛡️ Security Considerations

<div align="center">
  <table>
    <tr>
      <td width="33%" align="center"><b>Visual Security</b></td>
      <td width="33%" align="center"><b>Cryptographic Security</b></td>
      <td width="33%" align="center"><b>Format Security</b></td>
    </tr>
    <tr>
      <td>Changes to the image are imperceptible to human eyes and basic analysis tools</td>
      <td>Even if steganography is detected, the AES-256 encryption makes content unreadable without the key</td>
      <td>Output as PNG preserves all data bits, preventing compression losses that occur with JPEG</td>
    </tr>
  </table>
</div>

---

## 🤝 Contributing

Contributions make the open-source community amazing! Any contributions you make are **greatly appreciated**.

<div align="center">

```mermaid
gitGraph:
    commit id: "Initial"
    branch feature
    checkout feature
    commit id: "Feature"
    commit id: "Tests"
    checkout main
    merge feature
    commit id: "Release"
```

</div>

1. **Fork** the project
2. **Create** your feature branch (`git checkout -b feature/AmazingFeature`)
3. **Commit** your changes (`git commit -m 'Add some AmazingFeature'`)
4. **Push** to the branch (`git push origin feature/AmazingFeature`)
5. **Open** a Pull Request

---

<div align="center">

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

## 🙏 Acknowledgements

* Go Programming Language
* All the amazing contributors
* You, for checking out this project!

<br>

**Created with ❤️ by [pranaykumar2](https://github.com/pranaykumar2), [Vinai2003](https://github.com/vinai2003), [SGN/2100030721](https://github.com/2100030721)**

<br>

[![GitHub followers](https://img.shields.io/github/followers/pranaykumar2?style=social)](https://github.com/pranaykumar2)
[![LinkedIn](https://img.shields.io/badge/LinkedIn-Connect-blue)](https://www.linkedin.com/in/iamypranay/)
[![Go Reference](https://pkg.go.dev/badge/github.com/pranaykumar2/steg-go.svg)](https://pkg.go.dev/github.com/pranaykumar2/steg-go)
</div>
