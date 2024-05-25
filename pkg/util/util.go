package util

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"

	"github.com/google/uuid"
)

func ValidatePageID(pageID string) bool {
	// regex for only alphanumeric characters
	check := regexp.MustCompile(`^[a-zA-Z0-9]*$`)

	return check.MatchString(pageID)
}

func UUID() string {
	return uuid.New().String()
}

func SHA1(s string) string {
	return uuid.NewSHA1(uuid.New(), []byte(s)).String()
}

func Sha256Hex32(b []byte) string {
	hash := sha256.New()
	hash.Write(b)
	return hex.EncodeToString(hash.Sum(nil))[:32]
}

// PageID generates a unique 32-character hash for a given URL.
func PageID(url string) string {
	// Create a new SHA-256 hash.
	hash := sha256.New()
	// Write the URL to the hash.
	hash.Write([]byte(url))
	// Get the resulting hash as a byte slice.
	hashBytes := hash.Sum(nil)
	// Encode the hash bytes as a Base64 string.
	hashString := hex.EncodeToString(hashBytes)
	// Return the first 32 characters of the hexadecimal string.
	return hashString[:32]
}
