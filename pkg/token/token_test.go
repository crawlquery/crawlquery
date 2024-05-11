package token_test

import (
	"crawlquery/pkg/token"
	"reflect"
	"testing"
)

func TestTokenizeTerm(t *testing.T) {
	// Input text content
	textContent := "Hello World! This is a simple test. Numbers: 1234."

	// Expected output: slice of tokens
	expectedTokens := []string{"hello", "world", "this", "is", "a", "simple", "test", "numbers", "1234"}

	// Tokenize the input text content
	tokens := token.TokenizeTerm(textContent)

	// Check if the output matches the expected tokens
	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("TokenizeTerms() = %v, want %v", tokens, expectedTokens)
	}
}

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
