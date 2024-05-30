package service_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/testfactory"

	"crawlquery/pkg/util"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Run("creates a link with page", func(t *testing.T) {
		sf := testfactory.NewServiceFactory(
			testfactory.WithShard(&domain.Shard{ID: 0}),
		)

		linkService := sf.LinkService
		linkRepo := sf.LinkRepo
		eventService := sf.EventService

		var publishedEvent bool
		eventService.Subscribe("link.created", func(e domain.Event) {
			publishedEvent = true

			linkCreated := e.(*domain.LinkCreated)

			if linkCreated.DstURL != "https://cancreatealink.com/about" {
				t.Errorf("Expected DstURL to be 'https://cancreatealink.com/about', got '%s'", linkCreated.DstURL)
			}

			if linkCreated.Link.SrcID != util.PageID("https://cancreatealink.com") {
				t.Errorf("Expected SrcID to be 'https://cancreatealink.com', got '%s'", linkCreated.Link.SrcID)
			}

			if linkCreated.Link.DstID != util.PageID("https://cancreatealink.com/about") {
				t.Errorf("Expected DstID to be 'https://cancreatealink.com/about', got '%s'", linkCreated.Link.DstID)
			}
		})

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

		if !publishedEvent {
			t.Errorf("Expected event to be published")
		}
	})
}
