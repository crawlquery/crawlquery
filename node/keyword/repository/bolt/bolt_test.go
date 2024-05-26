package bolt_test

import (
	"crawlquery/node/keyword/repository/bolt"
	"os"
	"testing"
)

func TestAddPageKeywords(t *testing.T) {
	r, err := bolt.NewRepository("/tmp/inverted_get_test.db")
	defer os.Remove("/tmp/inverted_get_test.db")

	if err != nil {
		t.Fatalf("error creating repository: %v", err)
	}

	pageID := "pageID"
	keywords := []string{"keyword1", "keyword2"}

	r.AddPageKeywords(pageID, keywords)

	pages, err := r.GetPages("keyword1")

	if err != nil {
		t.Fatalf("error getting pages: %v", err)
	}

	if len(pages) != 1 || pages[0] != pageID {
		t.Errorf("Expected pageID to be in keyword1, got %v", pages)
	}

	pages, err = r.GetPages("keyword2")

	if err != nil {
		t.Fatalf("error getting pages: %v", err)
	}

	if len(pages) != 1 || pages[0] != pageID {
		t.Errorf("Expected pageID to be in keyword2, got %v", pages)
	}
}

func TestRemovePageKeywords(t *testing.T) {
	r, err := bolt.NewRepository("/tmp/inverted_remove_test.db")
	defer os.Remove("/tmp/inverted_remove_test.db")

	if err != nil {
		t.Fatalf("error creating repository: %v", err)
	}

	pageID := "pageID"
	keywords := []string{"keyword1", "keyword2"}

	r.AddPageKeywords(pageID, keywords)

	r.RemovePageKeywords(pageID)

	pages, err := r.GetPages("keyword1")

	if err != nil {
		t.Fatalf("error getting pages: %v", err)
	}

	if len(pages) != 0 {
		t.Errorf("Expected pageID to be removed from keyword1, got %v", pages)
	}

	pages, err = r.GetPages("keyword2")

	if err != nil {
		t.Fatalf("error getting pages: %v", err)
	}

	if len(pages) != 0 {
		t.Errorf("Expected pageID to be removed from keyword2, got %v", pages)
	}
}
