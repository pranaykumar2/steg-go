package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
)

const (
	// Iterations for PBKDF2. 10000 is a common default.
	// Adjust as needed for security vs. performance trade-off.
	pbkdf2Iterations = 10000
	// Key length for AES-256
	keyLength = 32 // 32 bytes for AES-256
)

// Encrypt derives a key from the password using PBKDF2, then encrypts the data using AES-256 GCM.
// It returns the ciphertext, the salt used for key derivation, and the nonce used for GCM.
func Encrypt(data []byte, password string) (encryptedData []byte, salt []byte, nonce []byte, err error) {
	// 1. Generate a cryptographically secure random salt (16 bytes)
	salt = make([]byte, 16) // Standard salt size
	if _, err = io.ReadFull(rand.Reader, salt); err != nil {
		return nil, nil, nil, err
	}

	// 2. Derive a 32-byte key from the password using PBKDF2
	// SHA3-256 is used as the hash function for PBKDF2 here.
	// You can choose other hash functions like sha256.
	key := pbkdf2.Key([]byte(password), salt, pbkdf2Iterations, keyLength, sha3.New256)

	// 3. Initialize AES-GCM cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, nil, err
	}

	// 4. Generate a cryptographically secure random nonce (GCM standard nonce size is 12 bytes)
	nonce = make([]byte, aesgcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, nil, err
	}

	// 5. Encrypt the data
	// The nonce is prepended to the ciphertext by Seal, or can be handled separately.
	// Here, we return it separately as per the function signature.
	encryptedData = aesgcm.Seal(nil, nonce, data, nil)

	return encryptedData, salt, nonce, nil
}

// Decrypt derives the key from the password and salt using PBKDF2, then decrypts the data using AES-GCM.
func Decrypt(encryptedData []byte, password string, salt []byte, nonce []byte) (data []byte, err error) {
	// 1. Derive the key from the password and the provided salt using PBKDF2
	key := pbkdf2.Key([]byte(password), salt, pbkdf2Iterations, keyLength, sha3.New256)

	// 2. Initialize AES-GCM cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 3. Decrypt the data
	// The nonce must be the same used for encryption.
	data, err = aesgcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, err
	}

	return data, nil
}
