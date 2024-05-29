package service_test

import (
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"testing"

	"crawlquery/api/domain"

	shardRepo "crawlquery/api/shard/repository/mem"
	shardService "crawlquery/api/shard/service"

	pageRepo "crawlquery/api/page/repository/mem"
	pageService "crawlquery/api/page/service"
)

func TestGet(t *testing.T) {
	t.Run("returns page", func(t *testing.T) {
		logger := testutil.NewTestLogger()
		pageRepo := pageRepo.NewRepository()
		shardRepo := shardRepo.NewRepository()

		shardService := shardService.NewService(
			shardService.WithRepo(shardRepo),
			shardService.WithLogger(logger),
		)

		pageService := pageService.NewService(
			pageService.WithPageRepo(pageRepo),
			pageService.WithShardService(shardService),
			pageService.WithLogger(logger),
		)

		url := domain.URL("http://example.com")
		pageID := util.PageID(url)

		pageRepo.Create(&domain.Page{
			ID:  pageID,
			URL: url,
		})

		got, err := pageService.Get(pageID)

		if err != nil {
			t.Fatalf("error getting page: %v", err)
		}

		if got.ID != pageID {
			t.Errorf("got page ID %s, want %s", got.ID, pageID)
		}

		if got.URL != url {
			t.Errorf("got page URL %s, want %s", got.URL, url)
		}
	})

	t.Run("returns err if page not found", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(
			pageService.WithPageRepo(pageRepo),
			pageService.WithLogger(testutil.NewTestLogger()),
		)

		_, err := pageService.Get("pageID")

		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("creates page", func(t *testing.T) {

		logger := testutil.NewTestLogger()
		shardRepo := shardRepo.NewRepository()
		shardRepo.Create(&domain.Shard{ID: 0})
		shardService := shardService.NewService(
			shardService.WithRepo(shardRepo),
			shardService.WithLogger(logger),
		)

		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(
			pageService.WithPageRepo(pageRepo),
			pageService.WithShardService(shardService),
			pageService.WithLogger(logger),
		)

		url := domain.URL("http://example.com")
		page, err := pageService.Create(url)

		if err != nil {
			t.Fatalf("error creating page: %v", err)
		}

		if page.ID != util.PageID(url) {
			t.Errorf("got page ID %s, want %s", page.ID, util.PageID(url))
		}

		shardID, err := shardService.GetURLShardID(url)

		if err != nil {
			t.Fatalf("error getting shard ID: %v", err)
		}

		if page.ShardID != shardID {
			t.Errorf("got page ShardID %d, want %d", page.ShardID, shardID)
		}

		if page.ShardID != 0 {
			t.Errorf("got page ShardID %d, want 0", page.ShardID)
		}

		if page.URL != url {
			t.Errorf("got page URL %s, want %s", page.URL, url)
		}
	})

	t.Run("returns error if page already exists", func(t *testing.T) {
		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(
			pageService.WithPageRepo(pageRepo),
			pageService.WithLogger(testutil.NewTestLogger()),
		)

		pageRepo.Create(&domain.Page{
			ID: util.PageID("http://google.com"),
		})

		_, err := pageService.Create("http://google.com")

		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
}
