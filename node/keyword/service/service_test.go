package service_test

import (
	keywordRepo "crawlquery/node/keyword/repository/mem"
	keywordService "crawlquery/node/keyword/service"
	"testing"
)

func TestUpdatePageKeywords(t *testing.T) {
	t.Run("removes old keywords and adds new ones", func(t *testing.T) {
		repo := keywordRepo.NewRepository()

		service := keywordService.NewService(repo)

		pageID := "pageID"
		keywords := [][]string{{"keyword1"}, {"keyword2"}}

		service.UpdatePageKeywords(pageID, keywords)

		pages, _ := repo.GetPages("keyword1")

		if len(pages) != 1 || pages[0] != pageID {
			t.Errorf("Expected pageID to be in keyword1, got %v", pages)
		}

		pages, _ = repo.GetPages("keyword2")

		if len(pages) != 1 || pages[0] != pageID {
			t.Errorf("Expected pageID to be in keyword2, got %v", pages)
		}

		keywords = [][]string{{"keyword3"}, {"keyword4"}}

		service.UpdatePageKeywords(pageID, keywords)

		pages, _ = repo.GetPages("keyword1")

		if len(pages) != 0 {
			t.Errorf("Expected pageID to be removed from keyword1, got %v", pages)
		}

		pages, _ = repo.GetPages("keyword2")

		if len(pages) != 0 {
			t.Errorf("Expected pageID to be removed from keyword2, got %v", pages)
		}

		pages, _ = repo.GetPages("keyword3")

		if len(pages) != 1 || pages[0] != pageID {
			t.Errorf("Expected pageID to be in keyword3, got %v", pages)
		}

		pages, _ = repo.GetPages("keyword4")

		if len(pages) != 1 || pages[0] != pageID {
			t.Errorf("Expected pageID to be in keyword4, got %v", pages)
		}
	})

}

func TestGetPageIDsByKeyword(t *testing.T) {
	t.Run("returns pageIDs for a keyword", func(t *testing.T) {
		repo := keywordRepo.NewRepository()

		service := keywordService.NewService(repo)

		pageID := "pageID"
		keywords := [][]string{{"keyword1"}, {"keyword2"}}

		service.UpdatePageKeywords(pageID, keywords)

		pages, _ := service.GetPageIDsByKeyword("keyword1")

		if len(pages) != 1 || pages[0] != pageID {
			t.Errorf("Expected pageID to be in keyword1, got %v", pages)
		}

		pages, _ = service.GetPageIDsByKeyword("keyword2")

		if len(pages) != 1 || pages[0] != pageID {
			t.Errorf("Expected pageID to be in keyword2, got %v", pages)
		}

		pages, _ = service.GetPageIDsByKeyword("keyword3")

		if len(pages) != 0 {
			t.Errorf("Expected pageID to be removed from keyword3, got %v", pages)
		}

		pages, _ = service.GetPageIDsByKeyword("keyword4")

		if len(pages) != 0 {
			t.Errorf("Expected pageID to be removed from keyword4, got %v", pages)
		}
	})
}
