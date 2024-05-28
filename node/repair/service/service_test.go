package service_test

import (
	"crawlquery/node/domain"
	"crawlquery/node/dto"
	repairJobRepo "crawlquery/node/repair/job/repository/mem"
	repairJobService "crawlquery/node/repair/service"
	"reflect"

	peerService "crawlquery/node/peer/service"
	"crawlquery/pkg/testutil"
	"testing"
	"time"

	pageRepo "crawlquery/node/page/repository/mem"
	pageService "crawlquery/node/page/service"

	keywordOccurrenceRepo "crawlquery/node/keyword/occurrence/repository/mem"
	keywordService "crawlquery/node/keyword/service"

	"github.com/h2non/gock"
)

func TestGetAllIndexMetas(t *testing.T) {
	t.Run("can get all index metas", func(t *testing.T) {
		defer gock.Off()

		now := time.Now().Round(time.Second)
		oneHourAgo := now.Add(-time.Hour).Round(time.Second)

		expectedMetas := []domain.IndexMeta{
			{
				PageID:        "1",
				PeerID:        "peer1",
				LastIndexedAt: now,
			},
			{
				PageID:        "2",
				PeerID:        "peer1",
				LastIndexedAt: oneHourAgo,
			},
		}

		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		pageRepo.Save("1", &domain.Page{
			ID:            "1",
			LastIndexedAt: &now,
		})

		pageRepo.Save("2", &domain.Page{
			ID:            "2",
			LastIndexedAt: &oneHourAgo,
		})

		pageRepo.Save("3", &domain.Page{
			ID:            "3",
			LastIndexedAt: nil,
		})

		peerService := peerService.NewService(nil, &domain.Peer{
			ID:       "peer1",
			Hostname: "node1",
			Port:     8080,
		}, testutil.NewTestLogger())

		repairJobService := repairJobService.NewService(nil, pageService, nil, peerService, testutil.NewTestLogger())

		metas, err := repairJobService.GetAllIndexMetas()

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(metas) != 2 {
			t.Fatalf("Expected 2 metas, got %v", len(metas))
		}

		for i, meta := range expectedMetas {
			if !reflect.DeepEqual(meta, metas[i]) {
				t.Errorf("Expected meta %v, got %v", expectedMetas[i], meta)
			}
		}
	})
}

func TestGetIndexMetas(t *testing.T) {
	t.Run("can get index metas", func(t *testing.T) {
		defer gock.Off()

		now := time.Now().Round(time.Second)
		oneHourAgo := now.Add(-time.Hour).Round(time.Second)

		expectedMetas := []domain.IndexMeta{
			{
				PageID:        "1",
				PeerID:        "peer1",
				LastIndexedAt: now,
			},
			{
				PageID:        "2",
				PeerID:        "peer1",
				LastIndexedAt: oneHourAgo,
			},
		}

		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		pageRepo.Save("1", &domain.Page{
			ID:            "1",
			LastIndexedAt: &now,
		})

		pageRepo.Save("2", &domain.Page{
			ID:            "2",
			LastIndexedAt: &oneHourAgo,
		})

		pageRepo.Save("3", &domain.Page{
			ID:            "3",
			LastIndexedAt: nil,
		})

		peerService := peerService.NewService(nil, &domain.Peer{
			ID:       "peer1",
			Hostname: "node1",
			Port:     8080,
		}, testutil.NewTestLogger())

		repairJobService := repairJobService.NewService(nil, pageService, nil, peerService, testutil.NewTestLogger())

		metas, err := repairJobService.GetIndexMetas([]string{"1", "2", "3"})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(metas) != 2 {
			t.Fatalf("Expected 2 metas, got %v", len(metas))
		}

		for i, meta := range expectedMetas {
			if !reflect.DeepEqual(meta, metas[i]) {
				t.Errorf("Expected meta %v, got %v", expectedMetas[i], meta)
			}
		}
	})
}

