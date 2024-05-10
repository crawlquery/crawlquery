package token_test

import (
	"crawlquery/pkg/token"
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	// Input HTML content
	htmlContent := `<html><head><title>Example</title></head><body><h1>Hello World!</h1><p>This is a simple test. Numbers: 1234.</p></body></html>`

	// Expected output: slice of tokens
	expectedTokens := []string{"hello", "world", "this", "is", "a", "simple", "test", "numbers", "1234"}

	// Tokenize the input HTML content
	tokens := token.Tokenize(htmlContent)

	// Check if the output matches the expected tokens
	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("Tokenize() = %v, want %v", tokens, expectedTokens)
	}
}
