package crypto

import (
  "crypto/aes"
  "crypto/cipher"
  "crypto/rand"
  "errors"
  "io"
)

const (
  keySize   = 32 
  nonceSize = 12 
)

type Encryptor struct {
  key []byte
}

func NewEncryptor() (*Encryptor, error) {
  key := make([]byte, keySize)
  if _, err := io.ReadFull(rand.Reader, key); err != nil {
    return nil, err
  }
  return &Encryptor{key: key}, nil
}

func NewEncryptorWithKey(key []byte) (*Encryptor, error) {
  if len(key) != keySize {
    return nil, errors.New("invalid key size")
  }
  return &Encryptor{key: key}, nil
}

func (e *Encryptor) Encrypt(data []byte) ([]byte, error) {
  block, err := aes.NewCipher(e.key)
  if err != nil {
    return nil, err
  }

  gcm, err := cipher.NewGCM(block)
  if err != nil {
    return nil, err
  }

  nonce := make([]byte, gcm.NonceSize())
  if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
    return nil, err
  }

  ciphertext := gcm.Seal(nonce, nonce, data, nil)
  return ciphertext, nil
}

func (e *Encryptor) Decrypt(ciphertext []byte) ([]byte, error) {
  block, err := aes.NewCipher(e.key)
  if err != nil {
    return nil, err
  }

  gcm, err := cipher.NewGCM(block)
  if err != nil {
    return nil, err
  }

  if len(ciphertext) < gcm.NonceSize() {
    return nil, errors.New("ciphertext too short")
  }

  nonce := ciphertext[:gcm.NonceSize()]
  ciphertext = ciphertext[gcm.NonceSize():]

  plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
  if err != nil {
    return nil, err
  }

  return plaintext, nil
}

func (e *Encryptor) GetKey() []byte {
  return e.key
}