func TestGetPageDumps(t *testing.T) {
	t.Run("can get page dumps", func(t *testing.T) {
		defer gock.Off()

		now := time.Now().Round(time.Second)

		expectedPageDumps := []domain.PageDump{
			{
				PeerID: "peer1",
				PageID: "1",
				Page: domain.Page{
					ID:            "1",
					URL:           "http://google.com",
					Title:         "Google",
					Description:   "Search engine",
					Language:      "English",
					LastIndexedAt: &now,
				},
				KeywordOccurrences: map[domain.Keyword]domain.KeywordOccurrence{
					"keyword1": {
						PageID:    "1",
						Frequency: 1,
						Positions: []int{1, 2, 3},
					},
				},
			},
			{
				PeerID: "peer1",
				PageID: "2",
				Page: domain.Page{
					ID:            "2",
					URL:           "http://example.com",
					Title:         "Example",
					Description:   "Description",
					Language:      "English",
					LastIndexedAt: &now,
				},
				KeywordOccurrences: map[domain.Keyword]domain.KeywordOccurrence{
					"keyword1": {
						PageID:    "2",
						Frequency: 1,
						Positions: []int{1},
					},
				},
			},
		}

		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		pageRepo.Save("1", &domain.Page{
			ID:            "1",
			URL:           "http://google.com",
			Title:         "Google",
			Description:   "Search engine",
			Language:      "English",
			LastIndexedAt: &now,
		})

		pageRepo.Save("2", &domain.Page{
			ID:            "2",
			URL:           "http://example.com",
			Title:         "Example",
			Description:   "Description",
			Language:      "English",
			LastIndexedAt: &now,
		})

		keywordOccurrenceRepo := keywordOccurrenceRepo.NewRepository()
		keywordService := keywordService.NewService(keywordOccurrenceRepo)

		err := keywordOccurrenceRepo.Add("keyword1", domain.KeywordOccurrence{
			PageID:    "1",
			Frequency: 1,
			Positions: []int{1, 2, 3},
		})

		if err != nil {
			t.Fatalf("Error adding keyword occurrence: %v", err)
		}

		err = keywordOccurrenceRepo.Add("keyword1", domain.KeywordOccurrence{
			PageID:    "2",
			Frequency: 1,
			Positions: []int{1},
		})

		if err != nil {
			t.Fatalf("Error adding keyword occurrence: %v", err)
		}

		peerService := peerService.NewService(nil, &domain.Peer{
			ID:       "peer1",
			Hostname: "node1",
			Port:     8080,
		}, testutil.NewTestLogger())

		repairJobService := repairJobService.NewService(nil, pageService, keywordService, peerService, testutil.NewTestLogger())

		pageDumps, err := repairJobService.GetPageDumps([]string{"1", "2"})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(pageDumps) != 2 {
			t.Fatalf("Expected 2 page dumps, got %v", len(pageDumps))
		}

		for i, dump := range expectedPageDumps {
			if !reflect.DeepEqual(dump, pageDumps[i]) {
				t.Errorf("Expected page dump %v, got %v", expectedPageDumps[i], pageDumps[i])
			}
		}
	})
}

func TestCreateRepairJobs(t *testing.T) {
	t.Run("can create repair jobs", func(t *testing.T) {
		repairJobRepo := repairJobRepo.NewRepository()
		repairJobService := repairJobService.NewService(repairJobRepo, nil, nil, nil, testutil.NewTestLogger())

		pages := []string{"1", "2", "3"}
		repairJobService.CreateRepairJobs(pages)

		for _, page := range pages {
			_, err := repairJobRepo.Get(page)
			if err != nil {
				t.Fatalf("Error getting repair job: %v", err)
			}
		}
	})
}

