package service_test

import (
	linkRepo "crawlquery/api/link/repository/mem"
	linkService "crawlquery/api/link/service"

	crawlJobRepo "crawlquery/api/crawl/job/repository/mem"
	crawlJobService "crawlquery/api/crawl/job/service"

	crawlRestrictionRepo "crawlquery/api/crawl/restriction/repository/mem"
	crawlRestrictionService "crawlquery/api/crawl/restriction/service"

	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Run("can create a link", func(t *testing.T) {
		// Arrange
		linkRepo := linkRepo.NewRepository()

		crawlRestrictionRepo := crawlRestrictionRepo.NewRepository()
		crawlRestrictionService := crawlRestrictionService.NewService(crawlRestrictionRepo, testutil.NewTestLogger())

		crawlJobRepo := crawlJobRepo.NewRepository()
		crawlJobService := crawlJobService.NewService(crawlJobRepo, nil, nil, crawlRestrictionService, nil, testutil.NewTestLogger())

		linkService := linkService.NewService(linkRepo, crawlJobService, testutil.NewTestLogger())

		src := "https://cancreatealink.com"
		dst := "https://cancreatealink.com/about"

		// Act
		link, err := linkService.Create(src, dst)

		// Assert
		if err != nil {
			t.Errorf("Error adding link: %v", err)
		}

		if link.SrcID != util.PageID(src) {
			t.Errorf("Expected srcID %s, got %s", util.PageID(src), link.SrcID)
		}

		if link.DstID != util.PageID(dst) {
			t.Errorf("Expected dstID %s, got %s", util.PageID(dst), link.DstID)
		}

		if link.CreatedAt.IsZero() {
			t.Errorf("Expected CreatedAt to be set")
		}

		repoCheck := linkRepo.GetAllBySrcID(link.SrcID)

		if len(repoCheck) != 1 {
			t.Errorf("Expected 1 link, got %d", len(repoCheck))
		}
	})

	t.Run("creates a crawl job if the link is new", func(t *testing.T) {
		// Arrange
		linkRepo := linkRepo.NewRepository()

		crawlRestrictionRepo := crawlRestrictionRepo.NewRepository()
		crawlRestrictionService := crawlRestrictionService.NewService(crawlRestrictionRepo, testutil.NewTestLogger())

		crawlJobRepo := crawlJobRepo.NewRepository()
		crawlJobService := crawlJobService.NewService(crawlJobRepo, nil, nil, crawlRestrictionService, nil, testutil.NewTestLogger())

		linkService := linkService.NewService(linkRepo, crawlJobService, testutil.NewTestLogger())

		src := "https://newcrawljob.com"
		dst := "https://newcrawljob.com/about"

		// Act
		_, err := linkService.Create(src, dst)

		// Assert
		if err != nil {
			t.Errorf("Error adding link: %v", err)
		}

		job, err := crawlJobRepo.First()

		if err != nil {
			t.Errorf("Error getting crawl job: %v", err)
		}

		if job == nil {
			t.Fatalf("Expected a crawl job to be created")
		}

		if job.URL != dst {
			t.Errorf("Expected URL %s, got %s", dst, job.URL)
		}
	})
}
