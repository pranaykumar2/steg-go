[![Banner](https://capsule-render.vercel.app/api?type=waving&color=gradient&height=150&section=header&text=Steg-Go%20-%20Image%20Steganography&fontSize=30&animation=fadeIn&fontAlignY=35&desc=Hide%20secrets%20in%20plain%20sight%20with%20military-grade%20encryption!&descAlignY=51&descAlignX=50)](https://github.com/pranaykumar2/steg-go)

<p align="center">
  <a href="https://golang.org"><img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go" alt="Go Version"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square" alt="License"></a>
  <a href="https://github.com/pranaykumar2/steg-go/stargazers"><img src="https://img.shields.io/github/stars/pranaykumar2/steg-go?style=social" alt="Stars"></a>
</p>

---

## 🚀 **What is Steg-Go?**
🔒 **Steg-Go** is a powerful **command-line tool** that hides **encrypted messages inside images** using **Least Significant Bit (LSB) Steganography** and **AES-256 encryption**. Protect your sensitive data by embedding it imperceptibly in PNG/JPEG images!

---

## 🌟 **Features**
<p align="center">
  <img src="https://img.shields.io/badge/Steganography-Hide%20messages%20in%20images-blueviolet?style=for-the-badge"/>
  <img src="https://img.shields.io/badge/Encryption-AES--256%20military--grade-red?style=for-the-badge"/>
  <img src="https://img.shields.io/badge/CLI%20Tool-Simple%20Command--line%20Interface-orange?style=for-the-badge"/>
</p>

✅ **Steganography**: Hide messages **imperceptibly** in PNG/JPEG images  
✅ **AES-256 Encryption**: Military-grade encryption for maximum security  
✅ **Intuitive CLI**: User-friendly, fast command-line tool  
✅ **Format Support**: Works with PNG and JPG/JPEG input images  
✅ **Data Integrity**: Ensures hidden data remains intact  

---

## 📦 **Installation**
### **🔧 From Source (Unix Based Systems)**
```bash
# Clone the repository
git clone https://github.com/pranaykumar2/steg-go.git
cd steg-go

# Change permissions to executable
chmod +x build.sh

# Build the application
./build.sh

# Run the application
./stego info


# 🌟 Running **Steg-Go** on Windows 🚀  

Make **Steg-Go** work seamlessly on Windows with these simple steps!  

---

## 📌 **1. Install Go on Windows** 🛠️  

1️⃣ Download the Go installer for Windows from 👉 [Go Official Site](https://golang.org/dl/)  
2️⃣ Run the downloaded `.msi` file and follow the installation steps  
3️⃣ Verify your installation by opening **Command Prompt** and running:  

   ```sh
   go version
   ```  
   ✅ If you see a version number, you're good to go!  

---

## 🔗 **2. Clone the Repository** 🌍  

Open **Git Bash** or **Command Prompt** and run:  

```sh
git clone https://github.com/pranaykumar2/steg-go.git
cd steg-go
```  

---

## ⚙️ **3. Build the Application** 🔨  

Since `build.sh` is meant for Unix systems, use these Windows-friendly commands instead:  

```sh
go mod tidy
go build -o stego.exe ./cmd/stego
```  

---

## 🚀 **4. Run the Application** 🎯  

Once built, you can run **Steg-Go** using:  

```sh
.\stego.exe info
```  

### ➡️ Other Available Commands  

| Command | Description |
|---------|------------|
| `.\stego.exe hide` | Hide data inside an image |
| `.\stego.exe extract` | Extract hidden data from an image |

---

## 📜 **Alternative: Windows Batch File**   

To simplify the process, create a batch file (`build.bat`) with the following content:  

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

Run the batch file by executing:  

```sh
build.bat
```  

---

## 🛠 **Troubleshooting** 🧐  

💡 **Facing issues? Try these fixes:**  

- 🔹 **Permission Denied?** → Run **Command Prompt** as Administrator  
- 📂 **Image Not Found?** → Ensure the image file is in the correct directory or use the full path  
- 🎨 **Weird Terminal UI?** → Try using **Windows Terminal** (supports ANSI color codes better)  

🚀 **Enjoy using Steg-Go on Windows!** 🖼️🔐  

---

## ▶️ **Run It on Replit**
Want to try **Steg-Go** instantly without installation? Click below to launch on **Replit**.

[![Run on Replit](https://replit.com/badge/github/pranaykumar2/steg-go)](https://replit.com/github/pranaykumar2/steg-go)

---

## 🚀 **Usage**
### **🔹 Hide a Secret Message**
```bash
./stego hide
```
📌 **Example Session:**
```yaml
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

---

### **🔹 Extract a Hidden Message**
```bash
./stego extract
```
📌 **Example Session:**
```yaml
➜ Enter image path: 
sample-hidden.png

➜ Enter encryption key (hex): 
5e365d1e972297e6f6b028a6720385a1ccf126463a111537687aa1713024c4c6

ℹ Extracting message...
✓ Message extracted successfully!

Extracted message: This is a top secret message!
```

---

## 🛠️ **How It Works**
### **🔹 Least Significant Bit (LSB) Steganography**
🖼️ **Steg-Go** modifies the **least significant bit** of each pixel's **RGB values** to hide messages:

```mermaid
flowchart LR
    subgraph "Pixel Value Modification"
        P["RGB Pixel<br/><b>Original:</b><br/>R: 100 = 01100100<br/>G: 150 = 10010110<br/>B: 200 = 11001000"] 
        
        M["<b>Message Bits:</b><br/>1, 1, 1"]
        
        N["<b>New Values:</b><br/>R: 101 = 01100101<br/>G: 151 = 10010111<br/>B: 201 = 11001001"]
    end
    
    P -->|"Replace LSB"| M
    M -->|"Embed"| N
```

### **🔹 AES-256 Encryption**
Before embedding, messages are **encrypted** using **AES-256**:

1. **Generate a random 32-byte key**
2. **Encrypt the message using AES-GCM**
3. **Embed both the encrypted data & nonce in the image**
4. **Same key is required for decryption**

---

## 🔒 **Security Considerations**
<p align="center">
  <img src="https://img.shields.io/badge/Invisible%20Data-Hidden%20in%20LSB-purple?style=for-the-badge"/>
  <img src="https://img.shields.io/badge/AES--256%20Encryption-Ultra%20Secure-red?style=for-the-badge"/>
  <img src="https://img.shields.io/badge/Output%20Format-PNG%20(Safe%20from%20Lossy%20Compression)-blue?style=for-the-badge"/>
</p>

✔ **Invisible Embedding**: Data is hidden at the **bit level**, undetectable to the **human eye**  
✔ **AES-256 Security**: Even if extracted, the message remains encrypted  
✔ **Lossless Image Format**: Output images are **PNG**, preventing **compression loss**  

---

## 🧪 **Technical Details**
| Component | Description |
|-----------|-------------|
| `cmd/stego` | Main application entry point |
| `internal/steganography` | Handles hiding and extracting data |
| `internal/crypto` | Manages encryption and decryption |
| `pkg/imageprocessing` | Processes and manipulates image data |
| `internal/ui` | Command-line user interface |

---

## 📸 **Before and After Comparison**
| Original Image | Image with Hidden Message |
|---------------|---------------------------|
| ![Original](sample-image.jpg) | ![Hidden](sample-hidden-image.png) |

🖼️ *Notice how the images appear identical? That's the power of LSB steganography!*

---

## 📝 **License**
This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

## 👥 **Contributing**
Contributions are **welcome**! Follow these steps to contribute:

1. **Fork** the repository  
2. **Create a feature branch** (`git checkout -b feature/amazing-feature`)  
3. **Commit your changes** (`git commit -m 'Add some amazing feature'`)  
4. **Push to the branch** (`git push origin feature/amazing-feature`)  
5. **Open a Pull Request**  

---

### 🎉 **Created with ❤️ by [pranaykumar2](https://github.com/pranaykumar2)**
