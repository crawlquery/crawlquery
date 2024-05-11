package service_test

import (
	"crawlquery/node/service"
	"crawlquery/pkg/factory"
	"crawlquery/pkg/index"
	"crawlquery/pkg/repository/index/mem"
	"testing"
)

func TestSearch(t *testing.T) {
	idx := index.NewIndex()
	memRepo := mem.NewMemoryRepository()
	memRepo.Save(idx)

	for _, page := range factory.TenPages() {
		idx.AddPage(page)
	}

	svc := service.NewIndexService(memRepo)

	res, err := svc.Search("homepage")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if len(res) != 1 {
		t.Errorf("Expected 1 result, got %v", len(res))
	}
}

// Output:
