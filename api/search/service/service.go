package service

import (
	"crawlquery/api/domain"
	"crawlquery/pkg/dto"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	sharedDomain "crawlquery/pkg/domain"

	"go.uber.org/zap"
)

type Service struct {
	nodeService domain.NodeService
	logger      *zap.SugaredLogger
}

func NewService(nodeService domain.NodeService, logger *zap.SugaredLogger) *Service {
	return &Service{
		nodeService: nodeService,
		logger:      logger,
	}
}

func (s *Service) Search(term string) ([]sharedDomain.Result, error) {
	shardNodes, err := s.nodeService.ListGroupByShard()

	if err != nil {
		return nil, err
	}

	results := []sharedDomain.Result{}
	resultsLock := sync.Mutex{}

	var wg sync.WaitGroup

	wg.Add(len(shardNodes))

	for _, nodes := range shardNodes {
		go func(nodes []*domain.Node) {
			defer wg.Done()
			for _, node := range nodes {
				endpoint := fmt.Sprintf("http://%s:%d/search?q=%s", node.Hostname, node.Port, strings.Replace(term, " ", "%20", -1))
				res, err := http.Get(endpoint)
				if err != nil {
					s.logger.Errorf("Search.Service.Search: Error searching node %s: %v", node.ID, err)
					continue
				}
				defer res.Body.Close()

				var response dto.NodeSearchResponse

				err = json.NewDecoder(res.Body).Decode(&response)

				if err != nil {
					s.logger.Errorf("Search.Service.Search: Error decoding response from node %s: %v", node.ID, err)
					continue
				}
				resultsLock.Lock()
				results = append(results, response.Results...)
				resultsLock.Unlock()

				break
			}
		}(nodes)
	}

	wg.Wait()

	return results, nil
}
