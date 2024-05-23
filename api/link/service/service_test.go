package service_test

import (
	"crawlquery/api/domain"
	linkRepo "crawlquery/api/link/repository/mem"
	linkService "crawlquery/api/link/service"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Run("can create a link", func(t *testing.T) {
		// Arrange
		linkRepo := linkRepo.NewRepository()
		linkService := linkService.NewService(linkRepo, testutil.NewTestLogger())

		expected := &domain.Link{
			SrcID: util.UUID(),
			DstID: util.UUID(),
		}

		// Act
		link, err := linkService.Create(expected.SrcID, expected.DstID)

		// Assert
		if err != nil {
			t.Errorf("Error adding link: %v", err)
		}

		if link.SrcID != expected.SrcID {
			t.Errorf("Expected srcID %s, got %s", expected.SrcID, link.SrcID)
		}

		if link.DstID != expected.DstID {
			t.Errorf("Expected dstID %s, got %s", expected.DstID, link.DstID)
		}

		if link.CreatedAt.IsZero() {
			t.Errorf("Expected CreatedAt to be set")
		}

		repoCheck := linkRepo.GetAllBySrcID(link.SrcID)

		if len(repoCheck) != 1 {
			t.Errorf("Expected 1 link, got %d", len(repoCheck))
		}
	})
}
