package mem

import "testing"

func TestAddPageKeywords(t *testing.T) {
	repo := NewRepository()
	pageID := "pageID"
	keywords := []string{"keyword1", "keyword2"}

	repo.AddPageKeywords(pageID, keywords)

	pages, _ := repo.GetPages("keyword1")
	if len(pages) != 1 || pages[0] != pageID {
		t.Errorf("Expected pageID to be in keyword1, got %v", pages)
	}

	pages, _ = repo.GetPages("keyword2")
	if len(pages) != 1 || pages[0] != pageID {
		t.Errorf("Expected pageID to be in keyword2, got %v", pages)
	}
}

func TestRemovePageKeywords(t *testing.T) {
	repo := NewRepository()
	pageID := "pageID"
	keywords := []string{"keyword1", "keyword2"}

	repo.AddPageKeywords(pageID, keywords)

	repo.RemovePageKeywords(pageID)

	pages, _ := repo.GetPages("keyword1")
	if len(pages) != 0 {
		t.Errorf("Expected pageID to be removed from keyword1, got %v", pages)
	}

	pages, _ = repo.GetPages("keyword2")
	if len(pages) != 0 {
		t.Errorf("Expected pageID to be removed from keyword2, got %v", pages)
	}
}

func TestRemovePageKeywordsNotInKeywords(t *testing.T) {
	repo := NewRepository()
	repo.keywords = map[string][]string{
		"keyword1": {"page1"},
		"keyword2": {"page2"},
	}

	pageID := "page1"

	repo.RemovePageKeywords(pageID)

	pages, _ := repo.GetPages("keyword1")

	if len(pages) != 0 {
		t.Errorf("Expected pageID to be removed from keyword1, got %v", pages)
	}

	pages, _ = repo.GetPages("keyword2")

	if len(pages) != 1 {
		t.Errorf("Expected pageID to be in keyword2, got %v", pages)
	}
}

func TestGetPages(t *testing.T) {
	repo := NewRepository()
	pageID := "pageID"
	keywords := []string{"keyword1", "keyword2"}

	repo.AddPageKeywords(pageID, keywords)

	pages, _ := repo.GetPages("keyword1")
	if len(pages) != 1 || pages[0] != pageID {
		t.Errorf("Expected pageID to be in keyword1, got %v", pages)
	}

	pages, _ = repo.GetPages("keyword2")
	if len(pages) != 1 || pages[0] != pageID {
		t.Errorf("Expected pageID to be in keyword2, got %v", pages)
	}
}
