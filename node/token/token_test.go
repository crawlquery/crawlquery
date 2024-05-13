package token_test

import (
	"crawlquery/node/token"
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
		"example": {0},
		"hello":   {1},
		"world":   {2},
		"this":    {3},
		"is":      {4},
		"a":       {5},
		"simple":  {6},
		"test":    {7},
		"numbers": {8},
		"1234":    {9},
	}

	// Tokenize the input HTML content
	tokensWithPositions := token.Tokenize(htmlContent)

	// Check if the output matches the expected tokens and positions
	for token, positions := range expectedTokens {
		if actualPositions, ok := tokensWithPositions[token]; ok {
			for i, pos := range positions {
				if pos != actualPositions[i] {
					t.Errorf("For token %s, expected position %d, got %d", token, pos, actualPositions[i])
				}
			}
		} else {
			t.Errorf("Expected token %s was not found", token)
		}
	}

	// Check for extra tokens not expected
	for token := range tokensWithPositions {
		if _, ok := expectedTokens[token]; !ok {
			t.Errorf("Unexpected token %s found", token)
		}
	}
}
