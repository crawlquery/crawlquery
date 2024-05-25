package node

import (
	"bytes"
	"crawlquery/node/dto"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
	}
}

func (c *Client) Crawl(pageID, url string) (*dto.Page, error) {
	req := dto.CrawlRequest{
		PageID: pageID,
		URL:    url,
	}

	jsonBody, err := json.Marshal(req)

	if err != nil {
		return nil, err
	}

	res, err := http.Post(
		c.BaseURL+"/crawl",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusBadRequest {

		var errRes dto.ErrorResponse

		err = json.NewDecoder(res.Body).Decode(&errRes)

		if err != nil {
			return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}

		defer res.Body.Close()

		return nil, fmt.Errorf("unexpected status code: %d (%s)", res.StatusCode, errRes.Error)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var crawlRes dto.CrawlResponse
	err = json.NewDecoder(res.Body).Decode(&crawlRes)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return crawlRes.Page, nil
}
