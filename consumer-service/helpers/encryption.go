package helpers

import (
	"consumer-service/logs"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Encrypt encrypts the given data using AES encryption with the provided key.
func Encrypt(data, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		logs.Log.Errorln("Failed to create new cipher:", err)
		return "", err
	}

	plainText := []byte(data)
	// Add padding if necessary
	plainText = PKCS7Padding(plainText, aes.BlockSize)

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize] // Initialization vector (IV)

	// Generate a random IV
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		logs.Log.Errorln("Failed to generate random IV:", err)
		return "", err
	}

	stream := cipher.NewCBCEncrypter(block, iv)
	stream.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	// Encode to base64 for easy transport
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt decrypts the given base64-encoded encrypted data using AES with the provided key.
func Decrypt(encryptedData, key string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		logs.Log.Errorln("Failed to decode base64:", err)
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		logs.Log.Errorln("Failed to create new cipher:", err)
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		logs.Log.Errorln("Ciphertext too short")
		return "", errors.New("ciphertext too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	// Decrypt the ciphertext
	stream := cipher.NewCBCDecrypter(block, iv)
	stream.CryptBlocks(cipherText, cipherText)

	// Remove padding
	plainText, err := PKCS7Unpadding(cipherText, aes.BlockSize)
	if err != nil {
		logs.Log.Errorln("Failed to unpad plaintext:", err)
		return "", err
	}

	return string(plainText), nil
}

// PKCS7Padding applies padding to the plaintext.
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7Unpadding removes padding from the plaintext.
func PKCS7Unpadding(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 || length%blockSize != 0 {
		logs.Log.Errorln("Invalid padding size")
		return nil, errors.New("invalid padding size")
	}

	padding := int(data[length-1])
	if padding > length || padding > blockSize {
		logs.Log.Errorln("Invalid padding")
		return nil, errors.New("invalid padding")
	}

	return data[:length-padding], nil
}
