package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func generateSecret(length int) string {
	secret := make([]byte, length)
	_, err := rand.Read(secret)
	if err != nil {
		panic("Failed to generate secure secret")
	}
	return base64.StdEncoding.EncodeToString(secret)
}

func main() {
	// Generate 32-byte (256-bit) secrets
	accessTokenSecret := generateSecret(32)
	refreshTokenSecret := generateSecret(32)

	fmt.Println("Generated Secrets (save these securely):")
	fmt.Println("=======================================")
	fmt.Printf("ACCESS_TOKEN_SECRET=%s\n", accessTokenSecret)
	fmt.Printf("REFRESH_TOKEN_SECRET=%s\n", refreshTokenSecret)
	fmt.Println("\nAdd these to your environment variables or .env file")
}
