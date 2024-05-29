package service_test

import (
	"crawlquery/api/domain"
	linkRepo "crawlquery/api/link/repository/mem"
	linkService "crawlquery/api/link/service"

	crawlJobRepo "crawlquery/api/crawl/job/repository/mem"
	crawlLogRepo "crawlquery/api/crawl/log/repository/mem"
	crawlService "crawlquery/api/crawl/service"

	shardRepo "crawlquery/api/shard/repository/mem"
	shardService "crawlquery/api/shard/service"

	pageRepo "crawlquery/api/page/repository/mem"
	pageService "crawlquery/api/page/service"

	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Run("creates a link with page", func(t *testing.T) {
		// Arrange
		linkRepo := linkRepo.NewRepository()
		crawlLogRepo := crawlLogRepo.NewRepository()

		crawlService := crawlService.NewService(
			crawlService.WithLogger(testutil.NewTestLogger()),
			crawlService.WithCrawlJobRepo(crawlJobRepo.NewRepository()),
			crawlService.WithCrawlLogRepo(crawlLogRepo),
		)

		shardRepo := shardRepo.NewRepository()
		shardService := shardService.NewService(
			shardService.WithLogger(testutil.NewTestLogger()),
			shardService.WithRepo(shardRepo),
		)
		shardRepo.Create(&domain.Shard{
			ID: 0,
		})

		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(
			pageService.WithLogger(testutil.NewTestLogger()),
			pageService.WithPageRepo(pageRepo),
			pageService.WithCrawlService(crawlService),
			pageService.WithShardService(shardService),
		)

		linkService := linkService.NewService(
			linkService.WithPageService(pageService),
			linkService.WithLinkRepo(linkRepo),
			linkService.WithLogger(testutil.NewTestLogger()),
		)

		src := util.PageID("https://cancreatealink.com")
		dst := domain.URL("https://cancreatealink.com/about")

		// Act
		link, err := linkService.Create(src, dst)

		// Assert
		if err != nil {
			t.Errorf("Error adding link: %v", err)
		}

		if link.SrcID != src {
			t.Errorf("Expected srcID %s, got %s", src, link.SrcID)
		}

		if link.DstID != util.PageID(dst) {
			t.Errorf("Expected dstID %s, got %s", util.PageID(dst), link.DstID)
		}

		if link.CreatedAt.IsZero() {
			t.Errorf("Expected CreatedAt to be set")
		}

		repoCheck, _ := linkRepo.GetAllBySrcID(link.SrcID)

		if len(repoCheck) != 1 {
			t.Errorf("Expected 1 link, got %d", len(repoCheck))
		}

		page, err := pageRepo.Get(util.PageID("https://cancreatealink.com/about"))

		if err != nil {
			t.Errorf("Error getting page: %v", err)
		}

		if page.URL != "https://cancreatealink.com/about" {
			t.Errorf("Expected page URL to be https://cancreatealink.com/about, got %s", page.URL)
		}
	})
}
