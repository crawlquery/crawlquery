package service

import (
	"context"
	"crawlquery/api/domain"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	nodeDomain "crawlquery/node/domain"

	"crawlquery/pkg/dto"

	"go.uber.org/zap"
)

type Service struct {
	nodeService     domain.NodeService
	pageRankService domain.PageRankService
	logger          *zap.SugaredLogger
}

func NewService(nodeService domain.NodeService, pageRankService domain.PageRankService, logger *zap.SugaredLogger) *Service {
	return &Service{
		nodeService:     nodeService,
		pageRankService: pageRankService,
		logger:          logger,
	}
}

// Search searches for the term and waits for the fastest node in each shard.
func (s *Service) Search(term string) ([]nodeDomain.Result, error) {

	// trim space either side of the term
	term = strings.TrimSpace(term)
	// remove duplicate spaces
	term = strings.Join(strings.Fields(term), " ")

	shardNodes, err := s.nodeService.ListGroupByShard()
	if err != nil {
		return nil, err
	}

	if len(shardNodes) == 0 {
		s.logger.Errorf("Search.Service.Search: No nodes found")
		return nil, domain.ErrInternalError
	}

	var results []nodeDomain.Result
	var resultsLock sync.Mutex
	var wg sync.WaitGroup

	wg.Add(len(shardNodes))

	for _, nodes := range shardNodes {
		go func(nodes []*domain.Node) {
			defer wg.Done()

			if len(nodes) > 10 {
				nodes = nodes[:10]
			}

			// Initialize results channel with buffer size of the number of nodes
			resultsChan := make(chan []nodeDomain.Result, len(nodes))
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			for _, node := range nodes {
				go func(node *domain.Node) {
					endpoint := fmt.Sprintf("http://%s:%d/search?q=%s", node.Hostname, node.Port, url.QueryEscape(term))
					res, err := http.Get(endpoint) // Simplified for clarity, consider handling with context
					if err != nil {
						s.logger.Errorf("Error searching node %s: %v", node.ID, err)
						return
					}
					defer res.Body.Close()

					var response dto.NodeSearchResponse
					if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
						s.logger.Errorf("Error decoding response from node %s: %v", node.ID, err)
						return
					}
					resultsChan <- response.Results
				}(node)
			}

			// Wait for the first result
			select {
			case res := <-resultsChan:
				resultsLock.Lock()
				results = append(results, res...)
				resultsLock.Unlock()
			case <-ctx.Done():
				s.logger.Errorf("Search timed out for shard %d", nodes[0].ShardID)
			}
		}(nodes)
	}

	wg.Wait()

	// filter out duplicate results
	uniqueResults := make(map[string]nodeDomain.Result)

	for _, res := range results {
		if _, ok := uniqueResults[res.PageID]; !ok {
			rank, err := s.pageRankService.GetPageRank(res.PageID)

			if err != nil {
				s.logger.Errorf("No pagerank found for %s: %v", res.PageID, err)
			}
			res.PageRank = rank
			uniqueResults[res.PageID] = res
		}
	}

	results = make([]nodeDomain.Result, 0, len(uniqueResults))

	for _, res := range uniqueResults {
		results = append(results, res)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if len(results) > 10 {
		results = results[:10]
	}

	return results, nil
}