func TestMapLatestPages(t *testing.T) {
	t.Run("can map latest pages with no current pages to refernece", func(t *testing.T) {
		repairJobRepo := repairJobRepo.NewRepository()
		repairJobService := repairJobService.NewService(repairJobRepo, nil, nil, nil, testutil.NewTestLogger())

		threeHoursAgo := time.Now().Add(-time.Hour * 3)
		currentPages := map[string]*domain.Page{
			"1": {
				ID:            "1",
				URL:           "http://google.com",
				Title:         "Google",
				Description:   "Search engine",
				Language:      "English",
				LastIndexedAt: &threeHoursAgo,
			},
			"3": {
				ID:            "3",
				URL:           "http://example.com",
				Title:         "Example",
				Description:   "Description",
				Language:      "English",
				LastIndexedAt: nil,
			},
		}
		metas := []domain.IndexMeta{
			{
				PageID:        "1",
				PeerID:        "peer1",
				LastIndexedAt: time.Now().Add(-time.Hour),
			},
			{
				PageID:        "1",
				PeerID:        "peer2",
				LastIndexedAt: time.Now(),
			},
			{
				PageID:        "2",
				PeerID:        "peer2",
				LastIndexedAt: time.Now().Add(-time.Hour * 2),
			},
			{
				PageID:        "2",
				PeerID:        "peer3",
				LastIndexedAt: time.Now().Add(-time.Hour),
			},
			{
				PageID:        "3",
				PeerID:        "peer3",
				LastIndexedAt: time.Time{},
			},
		}

		latestIndexedAtPeers := repairJobService.MapLatestPages(metas, currentPages)

		if len(latestIndexedAtPeers) != 2 {
			t.Fatalf("Expected 2 latest indexed pages, got %v", len(latestIndexedAtPeers))
		}

		if latestIndexedAtPeers["1"].PeerID != "peer2" {
			t.Fatalf("Expected peer ID to be peer2, got %s", latestIndexedAtPeers["1"].PeerID)
		}

		if latestIndexedAtPeers["2"].PeerID != "peer3" {
			t.Fatalf("Expected peer ID to be peer3, got %s", latestIndexedAtPeers["2"].PeerID)
		}
	})

	t.Run("only maps index meta when current page LastIndexedAt is older", func(t *testing.T) {
		repairJobRepo := repairJobRepo.NewRepository()
		repairJobService := repairJobService.NewService(repairJobRepo, nil, nil, nil, testutil.NewTestLogger())
		date := time.Now()
		currentPages := map[string]*domain.Page{
			"1": {
				ID:            "1",
				URL:           "http://google.com",
				Title:         "Google",
				Description:   "Search engine",
				Language:      "English",
				LastIndexedAt: &date,
			},
		}

		metas := []domain.IndexMeta{
			{
				PageID:        "1",
				PeerID:        "peer1",
				LastIndexedAt: time.Now().Add(-time.Hour),
			},
			{
				PageID:        "1",
				PeerID:        "peer2",
				LastIndexedAt: time.Now().Add(-time.Minute * 2),
			},
			{
				PageID:        "2",
				PeerID:        "peer2",
				LastIndexedAt: time.Now().Add(-time.Hour * 2),
			},
			{
				PageID:        "2",
				PeerID:        "peer3",
				LastIndexedAt: time.Now().Add(-time.Hour),
			},
		}

		latestIndexedAtPeers := repairJobService.MapLatestPages(metas, currentPages)

		if len(latestIndexedAtPeers) != 1 {
			t.Fatalf("Expected 2 latest indexed pages, got %v", len(latestIndexedAtPeers))
		}

		if latestIndexedAtPeers["2"].PeerID != "peer3" {
			t.Fatalf("Expected peer ID to be peer3, got %s", latestIndexedAtPeers["2"].PeerID)
		}
	})
}

func TestGroupPageIDsByThePeerID(t *testing.T) {
	t.Run("can group page IDs by the peer ID", func(t *testing.T) {
		repairJobRepo := repairJobRepo.NewRepository()
		repairJobService := repairJobService.NewService(repairJobRepo, nil, nil, nil, testutil.NewTestLogger())

		latestIndexedAtPeers := domain.LatestIndexedPages{
			"1": domain.PeerWithLatestIndexedAt{
				PeerID:          "peer1",
				LatestIndexedAt: time.Now(),
			},
			"2": domain.PeerWithLatestIndexedAt{
				PeerID:          "peer2",
				LatestIndexedAt: time.Now(),
			},
			"3": domain.PeerWithLatestIndexedAt{
				PeerID:          "peer1",
				LatestIndexedAt: time.Now(),
			},
		}

		peerPages := repairJobService.GroupPageIDsByThePeerID(latestIndexedAtPeers)

		check := peerPages["peer1"]

		if len(check) != 2 {
			t.Fatalf("Expected 2 page IDs, got %v", check)
		}

		if check[0] != "1" && check[0] != "3" {
			t.Fatalf("Expected page IDs to be 1 and 3, got %v", check)
		}

	})
}

