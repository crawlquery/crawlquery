package service_test

import (
	"crawlquery/node/domain"
	"crawlquery/node/dto"
	repairJobRepo "crawlquery/node/repair/job/repository/mem"
	repairJobService "crawlquery/node/repair/service"

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
			IndexMetas: []*dto.IndexMeta{
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
			PageDumps: []*dto.PageDump{
				{
					PeerID: "peer1",
					PageID: "3",
					Page: &dto.Page{
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
			IndexMetas: []*dto.IndexMeta{
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
			PageDumps: []*dto.PageDump{
				{
					PeerID: "peer1",
					PageID: "1",
					Page: &dto.Page{
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
					Page: &dto.Page{
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
				keywordOccurrences[0].PageID = "1"
				keywordOccurrences[0].Frequency = 1
				keywordOccurrences[0].Positions = []int{1, 2, 3}
			}

			if job == "2" {
				keywordOccurrences[0].PageID = "2"
				keywordOccurrences[0].Frequency = 1
				keywordOccurrences[0].Positions = []int{1}
			}

			if job == "3" {
				keywordOccurrences[0].PageID = "3"
				keywordOccurrences[0].Frequency = 1
				keywordOccurrences[0].Positions = []int{1}
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
