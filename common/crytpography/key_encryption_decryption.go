package crytpography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
)

// EncryptPrivateKeyHex encrypts a hex string private key using AES-GCM.
// Returns the encrypted data as a Base64-encoded string.
func EncryptPrivateKeyHex(hexPrivateKey string, password []byte) (string, error) {
	// Step 1: Decode the hex string into bytes
	privateKeyBytes, err := hex.DecodeString(hexPrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex string: %v", err)
	}

	// Step 2: Generate a random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %v", err)
	}

	// Step 3: Derive a 32-byte key using Argon2
	key := argon2.Key(password, salt, 3, 32*1024, 4, 32)

	// Step 4: Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher block: %v", err)
	}

	// Step 5: Create a GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM mode: %v", err)
	}

	// Step 6: Generate a random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %v", err)
	}

	// Step 7: Encrypt the private key bytes
	encryptedPrivateKey := gcm.Seal(nil, nonce, privateKeyBytes, nil)

	// Step 8: Combine salt, nonce, and encrypted data
	encryptedData := append(salt, nonce...)
	encryptedData = append(encryptedData, encryptedPrivateKey...)

	// Step 9: Encode the encrypted data as a Base64 string
	encryptedHexPrivateKey := base64.StdEncoding.EncodeToString(encryptedData)
	return encryptedHexPrivateKey, nil
}

// DecryptPrivateKeyHex decrypts a Base64-encoded encrypted private key.
// Returns the original hex string.
func DecryptPrivateKeyHex(encryptedHexPrivateKey string, password []byte) (string, error) {
	// Step 1: Decode the Base64-encoded encrypted data
	encryptedData, err := base64.StdEncoding.DecodeString(encryptedHexPrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode Base64 string: %v", err)
	}

	// Step 2: Extract salt, nonce, and encrypted private key
	if len(encryptedData) < 16+12 { // 16 bytes salt + 12 bytes nonce
		return "", errors.New("invalid encrypted data length")
	}
	salt := encryptedData[:16]
	nonce := encryptedData[16:28]
	encryptedPrivateKey := encryptedData[28:]

	// Step 3: Derive the key using Argon2
	key := argon2.Key(password, salt, 3, 32*1024, 4, 32)

	// Step 4: Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher block: %v", err)
	}

	// Step 5: Create a GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM mode: %v", err)
	}

	// Step 6: Decrypt the private key bytes
	privateKeyBytes, err := gcm.Open(nil, nonce, encryptedPrivateKey, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt private key: %v", err)
	}

	// Step 7: Encode the decrypted bytes as a hex string
	hexPrivateKey := hex.EncodeToString(privateKeyBytes)

	return hexPrivateKey, nil
}
