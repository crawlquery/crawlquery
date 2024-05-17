package mem

import (
	"crawlquery/pkg/domain"
	"testing"
)

func TestService(t *testing.T) {
	r := NewRepository()
	err := r.Save("page1", &domain.Page{
		ID: "page1",
	})

	if err != nil {
		t.Fatalf("error saving page: %v", err)
	}

	page, err := r.Get("page1")
	if err != nil {
		t.Fatalf("error getting page: %v", err)
	}

	if page == nil {
		t.Fatalf("expected page to be found")
	}

	if page.ID != "page1" {
		t.Fatalf("expected page ID to be 'page1', got '%s'", page.ID)
	}

	p2, err := r.Get("page2")

	if err == nil {
		t.Fatalf("expected error getting page2, got nil")
	}

	if p2 != nil {
		t.Fatalf("expected page2 to be nil, got %v", p2)
	}

}
