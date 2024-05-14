package token_test

import (
	"crawlquery/node/token"
	"reflect"
	"strings"
	"testing"

	testdataloader "github.com/peteole/testdata-loader"
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

func TestKeywords(t *testing.T) {

	testdata := testdataloader.GetTestFile("testdata/pages/google/search.html")

	if len(testdata) == 0 {
		t.Fatal("Failed to load test data")
	}

	//  [google aboutstore gmail images wrong history deleted choose youre giving feedback delete delete report inappropriate predictions dismiss watch google i/o learn latest innovations news ai updates united kingdom advertisingbusiness search works decade climate action join privacyterms settings search settings advanced search data search search history search send feedback dark theme google apps google account]
	expectedKeywords := []string{"google", "aboutstore", "gmail", "images", "wrong", "history", "deleted", "choose", "youre", "giving", "feedback", "delete", "report", "inappropriate", "predictions", "dismiss", "watch", "google", "i/o", "learn", "latest", "innovations", "news", "ai", "updates", "united", "kingdom", "advertisingbusiness", "search", "works", "decade", "climate", "action", "join", "privacyterms", "settings", "search", "settings", "advanced", "search", "data", "search", "search", "history", "search", "send", "feedback", "dark", "theme", "google", "apps", "google", "account"}

	// Extract keywords from the input data
	keywords := token.Keywords(token.Tokenize(string(testdata)))

	if len(expectedKeywords) != len(keywords) {
		t.Errorf("Keywords() = %v, want %v", keywords, expectedKeywords)
	}

	for _, expectedKeyword := range expectedKeywords {
		var found bool
		for _, keyword := range keywords {
			if keyword == expectedKeyword {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected keyword %s was not found", expectedKeyword)
		}
	}
}
