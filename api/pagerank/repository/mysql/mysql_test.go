package mysql_test

import (
	"crawlquery/api/domain"
	"crawlquery/api/migration"
	"crawlquery/pkg/testutil"
	"testing"
	"time"

	pageRankRepo "crawlquery/api/pagerank/repository/mysql"
)

func TestGet(t *testing.T) {
	db := testutil.CreateTestMysqlDB()
	defer db.Close()

	migration.Up(db)

	repo := pageRankRepo.NewRepository(db)

	// test cases
	tests := []struct {
		// input
		pageID domain.PageID
		rank   float64
	}{

		{"test1", 0.5},
		{"test2", 0.3},
		{"test3", 0.7},
	}

	// setup
	for _, test := range tests {
		_, err := db.Exec("INSERT INTO page_ranks (page_id, `rank`, `created_at`) VALUES (?, ?, ?)", test.pageID, test.rank, time.Now())
		defer db.Exec("DELETE FROM page_ranks WHERE page_id = ?", test.pageID)
		if err != nil {
			t.Fatal(err)
		}
	}

	// test
	for _, test := range tests {
		rank, err := repo.Get(test.pageID)
		if err != nil {
			t.Fatal(err)
		}
		if rank != test.rank {
			t.Errorf("Expected %f, got %f", test.rank, rank)
		}
	}
}

func TestUpdate(t *testing.T) {
	db := testutil.CreateTestMysqlDB()
	defer db.Close()

	migration.Up(db)

	repo := pageRankRepo.NewRepository(db)

	// test cases
	tests := []struct {
		// input
		pageID domain.PageID
		rank   float64
	}{

		{"test1", 0.5},
		{"test2", 0.3},
		{"test3", 0.7},
	}

	// test
	for _, test := range tests {
		err := repo.Update(test.pageID, test.rank, time.Now())
		defer db.Exec("DELETE FROM page_ranks WHERE page_id = ?", test.pageID)
		if err != nil {
			t.Fatal(err)
		}
	}

	// test
	for _, test := range tests {
		rank, err := repo.Get(test.pageID)
		if err != nil {
			t.Fatal(err)
		}
		if rank != test.rank {
			t.Errorf("Expected %f, got %f", test.rank, rank)
		}
	}
}
