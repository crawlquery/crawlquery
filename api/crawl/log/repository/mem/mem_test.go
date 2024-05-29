package mem

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/util"
	"testing"
)

func TestSave(t *testing.T) {
	t.Run("can save a crawl log", func(t *testing.T) {
		repo := NewRepository()
		cl := &domain.CrawlLog{
			ID:     domain.CrawlLogID(util.UUIDString()),
			PageID: domain.PageID("http://example.com"),
			Status: domain.CrawlStatusFailed,
		}

		err := repo.Save(cl)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(repo.logs) != 1 {
			t.Errorf("Expected 1 log, got %v", len(repo.logs))
		}

		if repo.logs[cl.ID] != cl {
			t.Errorf("Expected log to be saved, got %v", repo.logs[cl.ID])
		}
	})
}

func TestListByPageID(t *testing.T) {
	t.Run("can list logs by pageID", func(t *testing.T) {
		repo := NewRepository()
		pageID := domain.PageID("http://example.com")
		cl1 := &domain.CrawlLog{
			ID:     domain.CrawlLogID(util.UUIDString()),
			PageID: pageID,
			Status: domain.CrawlStatusFailed,
		}
		cl2 := &domain.CrawlLog{
			ID:     domain.CrawlLogID(util.UUIDString()),
			PageID: domain.PageID("http://example.com/2"),
			Status: domain.CrawlStatusSuccess,
		}

		repo.logs[cl1.ID] = cl1
		repo.logs[cl2.ID] = cl2

		logs, err := repo.ListByPageID(pageID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(logs) != 1 {
			t.Errorf("Expected 1 log, got %v", len(logs))
		}

		if logs[0] != cl1 {
			t.Errorf("Expected log to be returned, got %v", logs[0])
		}
	})
}
