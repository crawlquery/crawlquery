package token_test

import (
	"bytes"
	"crawlquery/node/token"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	testdataloader "github.com/peteole/testdata-loader"
)

func TestPositions(t *testing.T) {
	// Input text content
	textContent := "hello world this is a simple test numbers 1234 hello this is a"

	// Expected output: map of tokens to positions
	expectedPositions := map[string][]int{
		"hello":   {0, 9},
		"world":   {1},
		"this":    {2, 10},
		"is":      {3, 11},
		"a":       {4, 12},
		"simple":  {5},
		"test":    {6},
		"numbers": {7},
		"1234":    {8},
	}

	// Extract positions from the tokens
	positions := token.Positions(strings.Split(textContent, " "))

	// Check if the output matches the expected positions
	if !reflect.DeepEqual(positions, expectedPositions) {
		t.Errorf("Positions() = %v, want %v", positions, expectedPositions)
	}
}

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

func TestParseGoogle(t *testing.T) {
	testdata := testdataloader.GetTestFile("testdata/pages/google/search.html")

	if len(testdata) == 0 {
		t.Fatal("Failed to load test data")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(testdata))

	keywords := token.Keywords(doc)

	if err != nil {
		t.Errorf("Error parsing: %v", err)
	}

	expectKeywords := []string{"google", "search", "images", "news", "gmail"}

	for _, kw := range expectKeywords {
		if !slices.Contains(keywords, kw) {
			t.Errorf("Expected content to contain %s", kw)
		}
	}
}

func TestParseHowToMakeBologneseSauce(t *testing.T) {
	testdata := testdataloader.GetTestFile("testdata/pages/recipe/how-to-make-bolognese-sauce.html")

	if len(testdata) == 0 {
		t.Fatal("Failed to load test data")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(testdata))

	if err != nil {
		t.Errorf("Error parsing: %v", err)
	}

	keywords := token.Keywords(doc)

	expectKeywords := []string{"bolognese", "sauce", "recipe", "tomato", "beef", "pasta"}

	for _, kw := range expectKeywords {
		if !slices.Contains(keywords, kw) {
			t.Errorf("Expected content to contain %s", kw)
		}
	}
}
