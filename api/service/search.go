package service

import (
	"crawlquery/pkg/domain"
	"crawlquery/pkg/dto"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type SearchService struct {
	nodeService domain.NodeService
}

func NewSearchService(nodeService domain.NodeService) *SearchService {
	return &SearchService{
		nodeService: nodeService,
	}
}

func (s *SearchService) Search(term string) ([]domain.Result, error) {
	shardNodes, err := s.nodeService.AllByShard()

	if err != nil {
		return nil, err
	}

	results := []domain.Result{}
	resultsLock := sync.Mutex{}

	var wg sync.WaitGroup

	wg.Add(len(shardNodes))

	for _, nodes := range shardNodes {
		go func(nodes []*domain.Node) {
			defer wg.Done()
			for _, node := range nodes {
				endpoint := fmt.Sprintf("http://%s:%d/search?q=%s", node.Hostname, node.Port, term)
				res, err := http.Get(endpoint)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("made id?")
				defer res.Body.Close()

				var response dto.NodeSearchResponse

				err = json.NewDecoder(res.Body).Decode(&response)

				if err != nil {
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
