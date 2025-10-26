# stegano.cc
Steganography and things  
# 🛡️ Encryption Web Application

**stegano.cc** is a lightweight web application for performing **encryption**, **decryption**, and **steganographic embedding** of secret messages inside images.  
Built in **Go (Golang)** with a clean, modular architecture following the **Boundary–Control–Entity (BCE)** pattern, the project separates presentation, business logic, and core domain responsibilities.

Each core feature — `encrypt`, `decrypt`, `embed` — lives in its own top-level folder, containing:

- **Boundary layer** (`boundary/`): UI templates and HTTP handlers.  
- **Control layer** (`control/`): Application logic orchestrating the domain and boundaries.  
- **Entity layer** (shared): Core encryption, key derivation, and image processing logic.  

This organization makes the codebase easy to maintain, test, and extend.

---

### 🧩 Project Structure Overview

```plaintext
stegano.cc/
├── encrypt/      → Encryption feature (boundary + control)
├── decrypt/      → Decryption feature (boundary + control)
├── embed/        → Steganographic embedding (boundary + control)
├── home/         → Landing page and routing
├── key/          → AES key derivation utilities
├── logger/       → Structured logging
├── middleware/   → HTTP middlewares
├── public/       → Static assets (CSS, JS, favicon)
├── template/     → Base and layout HTML templates
├── config/       → Application configuration
├── build/        → Docker and compose setup
├── tests/        → Automated tests and sample images
└── main.go       → Application entrypoint (Fiber web server)
```

## 🐳 Running Locally

### with Docker Compose

You can easily run **stegano.cc** locally using Docker Compose.  
From the project root, execute:

```bash
docker compose -f ./build/docker-compose.yaml up --build
```
Once the build completes, the application will be available at:

👉 http://localhost:3000

To stop the container:
```bash
docker compose -f ./build/docker-compose.yaml down
```
### 💻 Running Locally (Development Mode)

If you prefer to run the app directly with Go (without Docker):

1. Install Go 1.25+
https://go.dev/dl

2. Clone the repository:
```bash
git clone https://github.com/yourusername/stegano.cc.git
cd stegano.cc
```
3. Run the application:
```bash
go run main.go
```
or 
```bash
make run
```
or with air (https://github.com/air-verse/air)
```bash
air
```
4. Open your browser at:

👉 http://localhost:3000

This mode is ideal for fast iteration and local debugging.

## 🚀 Features

- 🔐 **Encrypt** — Derive a 256-bit AES key from a password and encrypt a message securely.
- 🧩 **Embed** — (coming soon) Hide encrypted data inside image files.
- 🔓 **Decrypt** — Reverse the process to recover the original message.


## 🔐 Encrypt

### 📘 Overview

The **Encrypt** feature takes a user-provided **password** and **message**, then produces a secure, Base64-encoded ciphertext.  
The output can later be embedded in a file or decrypted using the same password.

### ⚙️ How it works

1. A **random salt** (16 bytes) is generated for key derivation.  
2. The password and salt are processed through **PBKDF2** with **SHA-256** to derive a **32-byte AES-256 key**.  
3. The message is encrypted using **AES-GCM**, a secure and authenticated cipher mode.  
4. A random **nonce** is generated for each encryption (ensuring uniqueness).  
5. The final ciphertext = `nonce + ciphertext + salt`.  
6. The entire binary blob is **Base64-encoded** to produce printable output.

## 🧩 Embed

### 📘 Overview

The **Embed** feature allows you to **hide encrypted text inside an image file** using **LSB (Least Significant Bit) steganography**.  
It takes a **message** (usually the encrypted Base64 string) and a **source image file**, then produces a **new PNG** image with the message embedded invisibly in its pixels.

This feature can be used to safely store or share encrypted data inside a regular-looking image.

### ⚙️ How it works

1. The user uploads an image (JPEG or PNG) and provides the message to embed.  
2. The system converts the image to **RGBA** format for pixel-level access.  
3. The message is **prefixed with a 4-byte length header** so the extractor knows how many bytes to recover later.  
4. Each bit of the payload is written into the **least significant bits (LSBs)** of the image’s RGB channels — one bit per channel.  
5. The result is re-encoded as a **PNG** file (lossless, ensuring message integrity).

## 🔓 Decrypt

### 📘 Overview

The **Decrypt** feature reverses both encryption and embedding operations.  
It can accept either:

- An **encrypted message string** (Base64-encoded output from the *Encrypt* feature),  
- Or an **image file** that contains an embedded encrypted message (produced by the *Embed* feature).

Given the correct **password**, it decrypts the hidden data and returns the original plaintext message.

If both a message **and** an image are provided, the function will attempt to decrypt **both** independently using the same password and return two separate plaintext results.

### ⚙️ How it works

1. **If an image file is provided:**
   - The system decodes the image and extracts the hidden payload using **LSB decoding** (reverse of embedding).
   - The first 4 bytes define the payload length.
   - The extracted bytes are Base64-decoded to recover the encrypted ciphertext + salt.
   - The **salt** (last 16 bytes) is used to derive the AES-256 key using PBKDF2.
   - The message is decrypted using **AES-GCM**, recovering the original plaintext.

2. **If a message string is provided:**
   - The Base64-encoded ciphertext is decoded.
   - The salt (last 16 bytes) is extracted and used to derive the same AES-256 key.
   - The ciphertext is decrypted with AES-GCM to retrieve the plaintext.

3. Both decrypted results (message and/or image-embedded message) are returned to the caller.

---

### ⚠️ Notes

- Both the encrypted message and embedded file use the **same password** for key derivation.  
- The decryptor can process **either** input type or both simultaneously.  
- If the wrong password is supplied, AES-GCM authentication will fail, preventing corrupted plaintext output.  
- For reliability, always use **lossless image formats (PNG)** when embedding or extracting hidden data.