func TestProcessRepairJobs(t *testing.T) {
	t.Run("can process repair jobs", func(t *testing.T) {

		defer gock.Off()

		expectedNode1MetasResponse := &dto.GetIndexMetasResponse{
			IndexMetas: []dto.IndexMeta{
				{
					PageID:        "1",
					PeerID:        "peer1",
					LastIndexedAt: time.Now().Add(-time.Hour),
				},
				{
					PageID:        "2",
					PeerID:        "peer1",
					LastIndexedAt: time.Now().Add(-time.Hour),
				},
				{
					PageID:        "3",
					PeerID:        "peer1",
					LastIndexedAt: time.Now(),
				},
				{
					PageID:        "4",
					PeerID:        "peer1",
					LastIndexedAt: time.Now().Add(-time.Minute * 30),
				},
			},
		}

		expectedNode1PageDumpsResponse := &dto.GetPageDumpsResponse{
			PageDumps: []dto.PageDump{
				{
					PeerID: "peer1",
					PageID: "3",
					Page: dto.Page{
						ID:            "3",
						URL:           "http://example.com",
						Title:         "Example",
						Description:   "Description",
						Language:      "English",
						LastIndexedAt: time.Now(),
					},
					KeywordOccurrences: map[string]dto.KeywordOccurrence{
						"keyword1": {
							PageID:    "3",
							Frequency: 1,
							Positions: []int{1},
						},
					},
				},
			},
		}

		expectedNode2MetasResponse := &dto.GetIndexMetasResponse{
			IndexMetas: []dto.IndexMeta{
				{
					PageID:        "1",
					PeerID:        "peer2",
					LastIndexedAt: time.Now(),
				},
				{
					PageID:        "2",
					PeerID:        "peer2",
					LastIndexedAt: time.Now(),
				},
				{
					PageID:        "3",
					PeerID:        "peer2",
					LastIndexedAt: time.Now().Add(-time.Hour),
				},
				{
					PageID:        "4",
					PeerID:        "peer2",
					LastIndexedAt: time.Now().Add(-time.Minute * 45),
				},
			},
		}

		expectedNode2PageDumpsResponse := &dto.GetPageDumpsResponse{
			PageDumps: []dto.PageDump{
				{
					PeerID: "peer1",
					PageID: "1",
					Page: dto.Page{
						ID:            "1",
						URL:           "http://google.com",
						Title:         "Google",
						Description:   "Search engine",
						Language:      "English",
						LastIndexedAt: time.Now(),
					},
					KeywordOccurrences: map[string]dto.KeywordOccurrence{
						"keyword1": {
							PageID:    "1",
							Frequency: 1,
							Positions: []int{1, 2, 3},
						},
					},
				},
				{
					PeerID: "peer1",
					PageID: "2",
					Page: dto.Page{
						ID:            "2",
						URL:           "http://example.com",
						Title:         "Example",
						Description:   "Description",
						Language:      "English",
						LastIndexedAt: time.Now(),
					},
					KeywordOccurrences: map[string]dto.KeywordOccurrence{
						"keyword1": {
							PageID:    "2",
							Frequency: 1,
							Positions: []int{1},
						},
					},
				},
			},
		}

		gock.New("http://node1:8080").
			Post("/repair/get-index-metas").
			JSON(&dto.GetIndexMetasRequest{
				PageIDs: []string{"1", "2", "3", "4"},
			}).
			Reply(200).
			JSON(expectedNode1MetasResponse)

		gock.New("http://node2:8080").
			Post("/repair/get-index-metas").
			JSON(&dto.GetIndexMetasRequest{
				PageIDs: []string{"1", "2", "3", "4"},
			}).
			Reply(200).
			JSON(expectedNode2MetasResponse)

		gock.New("http://node2:8080").
			Post("/repair/get-page-dumps").
			JSON(&dto.GetPageDumpsRequest{
				PageIDs: []string{"1", "2"},
			}).
			Reply(200).
			JSON(expectedNode2PageDumpsResponse)

		gock.New("http://node1:8080").
			Post("/repair/get-page-dumps").
			JSON(&dto.GetPageDumpsRequest{
				PageIDs: []string{"3"},
			}).
			Reply(200).
			JSON(expectedNode1PageDumpsResponse)

		jobs := []string{"1", "2", "3", "4"}

		repairJobRepo := repairJobRepo.NewRepository()

		peerService := peerService.NewService(nil, nil, testutil.NewTestLogger())

		peerService.AddPeer(&domain.Peer{
			ID:       "peer1",
			Hostname: "node1",
			Port:     8080,
		})

		peerService.AddPeer(&domain.Peer{
			ID:       "peer2",
			Hostname: "node2",
			Port:     8080,
		})

		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, peerService)

		keywordOccurrenceRepo := keywordOccurrenceRepo.NewRepository()
		keywordService := keywordService.NewService(keywordOccurrenceRepo)

		repairJobService := repairJobService.NewService(repairJobRepo, pageService, keywordService, peerService, testutil.NewTestLogger())

		now := time.Now()
		// Create page 4 with a last indexed at time of now and make sure it is not updated
		pageRepo.Save("4", &domain.Page{
			ID:            "4",
			URL:           "http://example.com",
			Title:         "Example",
			Description:   "Description",
			Language:      "English",
			LastIndexedAt: &now,
		})

		// add a keyword occurrence for page 3 and make sure it is removed when the page is updated
		err := keywordOccurrenceRepo.Add("check removed", domain.KeywordOccurrence{
			PageID:    "3",
			Frequency: 1,
			Positions: []int{1},
		})

		if err != nil {
			t.Fatalf("Error adding keyword occurrence: %v", err)
		}

		repairJobService.CreateRepairJobs(jobs)

		err = repairJobService.ProcessRepairJobs(jobs)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// check all 3 pages are updated
		for _, job := range jobs {
			page, err := pageRepo.Get(job)
			if err != nil {
				t.Fatalf("Error getting page: %v", err)
			}

			if page == nil {
				t.Fatalf("Expected page to be found")
			}

			if page.ID != job {
				t.Fatalf("Expected page ID to be %s, got %s", job, page.ID)
			}

			if job == "1" {
				if page.URL != "http://google.com" {
					t.Fatalf("Expected page URL to be http://google.com, got %s", page.URL)
				}

				if page.Title != "Google" {
					t.Fatalf("Expected page title to be Google, got %s", page.Title)
				}

				if page.Description != "Search engine" {
					t.Fatalf("Expected page description to be Search engine, got %s", page.Description)
				}

				if page.Language != "English" {
					t.Fatalf("Expected page language to be English, got %s", page.Language)
				}

				if page.LastIndexedAt.Round(time.Second) != expectedNode2PageDumpsResponse.PageDumps[0].Page.LastIndexedAt.Round(time.Second) {
					t.Fatalf("Expected last indexed at to be %v, got %v", expectedNode2PageDumpsResponse.PageDumps[0].Page.LastIndexedAt, page.LastIndexedAt)
				}
			}
		}

		// check all 3 pages have keyword occurrences
		for _, job := range jobs[:3] {
			keywordOccurrences, err := keywordOccurrenceRepo.GetForPageID(job)
			if err != nil {
				t.Fatalf("Error getting keyword occurrences: %v", err)
			}

			if keywordOccurrences == nil {
				t.Fatalf("Expected keyword occurrences to be found")
			}

			if job == "1" {
				if keywordOccurrences["keyword1"].PageID != "1" {
					t.Fatalf("Expected page ID to be 1, got %s", keywordOccurrences["keyword1"].PageID)
				}
				if keywordOccurrences["keyword1"].Frequency != 1 {
					t.Fatalf("Expected frequency to be 1, got %v", keywordOccurrences["keyword1"].Frequency)
				}
				if len(keywordOccurrences["keyword1"].Positions) != 3 {
					t.Fatalf("Expected 3 positions, got %v", len(keywordOccurrences["keyword1"].Positions))
				}
				for i, position := range keywordOccurrences["keyword1"].Positions {
					if position != i+1 {
						t.Fatalf("Expected position to be %v, got %v", i+1, position)
					}
				}
			}

			if job == "2" {
				if keywordOccurrences["keyword1"].PageID != "2" {
					t.Fatalf("Expected page ID to be 2, got %s", keywordOccurrences["keyword1"].PageID)
				}
				if keywordOccurrences["keyword1"].Frequency != 1 {
					t.Fatalf("Expected frequency to be 1, got %v", keywordOccurrences["keyword1"].Frequency)
				}
				if len(keywordOccurrences["keyword1"].Positions) != 1 {
					t.Fatalf("Expected 1 positions, got %v", len(keywordOccurrences["keyword1"].Positions))
				}
				if keywordOccurrences["keyword1"].Positions[0] != 1 {
					t.Fatalf("Expected position to be 1, got %v", keywordOccurrences["keyword1"].Positions[0])
				}
			}

			if job == "3" {
				if keywordOccurrences["keyword1"].PageID != "3" {
					t.Fatalf("Expected page ID to be 3, got %s", keywordOccurrences["keyword1"].PageID)
				}
				if keywordOccurrences["keyword1"].Frequency != 1 {
					t.Fatalf("Expected frequency to be 1, got %v", keywordOccurrences["keyword1"].Frequency)
				}
				if len(keywordOccurrences["keyword1"].Positions) != 1 {
					t.Fatalf("Expected 1 positions, got %v", len(keywordOccurrences["keyword1"].Positions))
				}
				if keywordOccurrences["keyword1"].Positions[0] != 1 {
					t.Fatalf("Expected position to be 1, got %v", keywordOccurrences["keyword1"].Positions[0])
				}
			}
		}

		checkRemoved, err := keywordOccurrenceRepo.GetAll("check removed")

		if err == nil || len(checkRemoved) > 0 {
			t.Fatalf("Expected no keyword occurrences to be found")
		}

		// check page 4 is not updated
		page, err := pageRepo.Get("4")

		if err != nil {
			t.Fatalf("Error getting page: %v", err)
		}

		if page == nil {
			t.Fatalf("Expected page to be found")
		}

		if page.ID != "4" {
			t.Fatalf("Expected page ID to be 4, got %s", page.ID)
		}

		if page.LastIndexedAt.Round(time.Second) != now.Round(time.Second) {
			t.Fatalf("Expected last indexed at to be %v, got %v", now, page.LastIndexedAt)
		}

		keywordOccurrences, err := keywordOccurrenceRepo.GetForPageID("4")

		if err != nil {
			t.Fatalf("Error getting keyword occurrences: %v", err)
		}

		if len(keywordOccurrences) > 0 {
			t.Fatalf("Expected no keyword occurrences to be found")
		}
	})
}

