package mem

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/util"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Run("can create a link", func(t *testing.T) {
		// Arrange
		repo := NewRepository()
		link := &domain.Link{
			SrcID: util.UUIDString(),
			DstID: util.UUIDString(),
		}

		// Act
		err := repo.Create(link)

		// Assert
		if err != nil {
			t.Errorf("Error adding link: %v", err)
		}
	})

	t.Run("cannot create duplicate link", func(t *testing.T) {
		// Arrange
		repo := NewRepository()
		link := &domain.Link{
			SrcID: util.UUIDString(),
			DstID: util.UUIDString(),
		}

		// Act
		err := repo.Create(link)
		if err != nil {
			t.Errorf("Error adding link: %v", err)
		}
		err = repo.Create(link)

		// Assert
		if err != domain.ErrLinkAlreadyExists {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestGetAll(t *testing.T) {
	t.Run("can get all links", func(t *testing.T) {
		// Arrange
		repo := NewRepository()
		link1 := &domain.Link{
			SrcID: util.UUIDString(),
			DstID: util.UUIDString(),
		}
		link2 := &domain.Link{
			SrcID: util.UUIDString(),
			DstID: util.UUIDString(),
		}
		link3 := &domain.Link{
			SrcID: util.UUIDString(),
			DstID: util.UUIDString(),
		}
		repo.Create(link1)
		repo.Create(link2)
		repo.Create(link3)

		// Act
		links, _ := repo.GetAll()

		// Assert
		if len(links) != 3 {
			t.Errorf("Expected 3 links, got %d", len(links))
		}
	})
}

func TestGetAllBySrcID(t *testing.T) {
	t.Run("can get all links by srcID", func(t *testing.T) {
		// Arrange
		repo := NewRepository()
		link1 := &domain.Link{
			SrcID: "srcID",
			DstID: "dstID1",
		}
		link2 := &domain.Link{
			SrcID: "srcID",
			DstID: "dstID2",
		}
		link3 := &domain.Link{
			SrcID: "srcID",
			DstID: "dstID3",
		}
		link4 := &domain.Link{
			SrcID: "srcID2",
			DstID: "dstID4",
		}
		repo.Create(link1)
		repo.Create(link2)
		repo.Create(link3)
		repo.Create(link4)

		// Act
		links, _ := repo.GetAllBySrcID("srcID")

		// Assert
		if len(links) != 3 {
			t.Errorf("Expected 3 links, got %d", len(links))
		}
	})
}
