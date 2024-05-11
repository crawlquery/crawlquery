package service

import (
	"crawlquery/pkg/domain"
	"encoding/json"
	"fmt"
	"net/http"
)

type SearchService struct{}

func NewSearchService() *SearchService {
	return &SearchService{}
}

func (s *SearchService) Search(term string) []domain.Result {
	nodes := []string{
		"http://localhost:9090",
	}

	results := []domain.Result{}

	for _, node := range nodes {
		res, err := http.Get(node + "/search?q=" + term)
		if err != nil {
			fmt.Println(err)
			continue
		}

		defer res.Body.Close()

		var response struct {
			Results []domain.Result `json:"results"`
		}

		err = json.NewDecoder(res.Body).Decode(&response)

		if err != nil {
			fmt.Println(err)
			continue
		}
		results = append(results, response.Results...)
	}

	return results
}
