package mysql_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/link/repository/mysql"
	"crawlquery/api/migration"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	t.Run("can create a link", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)
		// Arrange
		repo := mysql.NewRepository(db)
		link := &domain.Link{
			SrcID:     util.PageID("https://cancreatealink.com"),
			DstID:     util.PageID("https://cancreatealink.com/about"),
			CreatedAt: time.Now(),
		}

		defer db.Exec("DELETE FROM links WHERE src_id = ?", link.SrcID)

		// Act
		err := repo.Create(link)

		// Assert
		if err != nil {
			t.Errorf("Error adding link: %v", err)
		}
	})

	t.Run("cannot create duplicate link", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)
		// Arrange
		repo := mysql.NewRepository(db)
		link := &domain.Link{
			SrcID:     util.PageID("https://noduplicates.com"),
			DstID:     util.PageID("https://noduplicates.com/about"),
			CreatedAt: time.Now(),
		}

		defer db.Exec("DELETE FROM links WHERE src_id = ?", link.SrcID)

		// Act
		err := repo.Create(link)
		if err != nil {
			t.Errorf("Error adding link: %v", err)
		}
		err = repo.Create(link)

		// Assert
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestGetAll(t *testing.T) {
	t.Run("can get all links", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)
		// Arrange
		repo := mysql.NewRepository(db)
		link1 := &domain.Link{
			SrcID:     util.PageID("https://getalllinks.com"),
			DstID:     util.PageID("https://getalllinks.com/about"),
			CreatedAt: time.Now(),
		}
		link2 := &domain.Link{
			SrcID:     util.PageID("https://getalllinks.com"),
			DstID:     util.PageID("https://getalllinks.com/contact"),
			CreatedAt: time.Now(),
		}
		link3 := &domain.Link{
			SrcID:     util.PageID("https://getalllinks.com"),
			DstID:     util.PageID("https://getalllinks.com/faq"),
			CreatedAt: time.Now(),
		}

		defer db.Exec("DELETE FROM links WHERE src_id = ?", link1.SrcID)
		defer db.Exec("DELETE FROM links WHERE src_id = ?", link2.SrcID)
		defer db.Exec("DELETE FROM links WHERE src_id = ?", link3.SrcID)

		// Act
		err := repo.Create(link1)
		if err != nil {
			t.Errorf("Error adding link: %v", err)
		}
		err = repo.Create(link2)
		if err != nil {
			t.Errorf("Error adding link: %v", err)
		}
		err = repo.Create(link3)
		if err != nil {
			t.Errorf("Error adding link: %v", err)
		}

		links, err := repo.GetAll()

		// Assert
		if err != nil {
			t.Errorf("Error getting links: %v", err)
		}

		if len(links) != 3 {
			t.Errorf("Expected 3 links, got %d", len(links))
		}
	})
}

func TestGetAllBySrcID(t *testing.T) {
	t.Run("can get all links by srcID", func(t *testing.T) {
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)
		// Arrange
		repo := mysql.NewRepository(db)
		link1 := &domain.Link{
			SrcID:     util.PageID("https://getalllinks.com"),
			DstID:     util.PageID("https://getalllinks.com/about"),
			CreatedAt: time.Now(),
		}
		link2 := &domain.Link{
			SrcID:     util.PageID("https://getalllinks.com"),
			DstID:     util.PageID("https://getalllinks.com/contact"),
			CreatedAt: time.Now(),
		}
		link3 := &domain.Link{
			SrcID:     util.PageID("https://getalllinks.com"),
			DstID:     util.PageID("https://getalllinks.com/faq"),
			CreatedAt: time.Now(),
		}

		defer db.Exec("DELETE FROM links WHERE src_id = ?", link1.SrcID)
		defer db.Exec("DELETE FROM links WHERE src_id = ?", link2.SrcID)
		defer db.Exec("DELETE FROM links WHERE src_id = ?", link3.SrcID)

		// Act
		err := repo.Create(link1)
		if err != nil {
			t.Errorf("Error adding link: %v", err)
		}
		err = repo.Create(link2)
		if err != nil {
			t.Errorf("Error adding link: %v", err)
		}
		err = repo.Create(link3)
		if err != nil {
			t.Errorf("Error adding link: %v", err)
		}

		links, err := repo.GetAllBySrcID(link1.SrcID)

		// Assert
		if err != nil {
			t.Errorf("Error getting links: %v", err)
		}

		if len(links) != 3 {
			t.Errorf("Expected 3 links, got %d", len(links))
		}
	})
}
