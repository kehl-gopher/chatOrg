package data

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateSecureAPIKey creates a 32-byte secure API key
func GenerateSecureAPIKey() (string, error) {
	bytes := make([]byte, 40)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// func VerifyAPIKey(apiKey string) bool {
// 	return true
// }
