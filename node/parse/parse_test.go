package parse_test

import (
	"bytes"
	"crawlquery/node/parse"
	"slices"
	"testing"

	"github.com/PuerkitoBio/goquery"
	testdataloader "github.com/peteole/testdata-loader"
)

func TestKeywordsGoodle(t *testing.T) {
	testdata := testdataloader.GetTestFile("testdata/pages/google/search.html")

	if len(testdata) == 0 {
		t.Fatal("Failed to load test data")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(testdata))

	if err != nil {
		t.Fatalf("Error loading document: %v", err)
	}

	keywords := parse.Keywords(doc)

	expectKeywords := []string{"google", "search", "images", "news", "gmail"}

	for _, kw := range expectKeywords {
		if !slices.Contains(keywords, kw) {
			t.Errorf("Expected content to contain %s", kw)
		}
	}
}

func TestTitleGoogle(t *testing.T) {
	testdata := testdataloader.GetTestFile("testdata/pages/google/search.html")

	if len(testdata) == 0 {
		t.Fatal("Failed to load test data")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(testdata))

	if err != nil {
		t.Fatalf("Error loading document: %v", err)
	}

	title := parse.Title(doc)

	if title != "Google" {
		t.Errorf("Expected title to be Google, got %s", title)
	}

}

func TestTitleHowToMakeBologneseSauce(t *testing.T) {
	testdata := testdataloader.GetTestFile("testdata/pages/recipe/how-to-make-bolognese-sauce.html")

	if len(testdata) == 0 {
		t.Fatal("Failed to load test data")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(testdata))

	if err != nil {
		t.Fatalf("Error loading document: %v", err)
	}

	title := parse.Title(doc)

	if title != "The best spaghetti bolognese recipe | Good Food" {
		t.Errorf("Expected title to be The best spaghetti bolognese recipe | Good Food, got %s", title)
	}
}

func TestKeywordsHowToMakeBologneseSauce(t *testing.T) {
	testdata := testdataloader.GetTestFile("testdata/pages/recipe/how-to-make-bolognese-sauce.html")

	if len(testdata) == 0 {
		t.Fatal("Failed to load test data")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(testdata))

	if err != nil {
		t.Fatalf("Error loading document: %v", err)
	}

	keywords := parse.Keywords(doc)

	if err != nil {
		t.Errorf("Error parsing: %v", err)
	}

	expectKeywords := []string{"bolognese", "sauce", "recipe", "tomato", "beef", "pasta"}

	for _, kw := range expectKeywords {
		if !slices.Contains(keywords, kw) {
			t.Errorf("Expected content to contain %s", kw)
		}
	}
}
