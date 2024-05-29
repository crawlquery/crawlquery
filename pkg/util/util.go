package util

import (
	"crawlquery/api/domain"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"regexp"

	"github.com/google/uuid"
)

func MakeAbsoluteIfRelative(base, link string) (string, error) {
	// Parse the base URL
	baseParsed, err := url.Parse(base)
	if err != nil {
		return "", fmt.Errorf("failed to parse base URL: %w", err)
	}

	// Parse the link URL
	linkParsed, err := url.Parse(link)
	if err != nil {
		return "", fmt.Errorf("failed to parse link: %w", err)
	}

	// If the link is already absolute, check its scheme
	if linkParsed.IsAbs() {
		if linkParsed.Scheme != "http" && linkParsed.Scheme != "https" {
			return "", fmt.Errorf("unsupported URL scheme: %s", linkParsed.Scheme)
		}
		return link, nil
	}

	if len(linkParsed.String()) > 0 && linkParsed.String()[0] == '#' {
		return "", fmt.Errorf("unsupported URL scheme: %s", linkParsed.Scheme)
	}

	// Resolve the relative URL against the base URL
	resolvedURL := baseParsed.ResolveReference(linkParsed)

	if resolvedURL.Scheme != "http" && resolvedURL.Scheme != "https" {
		return "", fmt.Errorf("unsupported URL scheme: %s", resolvedURL.Scheme)
	}
	return resolvedURL.String(), nil
}

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
func PageID(url domain.URL) domain.PageID {
	// Create a new SHA-256 hash.
	hash := sha256.New()
	// Write the URL to the hash.
	hash.Write([]byte(url))
	// Get the resulting hash as a byte slice.
	hashBytes := hash.Sum(nil)
	// Encode the hash bytes as a Base64 string.
	hashString := hex.EncodeToString(hashBytes)
	// Return the first 32 characters of the hexadecimal string.
	return domain.PageID(hashString[:32])
}
