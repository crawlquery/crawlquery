package parse_test

import (
	"crawlquery/node/parse"
	"slices"
	"testing"

	testdataloader "github.com/peteole/testdata-loader"
)

func TestParseGoogle(t *testing.T) {
	testdata := testdataloader.GetTestFile("testdata/pages/google/search.html")

	if len(testdata) == 0 {
		t.Fatal("Failed to load test data")
	}

	page, err := parse.Parse(testdata, "http://google.com")

	if err != nil {
		t.Errorf("Error parsing: %v", err)
	}

	if page.URL != "http://google.com" {
		t.Errorf("Expected URL to be http://google.com, got %s", page.URL)
	}

	if page.Title != "Google" {
		t.Errorf("Expected title to be Google, got %s", page.Title)
	}

	if page.MetaDescription != "" {
		t.Errorf("Expected meta description to be empty, got %s", page.MetaDescription)
	}

	expectKeywords := []string{"google", "search", "images", "news", "gmail"}

	for _, kw := range expectKeywords {
		if !slices.Contains(page.Keywords, kw) {
			t.Errorf("Expected content to contain %s", kw)
		}
	}
}

func TestParseHowToMakeBologneseSauce(t *testing.T) {
	testdata := testdataloader.GetTestFile("testdata/pages/recipe/how-to-make-bolognese-sauce.html")

	if len(testdata) == 0 {
		t.Fatal("Failed to load test data")
	}

	page, err := parse.Parse(testdata, "http://example.com/recipe")

	if err != nil {
		t.Errorf("Error parsing: %v", err)
	}

	if page.URL != "http://example.com/recipe" {
		t.Errorf("Expected URL to be http://example.com/recipe, got %s", page.URL)
	}

	if page.Title != "The best spaghetti bolognese recipe | Good Food" {
		t.Errorf("Expected title to be The best spaghetti bolognese recipe | Good Food, got %s", page.Title)
	}

	expectKeywords := []string{"bolognese", "sauce", "recipe", "tomato", "beef", "pasta"}

	for _, kw := range expectKeywords {
		if !slices.Contains(page.Keywords, kw) {
			t.Errorf("Expected content to contain %s", kw)
		}
	}
}
