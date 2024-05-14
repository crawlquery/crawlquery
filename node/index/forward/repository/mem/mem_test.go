package mem

import (
	"crawlquery/pkg/domain"
	"testing"
)

func TestSave(t *testing.T) {
	r := NewRepository()
	err := r.Save("page1", &domain.Page{})

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
}
