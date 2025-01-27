package service_test

import (
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"testing"

	"crawlquery/api/domain"
	"crawlquery/api/testfactory"

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

func TestHandleLinkCreatedEvent(t *testing.T) {
	t.Run("creates page", func(t *testing.T) {
		sf := testfactory.NewServiceFactory(
			testfactory.WithShard(&domain.Shard{ID: 0}),
		)

		pageService.NewService(
			pageService.WithPageRepo(sf.PageRepo),
			pageService.WithShardService(sf.ShardService),
			pageService.WithEventService(sf.EventService),
			pageService.WithLogger(testutil.NewTestLogger()),
			pageService.WithEventListeners(),
		)

		linkCreated := &domain.LinkCreated{
			DstURL: "http://example.com",
		}

		sf.EventService.Publish(linkCreated)

		repoCheck, err := sf.PageRepo.Get(util.PageID(linkCreated.DstURL))

		if err != nil {
			t.Fatalf("error getting page: %v", err)
		}

		if repoCheck.URL != linkCreated.DstURL {
			t.Fatalf("expected page to not exist, got %v", repoCheck)
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("creates page", func(t *testing.T) {

		sf := testfactory.NewServiceFactory(
			testfactory.WithShard(&domain.Shard{ID: 0}),
		)
		pageService := sf.PageService
		shardService := sf.ShardService
		eventService := sf.EventService

		var pageEventPublished bool
		eventService.Subscribe("page.created", func(event domain.Event) {
			pageEventPublished = true
		})

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

		if !pageEventPublished {
			t.Error("expected page event to be published")
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