func TestAuditAndRepair(t *testing.T) {
	t.Run("can audit and repair pages (all out of date)", func(t *testing.T) {
		defer gock.Off()

		now := time.Now().Round(time.Second)
		oneHourAgo := now.Add(-time.Hour).Round(time.Second)
		twoHoursAgo := now.Add(-time.Hour * 2).Round(time.Second)
		peer1 := "peer1"
		peer2 := "peer2"

		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		keywordOccurrenceRepo := keywordOccurrenceRepo.NewRepository()
		keywordService := keywordService.NewService(keywordOccurrenceRepo)

		pageRepo.Save("1", &domain.Page{
			ID:            "1",
			URL:           "http://google.com",
			Title:         "Google",
			Description:   "Search engine",
			Language:      "English",
			LastIndexedAt: &twoHoursAgo,
		})

		pageRepo.Save("2", &domain.Page{
			ID:            "2",
			URL:           "http://example.com",
			Title:         "Example",
			Description:   "Description",
			Language:      "English",
			LastIndexedAt: &twoHoursAgo,
		})

		pageRepo.Save("3", &domain.Page{
			ID:            "3",
			URL:           "http://example.org",
			Title:         "Example Org",
			Description:   "Another Description",
			Language:      "English",
			LastIndexedAt: &twoHoursAgo,
		})

		peerService := peerService.NewService(nil, &domain.Peer{
			ID:       "node4",
			Hostname: "node4",
			Port:     8080,
		}, testutil.NewTestLogger())

		peerService.AddPeer(&domain.Peer{
			ID:       peer1,
			Hostname: "peer1",
			Port:     8080,
		})

		peerService.AddPeer(&domain.Peer{
			ID:       "peer2",
			Hostname: "peer2",
			Port:     8080,
		})

		repairJobRepo := repairJobRepo.NewRepository()
		repairJobService := repairJobService.NewService(repairJobRepo, pageService, keywordService, peerService, testutil.NewTestLogger())

		expectedMetas := []domain.IndexMeta{
			{
				PageID:        "1",
				PeerID:        domain.PeerID(peer1),
				LastIndexedAt: now,
			},
			{
				PageID:        "2",
				PeerID:        domain.PeerID(peer1),
				LastIndexedAt: oneHourAgo,
			},
			{
				PageID:        "3",
				PeerID:        domain.PeerID(peer2),
				LastIndexedAt: now,
			},
		}

		expectedPeer1Metas := []dto.IndexMeta{
			{
				PageID:        "1",
				PeerID:        peer1,
				LastIndexedAt: now,
			},
			{
				PageID:        "2",
				PeerID:        peer1,
				LastIndexedAt: oneHourAgo,
			},
		}

		expectedPeer1PageDumps := []dto.PageDump{
			{
				PeerID: "peer1",
				PageID: "1",
				Page: dto.Page{
					ID:            "1",
					Title:         "Google",
					LastIndexedAt: now,
				},
				KeywordOccurrences: map[string]dto.KeywordOccurrence{
					"keyword1": {
						PageID:    "1",
						Frequency: 1,
						Positions: []int{1, 2, 3},
					},
				},
			},

			{
				PeerID: "peer1",
				PageID: "2",
				Page: dto.Page{
					ID:            "2",
					Title:         "Example",
					LastIndexedAt: oneHourAgo,
				},
				KeywordOccurrences: map[string]dto.KeywordOccurrence{
					"keyword1": {
						PageID:    "2",
						Frequency: 1,
						Positions: []int{1},
					},
				},
			},
		}

		expectedPeer2Metas := []dto.IndexMeta{
			{
				PageID:        "3",
				PeerID:        peer2,
				LastIndexedAt: now,
			},
		}

		expectedPeer2PageDumps := []dto.PageDump{
			{
				PeerID: "peer2",
				PageID: "3",
				Page: dto.Page{
					ID:            "3",
					Title:         "Example Org",
					LastIndexedAt: now,
				},
				KeywordOccurrences: map[string]dto.KeywordOccurrence{
					"keyword1": {
						PageID:    "3",
						Frequency: 1,
						Positions: []int{1},
					},
				},
			},
		}

		expectedPeer1MetasResponse := &dto.GetIndexMetasResponse{
			IndexMetas: expectedPeer1Metas,
		}
		expectedPeer2MetasResponse := &dto.GetIndexMetasResponse{
			IndexMetas: expectedPeer2Metas,
		}

		gock.New("http://peer1:8080").
			Get("/repair/get-all-index-metas").
			Reply(200).
			JSON(expectedPeer1MetasResponse)

		gock.New("http://peer2:8080").
			Get("/repair/get-all-index-metas").
			Reply(200).
			JSON(expectedPeer2MetasResponse)

		gock.New("http://peer1:8080").
			Post("/repair/get-index-metas").
			Reply(200).
			JSON(expectedPeer1MetasResponse)

		gock.New("http://peer2:8080").
			Post("/repair/get-index-metas").
			Reply(200).
			JSON(expectedPeer2MetasResponse)

		gock.New("http://peer1:8080").
			Post("/repair/get-page-dumps").
			JSON(&dto.GetPageDumpsRequest{
				PageIDs: []string{"1", "2"},
			}).
			Reply(200).
			JSON(&dto.GetPageDumpsResponse{
				PageDumps: expectedPeer1PageDumps,
			})

		gock.New("http://peer2:8080").
			Post("/repair/get-page-dumps").
			JSON(&dto.GetPageDumpsRequest{
				PageIDs: []string{"3"},
			}).
			Reply(200).
			JSON(&dto.GetPageDumpsResponse{
				PageDumps: expectedPeer2PageDumps,
			})

		err := repairJobService.AuditAndRepair()

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		for _, meta := range expectedMetas {
			page, err := pageRepo.Get(string(meta.PageID))
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if page.LastIndexedAt == nil || !page.LastIndexedAt.Round(time.Second).Equal(meta.LastIndexedAt.Round(time.Second)) {
				t.Errorf("Expected page %s last indexed at to be %v, got %v", meta.PageID, meta.LastIndexedAt, page.LastIndexedAt)
			}
		}

		if !gock.IsDone() {
			t.Errorf("Expected all mocks to be called")
		}
	})

	t.Run("can audit and repair pages (non exist)", func(t *testing.T) {
		defer gock.Off()

		now := time.Now().Round(time.Second)
		oneHourAgo := now.Add(-time.Hour).Round(time.Second)
		peer1 := "peer1"
		peer2 := "peer2"

		pageRepo := pageRepo.NewRepository()
		pageService := pageService.NewService(pageRepo, nil)

		keywordOccurrenceRepo := keywordOccurrenceRepo.NewRepository()
		keywordService := keywordService.NewService(keywordOccurrenceRepo)

		peerService := peerService.NewService(nil, &domain.Peer{
			ID:       "node4",
			Hostname: "node4",
			Port:     8080,
		}, testutil.NewTestLogger())

		peerService.AddPeer(&domain.Peer{
			ID:       peer1,
			Hostname: "peer1",
			Port:     8080,
		})

		peerService.AddPeer(&domain.Peer{
			ID:       "peer2",
			Hostname: "peer2",
			Port:     8080,
		})

		repairJobRepo := repairJobRepo.NewRepository()
		repairJobService := repairJobService.NewService(repairJobRepo, pageService, keywordService, peerService, testutil.NewTestLogger())

		expectedMetas := []domain.IndexMeta{
			{
				PageID:        "1",
				PeerID:        domain.PeerID(peer1),
				LastIndexedAt: now,
			},
			{
				PageID:        "2",
				PeerID:        domain.PeerID(peer1),
				LastIndexedAt: oneHourAgo,
			},
			{
				PageID:        "3",
				PeerID:        domain.PeerID(peer2),
				LastIndexedAt: now,
			},
		}

		expectedPeer1Metas := []dto.IndexMeta{
			{
				PageID:        "1",
				PeerID:        peer1,
				LastIndexedAt: now,
			},
			{
				PageID:        "2",
				PeerID:        peer1,
				LastIndexedAt: oneHourAgo,
			},
		}

		expectedPeer1PageDumps := []dto.PageDump{
			{
				PeerID: "peer1",
				PageID: "1",
				Page: dto.Page{
					ID:            "1",
					Title:         "Google",
					LastIndexedAt: now,
				},
				KeywordOccurrences: map[string]dto.KeywordOccurrence{
					"keyword1": {
						PageID:    "1",
						Frequency: 1,
						Positions: []int{1, 2, 3},
					},
				},
			},

			{
				PeerID: "peer1",
				PageID: "2",
				Page: dto.Page{
					ID:            "2",
					Title:         "Example",
					LastIndexedAt: oneHourAgo,
				},
				KeywordOccurrences: map[string]dto.KeywordOccurrence{
					"keyword1": {
						PageID:    "2",
						Frequency: 1,
						Positions: []int{1},
					},
				},
			},
		}

		expectedPeer2Metas := []dto.IndexMeta{
			{
				PageID:        "3",
				PeerID:        peer2,
				LastIndexedAt: now,
			},
		}

		expectedPeer2PageDumps := []dto.PageDump{
			{
				PeerID: "peer2",
				PageID: "3",
				Page: dto.Page{
					ID:            "3",
					Title:         "Example Org",
					LastIndexedAt: now,
				},
				KeywordOccurrences: map[string]dto.KeywordOccurrence{
					"keyword1": {
						PageID:    "3",
						Frequency: 1,
						Positions: []int{1},
					},
				},
			},
		}

		expectedPeer1MetasResponse := &dto.GetIndexMetasResponse{
			IndexMetas: expectedPeer1Metas,
		}
		expectedPeer2MetasResponse := &dto.GetIndexMetasResponse{
			IndexMetas: expectedPeer2Metas,
		}

		gock.New("http://peer1:8080").
			Get("/repair/get-all-index-metas").
			Reply(200).
			JSON(expectedPeer1MetasResponse)

		gock.New("http://peer2:8080").
			Get("/repair/get-all-index-metas").
			Reply(200).
			JSON(expectedPeer2MetasResponse)

		gock.New("http://peer1:8080").
			Post("/repair/get-index-metas").
			Reply(200).
			JSON(expectedPeer1MetasResponse)

		gock.New("http://peer2:8080").
			Post("/repair/get-index-metas").
			Reply(200).
			JSON(expectedPeer2MetasResponse)

		gock.New("http://peer1:8080").
			Post("/repair/get-page-dumps").
			JSON(&dto.GetPageDumpsRequest{
				PageIDs: []string{"1", "2"},
			}).
			Reply(200).
			JSON(&dto.GetPageDumpsResponse{
				PageDumps: expectedPeer1PageDumps,
			})

		gock.New("http://peer2:8080").
			Post("/repair/get-page-dumps").
			JSON(&dto.GetPageDumpsRequest{
				PageIDs: []string{"3"},
			}).
			Reply(200).
			JSON(&dto.GetPageDumpsResponse{
				PageDumps: expectedPeer2PageDumps,
			})

		err := repairJobService.AuditAndRepair()

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		for _, meta := range expectedMetas {
			page, err := pageRepo.Get(string(meta.PageID))
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if page.LastIndexedAt == nil || !page.LastIndexedAt.Round(time.Second).Equal(meta.LastIndexedAt.Round(time.Second)) {
				t.Errorf("Expected page %s last indexed at to be %v, got %v", meta.PageID, meta.LastIndexedAt, page.LastIndexedAt)
			}
		}

		if !gock.IsDone() {
			t.Errorf("Expected all mocks to be called")
		}
	})
}
