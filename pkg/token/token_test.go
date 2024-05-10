package token_test

import (
	"crawlquery/pkg/token"
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	// Input HTML content
	htmlContent := `<html><head><title>Example</title></head><body><h1>Hello World!</h1><p>This is a simple test. Numbers: 1234.</p></body></html>`

	// Expected output: map of tokens and their positions
	expectedTokens := map[string][]int{
		"hello":   {0},
		"world":   {1},
		"this":    {2},
		"is":      {3},
		"a":       {4},
		"simple":  {5},
		"test":    {6},
		"numbers": {7},
		"1234":    {8},
	}

	// Tokenize the input HTML content
	tokensWithPositions := token.Tokenize(htmlContent)

	// Check if the output matches the expected tokens and positions
	if !reflect.DeepEqual(tokensWithPositions, expectedTokens) {
		t.Errorf("Tokenize() = %v, want %v", tokensWithPositions, expectedTokens)
	}
}
